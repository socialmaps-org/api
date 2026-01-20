package method

import (
	"log/slog"
	"net/http"
	"regexp"
	"strconv"

	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/web"
	"github.com/danielgtaylor/huma/v2"
)

func GetAuthMiddleware(
	api huma.API,
	authr web.Authenticator,
	qs *model.Queries,
) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		re := regexp.MustCompile(`^Bearer\s+(?P<token>[a-zA-Z0-9._~+/=-]+)$`)

		authz := ctx.Header("Authorization")
		matches := re.FindStringSubmatch(authz)
		idx := re.SubexpIndex("token")

		if len(matches) < re.SubexpIndex("token") {
			slog.Info("CANONICAL-AUTH-LINE",
				"status", "error",
				"error", "malformed or missing Authorization header",
			)
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "bearer token is missing or malformed")
			// w.Header().Set(
			// 	"WWW-Authenticate",
			// 	`Bearer error="invalid_token", error_description="The access token is missing."`,
			// )
			return
		}
		token := matches[idx]

		authn, err := authr.Introspect(token)
		if err != nil {
			slog.Info("CANONICAL-AUTH-LINE",
				"status", "error",
				"error", err.Error(),
			)
			huma.WriteErr(api, ctx, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		slog.Info("CANONICAL-AUTH-LINE",
			"status", "success",
			"active", authn.Active,
			"username", authn.Username,
		)

		if !authn.Active {
			// w.Header().Set(
			// 	"WWW-Authenticate",
			// 	`Bearer error="invalid_token", error_description="The access token is invalid or has expired."`,
			// )
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "bearer token is invalid or expired")
			return
		}

		osmID, err := strconv.ParseInt(authn.OpenStreetMapSub, 10, 64)
		if err != nil {
			slog.Info("CANONICAL-AUTH-LINE",
				"status", "error",
				"error", err.Error(),
			)
			huma.WriteErr(api, ctx, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		usr, err := qs.CreateUser(ctx.Context(), osmID, authn.Username)
		if err != nil {
			slog.Info("CANONICAL-AUTH-LINE",
				"status", "error",
				"error", err.Error(),
			)
			huma.WriteErr(api, ctx, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		next(huma.WithValue(ctx, "auth.user", usr))
	}
}
