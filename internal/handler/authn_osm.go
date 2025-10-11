package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

type AuthnOSM struct {
	Common
	OAuth2Config *oauth2.Config
}

func (h *AuthnOSM) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	authCodeURL := cfg.AuthCodeURL("foo", oauth2.ApprovalForce)
	log.Println(authCodeURL)
	http.Redirect(w, r, authCodeURL, http.StatusSeeOther)
}
