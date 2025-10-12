package handler

import (
	"encoding/json"
	"net/http"
	"net/url"

	"codeberg.org/socialmaps/auth/internal/contrib/nonce"
	"codeberg.org/socialmaps/auth/internal/model"
	"codeberg.org/socialmaps/auth/internal/resource"
	"codeberg.org/socialmaps/auth/internal/session"
	"codeberg.org/socialmaps/auth/internal/web"
	"golang.org/x/oauth2"
)

type AuthnCallback struct {
	Common
	NonceService *nonce.NonceService
	OAuth2Config *oauth2.Config
}

func (h *AuthnCallback) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Create a copy of the OAuth2 config as we're changing its RedirectURL
	// based on the current request specifically.
	cfg := *h.OAuth2Config
	redirectURI := r.FormValue("redirect_uri")
	if redirectURI != "" {
		cfg.RedirectURL += "?redirect_uri=" + url.QueryEscape(redirectURI)
	}

	code := r.FormValue("code")
	state := r.FormValue("state")

	if ok := h.NonceService.Valid(state); !ok {
		web.JSON(w, http.StatusBadRequest, &resource.Error{
			Message: "query parameter `state` is invalid or missing",
		})
		return
	}

	tok, err := cfg.Exchange(ctx, code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		web.Flush(w)
		panic(err)
	}

	client := cfg.Client(ctx, tok)
	res, err := client.Get("https://www.openstreetmap.org/oauth2/userinfo")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		web.Flush(w)
		panic(err)
	}

	var userinfo struct {
		Sub               string `json:"sub"`
		PreferredUsername string `json:"preferred_username"`
	}
	err = json.NewDecoder(res.Body).Decode(&userinfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		web.Flush(w)
		panic(err)
	}

	user := model.CreateOrUpdateUser(ctx, h.DB, "org.openstreetmap", userinfo.Sub, userinfo.PreferredUsername)

	ses := model.CreateSession(ctx, h.DB, user.ID)
	sesCookie := session.ToCookie(h.Env.CookieSecret, ses.ID)

	http.SetCookie(w, sesCookie)

	if redirectURI != "" {
		// TODO: validate redirectURI
		http.Redirect(w, r, redirectURI, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
