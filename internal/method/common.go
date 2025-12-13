package method

import (
	"context"
	"database/sql"

	"codeberg.org/socialmaps/api/internal/model"
)

type Common struct {
	DB *sql.DB
}

type Response[B any] struct {
	Body B
}

func GetAuthUser(ctx context.Context) *model.User {
	usr, ok := ctx.Value("auth.user").(*model.User)
	if !ok || usr == nil {
		panic("cannot get auth user from context")
	}
	return usr
}
