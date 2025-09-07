package main

import (
	"encoding/base64"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"codeberg.org/socialmaps/auth/internal/database"
	"codeberg.org/socialmaps/auth/internal/model"
	"codeberg.org/socialmaps/auth/internal/session"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

type userinfo struct {
	Sub               string `json:"sub"`
	PreferredUsername string `json:"preferred_username"`
}

func main() {
	OSMClientID, ok := os.LookupEnv("OSM_CLIENT_ID")
	if !ok {
		panic("Environment variable `OSM_CLIENT_ID` is not set")
	}
	OSMClientSecret, ok := os.LookupEnv("OSM_CLIENT_SECRET")
	if !ok {
		panic("Environment variable `OSM_CLIENT_SECRET` is not set")
	}
	SESSION_SECRET_B64, ok := os.LookupEnv("SESSION_SECRET")
	if !ok {
		panic("Environment variable `SESSION_SECRET` is not set")
	}
	SESSION_SECRET, err := base64.StdEncoding.DecodeString(SESSION_SECRET_B64)
	if err != nil {
		panic("Could not decode `SESSION_SECRET`")
	}
	if len(SESSION_SECRET) < 32 {
		panic("`SESSION_SECRET` is too short; must be at least 32 bytes when decoded")
	}

	oauth2Config := &oauth2.Config{
		ClientID:     OSMClientID,
		ClientSecret: OSMClientSecret,
		Scopes:       []string{"openid"},
		Endpoint:     endpoints.OpenStreetMap,
		RedirectURL:  "http://127.0.0.1:8080/auth/callback",
	}

	loginTemplate := template.Must(template.New("login.gohtml").ParseFiles("web/template/login.gohtml"))

	db := database.Open("db.sqlite3")
	defer db.Close()

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		type loginTemplateData struct {
			User *model.User
		}

		var usr *model.User
		cookie, err := r.Cookie(session.COOKIE_NAME)
		if err != nil && err != http.ErrNoCookie {
			panic(err)
		}

		if cookie != nil {
			sessionID := session.FromCookie(SESSION_SECRET, cookie)
			session := model.LoadActiveSession(db, sessionID)
			if session != nil {
				usr = model.LoadUser(db, session.UserID)
			}
		}

		err = loginTemplate.Execute(
			w,
			loginTemplateData{
				User: usr,
			},
		)
		if err != nil {
			log.Println(err)
		}
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(session.COOKIE_NAME)
		if err != nil && err != http.ErrNoCookie {
			panic(err)
		}

		if cookie == nil {
			return
		}

		sessionID := session.FromCookie(SESSION_SECRET, cookie)
		ses := model.LoadActiveSession(db, sessionID)
		if ses == nil {
			return
		}

		ses.Revoke()

		http.SetCookie(w, session.EmptyCookie())

		w.Write([]byte("logged out!"))
	})

	http.HandleFunc("/auth/openstreetmap", func(w http.ResponseWriter, r *http.Request) {
		authCodeURL := oauth2Config.AuthCodeURL("foo", oauth2.ApprovalForce)
		log.Println(authCodeURL)
		http.Redirect(w, r, authCodeURL, http.StatusSeeOther)
	})

	http.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		tok, err := oauth2Config.Exchange(r.Context(), code)
		if err != nil {
			panic(err)
		}

		client := oauth2Config.Client(r.Context(), tok)

		res, err := client.Get("https://www.openstreetmap.org/oauth2/userinfo")
		if err != nil {
			panic(err)
		}

		var userinfo userinfo
		err = json.NewDecoder(res.Body).Decode(&userinfo)
		if err != nil {
			panic(err)
		}

		user := model.CreateOrUpdateUser(db, "org.openstreetmap", userinfo.Sub, userinfo.PreferredUsername)

		ses := model.CreateSession(db, user.ID)
		sesCookie := session.ToCookie(SESSION_SECRET, ses.ID)

		http.SetCookie(w, sesCookie)

		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			panic(err)
		}
	})

	log.Println("serving...")
	err = http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
