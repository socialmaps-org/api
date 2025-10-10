package web

import (
	"context"
	"net/http"

	"codeberg.org/socialmaps/auth/internal/resource"
)

type Method[A any] interface {
	NewArgs() *A
	Validate(args *A) *Response
	Execute(ctx context.Context, args *A) *Response
}

func MethodHandler[A any](method Method[A]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		args := method.NewArgs()
		if args == nil {
			panic("method.NewArgs() returned nil")
		}

		if err := Parse(r, args); err != nil {
			NewJSONResponse(
				http.StatusBadRequest,
				resource.Error{
					Message: err.Error(),
				},
			).send(w)
			return
		}

		if res := method.Validate(args); res != nil {
			res.send(w)
			return
		}

		res := method.Execute(r.Context(), args)
		if res != nil {
			res.send(w)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	})
}
