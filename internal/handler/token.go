package handler

import (
	"net/http"
	"time"

	"codeberg.org/socialmaps/auth/internal/web"
	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"
)

type Token struct {
	Common
	OAuth2Server fosite.OAuth2Provider
}

func (h *Token) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	accessRequest, err := h.OAuth2Server.NewAccessRequest(ctx, r, mySessionData)
	if err != nil {
		h.OAuth2Server.WriteAccessError(ctx, w, accessRequest, err)
		web.Flush(w)
		panic(err)
	}

	accessRequest.GrantScope("openid")

	response, err := h.OAuth2Server.NewAccessResponse(ctx, accessRequest)
	if err != nil {
		h.OAuth2Server.WriteAccessError(ctx, w, accessRequest, err)
		web.Flush(w)
		panic(err)
	}

	h.OAuth2Server.WriteAccessResponse(ctx, w, accessRequest, response)
}
