package render

import (
	"codeberg.org/socialmaps/auth/internal/model"
	"codeberg.org/socialmaps/auth/internal/resource"
)

func Review(m *model.Review) *resource.Review {
	return &resource.Review{
		Object:  "review",
		ID:      m.ID,
		Created: m.Created,
		Place: resource.Place{
			ID: m.PlaceID,
		},
		User: resource.User{
			ID: m.UserID,
		},
		Liked:   m.Liked,
		Comment: m.Comment,
	}
}
