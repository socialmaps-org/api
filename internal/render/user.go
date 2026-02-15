package render

import (
	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/resource"
)

func User(m model.User) resource.User {
	return resource.User{
		ID:          m.ID,
		DisplayName: m.DisplayName,
	}
}
