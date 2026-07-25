package render

import (
	"golang.socialmaps.org/api/internal/model"
	"golang.socialmaps.org/api/internal/resource"
)

func User(m model.User) resource.User {
	return resource.User{
		ID:          m.ID,
		DisplayName: m.DisplayName,
	}
}
