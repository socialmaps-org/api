package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"codeberg.org/socialmaps/auth/internal/model"
	"codeberg.org/socialmaps/auth/internal/session"
	"codeberg.org/socialmaps/auth/internal/web"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

type AuthnCallback struct {
	Common
	OAuth2Config *oauth2.Config
}

func (h *AuthnCallback) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	callbackURL := fmt.Sprintf("http://%s/auth/callback", h.Env.Host)

	redirectURI := r.FormValue("redirect_uri")
	if redirectURI != "" {
		callbackURL += "?redirect_uri=" + url.QueryEscape(redirectURI)
	}

	cfg := oauth2.Config{
		ClientID:     h.Env.OSMClientID,
		ClientSecret: h.Env.OSMClientSecret,
		Scopes:       []string{"openid"},
		Endpoint:     endpoints.OpenStreetMap,
		RedirectURL:  callbackURL,
	}

	code := r.FormValue("code")
	tok, err := cfg.Exchange(ctx, code)
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

	client := cfg.Client(ctx, tok)

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
