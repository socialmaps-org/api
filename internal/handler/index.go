package handler

import (
	"log/slog"
	"net/http"

	"codeberg.org/socialmaps/auth/internal/model"
	"codeberg.org/socialmaps/auth/internal/session"
	"codeberg.org/socialmaps/auth/internal/templates"
)

type Index struct {
	Common
}

func (h *Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type indexTemplateData struct {
		User *model.User
	}

	var usr *model.User
	cookie, err := r.Cookie(session.COOKIE_NAME)
	if err != nil && err != http.ErrNoCookie {
		panic(err)
	}

	var sessionID string
	var ses *model.Session
	if cookie != nil {
		sessionID = session.FromCookie(h.Env.CookieSecret, cookie)
		ses = model.LoadActiveSession(ctx, h.DB, sessionID)
		if ses != nil {
			usr = model.LoadUser(ctx, h.DB, ses.UserID)
		}
	}

	slog.InfoContext(ctx, "CANONICAL-AUTH-REQUEST",
		"handler", "index",
		"cookie_present", cookie != nil,
		"session_id", sessionID,
		"session_found", ses != nil,
		"user_found", usr != nil,
	)

	templates.Index.Execute(
		w,
		indexTemplateData{
			User: usr,
		},
	)
}
