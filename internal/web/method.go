package web

import (
	"context"
	"log/slog"
	"net/http"

	"codeberg.org/socialmaps/api/internal/resource"
)

type Method[A any] interface {
	NewArgs() *A
	Validate(args *A) *Response
	Execute(ctx context.Context, args *A) *Response
}

func MethodHandler[A any](method Method[A]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		args := method.NewArgs()
		if args != nil {
			if err := parseRequest(r, args); err != nil {
				res := NewJSONResponse(
					http.StatusBadRequest,
					resource.Error{
						Message: err.Error(),
					},
				)
				slog.InfoContext(ctx, "CANONICAL-METHOD-LINE",
					"http_status", res.StatusCode,
				)
				res.send(w)
				return
			}

			if res := method.Validate(args); res != nil {
				slog.InfoContext(ctx, "CANONICAL-METHOD-LINE",
					"http_status", res.StatusCode,
				)
				res.send(w)
				return
			}
		}

		res := method.Execute(ctx, args)
		if res == nil {
			panic("method returned empty result")
		}

		slog.InfoContext(ctx, "CANONICAL-METHOD-LINE",
			"http_status", res.StatusCode,
		)
		res.send(w)
	})
}
