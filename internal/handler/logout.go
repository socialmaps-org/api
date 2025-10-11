package handler

import (
	"net/http"

	"codeberg.org/socialmaps/auth/internal/model"
	"codeberg.org/socialmaps/auth/internal/session"
	"codeberg.org/socialmaps/auth/internal/templates"
)

type Logout struct {
	Common
}

func (h *Logout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	http.SetCookie(w, session.EmptyCookie())

	cookie, err := r.Cookie(session.COOKIE_NAME)
	if err != nil && err != http.ErrNoCookie {
		panic(err)
	}

	if cookie != nil {
		sessionID := session.FromCookie(h.Env.CookieSecret, cookie)
		ses := model.LoadActiveSession(ctx, h.DB, sessionID)
		if ses != nil {
			model.RevokeSession(ctx, h.DB, ses.ID)
		}
	}

	err = templates.Logout.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}
