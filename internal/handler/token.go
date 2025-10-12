package handler

import (
	"fmt"
	"net/http"

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

	session := &openid.DefaultSession{
		Claims:  &jwt.IDTokenClaims{},
		Headers: &jwt.Headers{},
	}

	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\tform: %+v\n", r.Form)

	accessRequest, err := h.OAuth2Server.NewAccessRequest(ctx, r, session)
	if err != nil {
		h.OAuth2Server.WriteAccessError(ctx, w, accessRequest, err)
		web.Flush(w)
		panic(err)
	}

	fmt.Printf("\tsession: %+v\n", session)

	response, err := h.OAuth2Server.NewAccessResponse(ctx, accessRequest)
	if err != nil {
		h.OAuth2Server.WriteAccessError(ctx, w, accessRequest, err)
		web.Flush(w)
		panic(err)
	}

	h.OAuth2Server.WriteAccessResponse(ctx, w, accessRequest, response)
}
