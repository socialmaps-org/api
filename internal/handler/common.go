package handler

import (
	"database/sql"

	"codeberg.org/socialmaps/auth/internal/env"
)

type Common struct {
	DB  *sql.DB
	Env env.AuthEnv
}
