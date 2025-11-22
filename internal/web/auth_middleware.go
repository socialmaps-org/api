package web

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"regexp"

	"codeberg.org/socialmaps/api/internal/model"
)

type contextKeyType int

const (
	CtxKeyUser contextKeyType = iota
)

var re = regexp.MustCompile(`^Bearer\s+(?P<token>[a-zA-Z0-9._~+/=-]+)$`)

func GetAuthUser(ctx context.Context) *model.User {
	usr, ok := ctx.Value(CtxKeyUser).(*model.User)
	if !ok || usr == nil {
		panic("cannot get auth user from context")
	}
	return usr
}

func AuthMiddleware(db *sql.DB, authr Authenticator, scope string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authz := r.Header.Get("Authorization")
		matches := re.FindStringSubmatch(authz)
		idx := re.SubexpIndex("token")

		if len(matches) < re.SubexpIndex("token") {
			slog.Info("CANONICAL-AUTH-LINE",
				"status", "error",
				"error", "malformed or missing Authorization header",
			)
			w.Header().Set(
				"WWW-Authenticate",
				`Bearer error="invalid_token", error_description="The access token is missing."`,
			)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := matches[idx]

		authn, err := authr.Introspect(token, scope)
		if err != nil {
			slog.Info("CANONICAL-AUTH-LINE",
				"status", "error",
				"error", err.Error(),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		slog.Info("CANONICAL-AUTH-LINE",
			"status", "success",
			"active", authn.Active,
			"username", authn.Username,
		)

		if !authn.Active {
			w.Header().Set(
				"WWW-Authenticate",
				`Bearer error="invalid_token", error_description="The access token is invalid or has expired."`,
			)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		usr := model.UpsertUser(ctx, db, authn.OpenStreetMapSub, authn.Username)
		nctx := context.WithValue(ctx, CtxKeyUser, usr)

		next.ServeHTTP(w, r.WithContext(nctx))
	})
}
