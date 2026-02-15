package render

import (
	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/resource"
)

func Review(m model.Review) resource.Review {
	var comment string
	if m.Comment.Valid {
		comment = m.Comment.String
	}

	return resource.Review{
		Object:  "review",
		ID:      m.ID,
		Created: m.Created,
		Place: resource.PlaceStub{
			ID: m.PlaceID,
		},
		User: resource.UserStub{
			ID: m.UserID,
		},
		Liked:   m.Liked,
		Comment: comment,
	}
}

func ReviewWithUser(rvw model.Review, usr model.User) resource.ReviewWithUser {
	return resource.ReviewWithUser{
		Review: Review(rvw),
		User:   User(usr),
	}
}
