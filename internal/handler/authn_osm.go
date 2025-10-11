package handler

import (
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

type AuthnOSM struct {
	Common
	OAuth2Config *oauth2.Config
}

func (h *AuthnOSM) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authCodeURL := h.OAuth2Config.AuthCodeURL("foo", oauth2.ApprovalForce)
	log.Println(authCodeURL)
	http.Redirect(w, r, authCodeURL, http.StatusSeeOther)
}
