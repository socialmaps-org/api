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
		Place: resource.Place{
			ID: m.PlaceID,
		},
		User: resource.User{
			ID: m.UserID,
		},
		Liked:   m.Liked,
		Comment: comment,
	}
}

func Reviews(ms []model.Review) resource.List[resource.Review] {
	var rs []resource.Review
	for _, m := range ms {
		rs = append(rs, Review(m))
	}

	return resource.List[resource.Review]{
		Object: "list",
		Data:   rs,
	}
}
