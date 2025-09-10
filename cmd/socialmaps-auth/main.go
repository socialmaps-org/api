package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"

	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/storage"
	"github.com/ory/fosite/token/jwt"

	"codeberg.org/socialmaps/auth/internal/database"
	"codeberg.org/socialmaps/auth/internal/model"
	"codeberg.org/socialmaps/auth/internal/session"
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

	indexTemplate := template.Must(template.ParseFiles("web/template/index.gohtml", "web/template/base.gohtml"))
	loginTemplate := template.Must(template.ParseFiles("web/template/login.gohtml", "web/template/base.gohtml"))
	logoutTemplate := template.Must(template.ParseFiles("web/template/logout.gohtml", "web/template/base.gohtml"))

	oauth2Config := &oauth2.Config{
		ClientID:     OSMClientID,
		ClientSecret: OSMClientSecret,
		Scopes:       []string{"openid"},
		Endpoint:     endpoints.OpenStreetMap,
		RedirectURL:  "http://127.0.0.1:8080/auth/callback",
	}

	store := storage.NewExampleStore()
	secret := []byte("ndZC4kpqi6f*!3!v2TDYpfqbAyMATT")
	hashedSecret, err := bcrypt.GenerateFromPassword(secret, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	store.Clients["my-client"] = &fosite.DefaultClient{
		ID:           "my-client",
		Secret:       hashedSecret,
		RedirectURIs: []string{"http://127.0.0.1:8000/auth/callback"},
		Scopes:       []string{"openid"},
	}

	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	var oauth2Server = compose.ComposeAllEnabled(
		&fosite.Config{
			IDTokenIssuer:              "org.socialmaps",
			EnforcePKCE:                false,
			AccessTokenLifespan:        time.Minute * 30,
			GlobalSecret:               []byte("some-cool-secret-that-is-32bytes"),
			SendDebugMessagesToClients: true,
			// ...
		},
		store,
		privateKey,
	)

	db := database.Open("db.sqlite3")
	defer db.Close()

	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		type indexTemplateData struct {
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

		err = indexTemplate.Execute(
			w,
			indexTemplateData{
				User: usr,
			},
		)
		if err != nil {
			log.Println(err)
		}
	})

	http.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
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

	http.HandleFunc("GET /logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, session.EmptyCookie())

		cookie, err := r.Cookie(session.COOKIE_NAME)
		if err != nil && err != http.ErrNoCookie {
			panic(err)
		}

		if cookie != nil {
			sessionID := session.FromCookie(SESSION_SECRET, cookie)
			ses := model.LoadActiveSession(db, sessionID)
			if ses != nil {
				ses.Revoke()
			}
		}

		err = logoutTemplate.Execute(w, nil)
		if err != nil {
			panic(err)
		}
	})

	http.HandleFunc("GET /auth/openstreetmap", func(w http.ResponseWriter, r *http.Request) {
		authCodeURL := oauth2Config.AuthCodeURL("foo", oauth2.ApprovalForce)
		log.Println(authCodeURL)
		http.Redirect(w, r, authCodeURL, http.StatusSeeOther)
	})

	http.HandleFunc("GET /auth/callback", func(w http.ResponseWriter, r *http.Request) {
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

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// OAuth2 Server
	// =========================================================================
	http.HandleFunc("/oauth2/authorize", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// AuthorizeRequest will analyze the request and extract important
		// information like scopes, response type and others.
		ar, err := oauth2Server.NewAuthorizeRequest(ctx, r)
		if err != nil {
			panic(err)
			oauth2Server.WriteAuthorizeError(ctx, w, ar, err)
		}

		ar.GrantScope("openid")

		now := time.Now().UTC()
		mySessionData := &openid.DefaultSession{
			Claims: &jwt.IDTokenClaims{
				Issuer:      "https://auth.socialmaps.org",
				Subject:     "bora",
				Audience:    []string{"http://127.0.0.1:8000"},
				ExpiresAt:   now.Add(time.Hour * 6),
				IssuedAt:    now,
				RequestedAt: now,
				AuthTime:    now,
			},
		}

		response, err := oauth2Server.NewAuthorizeResponse(ctx, ar, mySessionData)
		if err != nil {
			panic(err)
		}

		oauth2Server.WriteAuthorizeResponse(ctx, w, ar, response)
	})

	http.HandleFunc("/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		// This context will be passed to all methods.
		ctx := r.Context()

		now := time.Now().UTC()
		mySessionData := &openid.DefaultSession{
			Claims: &jwt.IDTokenClaims{
				Issuer:      "https://auth.socialmaps.org",
				Audience:    []string{"http://127.0.0.1:8000"},
				ExpiresAt:   now.Add(time.Hour * 6),
				IssuedAt:    now,
				RequestedAt: now,
				AuthTime:    now,
			},
		}

		accessRequest, err := oauth2Server.NewAccessRequest(ctx, r, mySessionData)
		if err != nil {
			oauth2Server.WriteAccessError(ctx, w, accessRequest, err)
		}

		accessRequest.GrantScope("openid")

		response, err := oauth2Server.NewAccessResponse(ctx, accessRequest)
		if err != nil {
			oauth2Server.WriteAccessError(ctx, w, accessRequest, err)
		}

		oauth2Server.WriteAccessResponse(ctx, w, accessRequest, response)
	})

	log.Println("serving...")
	err = http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
