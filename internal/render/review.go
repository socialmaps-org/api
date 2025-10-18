package render

import (
	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/resource"
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

func Reviews(ms []*model.Review) *resource.List {
	data := make([]any, 0, len(ms))
	for _, m := range ms {
		data = append(data, Review(m))
	}

	return &resource.List{
		Object: "list",
		Data:   data,
	}
}
