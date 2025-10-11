package handler

import (
	"encoding/json"
	"net/http"

	"codeberg.org/socialmaps/auth/internal/model"
	"codeberg.org/socialmaps/auth/internal/session"
	"codeberg.org/socialmaps/auth/internal/web"
	"golang.org/x/oauth2"
)

type AuthnCallback struct {
	Common
	OAuth2Config *oauth2.Config
}

func (h *AuthnCallback) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	tok, err := h.OAuth2Config.Exchange(r.Context(), code)
	if err != nil {
		rerr := err.(*oauth2.RetrieveError)
		web.JSON(w, http.StatusBadRequest, struct {
			ErrorCode        string `json:"error"`
			ErrorDescription string `json:"error_description"`
			ErrorURI         string `json:"error_uri,omitempty"`
		}{
			ErrorCode:        rerr.ErrorCode,
			ErrorDescription: rerr.ErrorDescription,
			ErrorURI:         rerr.ErrorURI,
		})
		return
	}

	client := h.OAuth2Config.Client(r.Context(), tok)

	res, err := client.Get("https://www.openstreetmap.org/oauth2/userinfo")
	if err != nil {
		panic(err)
	}

	var userinfo struct {
		Sub               string `json:"sub"`
		PreferredUsername string `json:"preferred_username"`
	}
	err = json.NewDecoder(res.Body).Decode(&userinfo)
	if err != nil {
		panic(err)
	}

	user := model.CreateOrUpdateUser(r.Context(), h.DB, "org.openstreetmap", userinfo.Sub, userinfo.PreferredUsername)

	ses := model.CreateSession(r.Context(), h.DB, user.ID)
	sesCookie := session.ToCookie(h.Env.CookieSecret, ses.ID)

	http.SetCookie(w, sesCookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
