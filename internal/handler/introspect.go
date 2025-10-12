package handler

import (
	"fmt"
	"net/http"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"
)

type Introspect struct {
	Common
	OAuth2Server fosite.OAuth2Provider
}

func (h *Introspect) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session := &openid.DefaultSession{
		Claims:  &jwt.IDTokenClaims{},
		Headers: &jwt.Headers{},
	}

	response, err := h.OAuth2Server.NewIntrospectionRequest(ctx, r, session)
	if err != nil && err != fosite.ErrInactiveToken {
		h.OAuth2Server.WriteIntrospectionError(ctx, w, err)
		return
	}

	fmt.Printf("\tsession: %+v\n", session)

	h.OAuth2Server.WriteIntrospectionResponse(ctx, w, response)
}
