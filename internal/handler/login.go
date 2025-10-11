package handler

import (
	"log"
	"net/http"

	"codeberg.org/socialmaps/auth/internal/model"
	"codeberg.org/socialmaps/auth/internal/session"
	"codeberg.org/socialmaps/auth/internal/templates"
)

type Login struct {
	Common
}

func (h *Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type loginTemplateData struct {
		User *model.User
	}

	var usr *model.User
	cookie, err := r.Cookie(session.COOKIE_NAME)
	if err != nil && err != http.ErrNoCookie {
		panic(err)
	}

	if cookie != nil {
		sessionID := session.FromCookie(h.Env.CookieSecret, cookie)
		session := model.LoadActiveSession(r.Context(), h.DB, sessionID)
		if session != nil {
			usr = model.LoadUser(r.Context(), h.DB, session.UserID)
		}
	}

	err = templates.Login.Execute(
		w,
		loginTemplateData{
			User: usr,
		},
	)
	if err != nil {
		log.Println(err)
	}
}
