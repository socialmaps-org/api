package method

import (
	"context"

	"codeberg.org/socialmaps/api/internal/model"
)

type Common struct {
	QS *model.Queries
}

type Response[B any] struct {
	Body B
}

func GetAuthUser(ctx context.Context) model.User {
	usr, ok := ctx.Value("auth.user").(model.User)
	if !ok {
		panic("cannot get auth user from context")
	}
	return usr
}
