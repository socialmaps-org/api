package handler

import (
	"log"
	"net/http"
	"net/url"

	"codeberg.org/socialmaps/auth/internal/contrib/nonce"
	"golang.org/x/oauth2"
)

type AuthnOSM struct {
	Common
	NonceService *nonce.NonceService
	OAuth2Config *oauth2.Config
}

func (h *AuthnOSM) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Create a copy of the OAuth2 config as we're changing its RedirectURL
	// based on the current request specifically.
	cfg := *h.OAuth2Config
	redirectURI := r.FormValue("redirect_uri")
	if redirectURI != "" {
		cfg.RedirectURL += "?redirect_uri=" + url.QueryEscape(redirectURI)
	}

	state, err := h.NonceService.Nonce()
	if err != nil {
		panic(err)
	}

	authCodeURL := cfg.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	log.Println(authCodeURL)
	http.Redirect(w, r, authCodeURL, http.StatusSeeOther)
}
