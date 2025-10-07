package web

import (
	"codeberg.org/socialmaps/auth/internal/resource"
)

type Error struct {
	StatusCode int
	Resource   resource.Error
}
