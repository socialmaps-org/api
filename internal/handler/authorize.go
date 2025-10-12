package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"codeberg.org/socialmaps/auth/internal/model"
	"codeberg.org/socialmaps/auth/internal/session"
	"codeberg.org/socialmaps/auth/internal/web"
	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"
)

type Authorize struct {
	Common
	OAuth2Server fosite.OAuth2Provider
}

func (h *Authorize) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, err := r.Cookie(session.COOKIE_NAME)
	if err != nil && err != http.ErrNoCookie {
		panic(err)
	}

	if cookie == nil {
		redirectURI := url.QueryEscape(fmt.Sprintf("%s?%s", r.URL.Path, r.URL.RawQuery))
		http.Redirect(w, r, "/login?redirect_uri="+redirectURI, http.StatusSeeOther)
		return
	}

	sessionID := session.FromCookie(h.Env.CookieSecret, cookie)
	session := model.LoadActiveSession(r.Context(), h.DB, sessionID)
	var usr *model.User
	if session != nil {
		usr = model.LoadUser(r.Context(), h.DB, session.UserID)
	}

	if usr == nil {
		panic("foo")
	}

	// AuthorizeRequest will analyze the request and extract important
	// information like scopes, response type and others.
	ar, err := h.OAuth2Server.NewAuthorizeRequest(ctx, r)
	if err != nil {
		h.OAuth2Server.WriteAuthorizeError(ctx, w, ar, err)
		web.Flush(w)
		panic(fmt.Errorf("%+v", err))
	}

	for _, scope := range ar.GetRequestedScopes() {
		if scope == "review" || scope == "offline_access" {
			ar.GrantScope(scope)
		}
	}

	now := time.Now().UTC()
	mySessionData := &openid.DefaultSession{
		Claims: &jwt.IDTokenClaims{
			Issuer:      "https://auth.socialmaps.org",
			Subject:     usr.ID,
			Audience:    []string{ar.GetClient().GetID()},
			ExpiresAt:   now.Add(time.Hour * 6),
			IssuedAt:    now,
			RequestedAt: now,
			AuthTime:    now,
		},
		Username: usr.Username,
		Subject:  usr.ID,
	}

	response, err := h.OAuth2Server.NewAuthorizeResponse(ctx, ar, mySessionData)
	if err != nil {
		h.OAuth2Server.WriteAuthorizeError(ctx, w, ar, err)
		web.Flush(w)
		panic(fmt.Errorf("%+v", err))
	}

	h.OAuth2Server.WriteAuthorizeResponse(ctx, w, ar, response)
}
