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

func Reviews(ms []*model.Review) *resource.List[*model.Review] {
	var startingAfter, endingBefore *string
	if len(ms) != 0 {
		startingAfter = &ms[len(ms)-1].ID
		endingBefore = &ms[0].ID
	}

	return &resource.List[*model.Review]{
		Object:        "list",
		Data:          ms,
		StartingAfter: startingAfter,
		EndingBefore:  endingBefore,
	}
}
