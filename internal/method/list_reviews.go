package method

import (
	"context"
	"database/sql"
	"math"

	"codeberg.org/socialmaps/api/internal/render"
	"codeberg.org/socialmaps/api/internal/resource"
)

type ListReviews struct {
	Common
}

type listReviewsArgs struct {
	PlaceID int64 `path:"place_id" minimum:"1" doc:"Unique identifier for the **Place** the user is listing **Review** for."`
	Limit   int64 `query:"limit" minimum:"1" maximum:"100" default:"20" doc:"Maximum number of **Review**s to return at a time."`

	LastCreated int64 `query:"last_created" hidden:"true" dependentRequired:"LastID"`
	LastID      int64 `query:"last_id" hidden:"true"`

	FirstCreated int64 `query:"first_created" hidden:"true" dependentRequired:"FirstID"`
	FirstID      int64 `query:"first_id" hidden:"true"`
}

func (m *ListReviews) Execute(ctx context.Context, args *listReviewsArgs) (*Response[resource.List[resource.ReviewWithUser]], error) {
	var rvwRs []resource.ReviewWithUser
	if args.FirstCreated != 0 {
		results, err := m.QS.ListLatestApprovedReviewsOfPlaceReverse(
			ctx, args.PlaceID, args.FirstCreated, args.FirstID, args.Limit,
		)

		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		for _, res := range results {
			rvwR := render.ReviewWithUser(res.Review, res.User)
			rvwRs = append(rvwRs, rvwR)
		}
	} else {
		var lastCreated, lastID int64
		if args.LastCreated == 0 {
			lastCreated, lastID = math.MaxInt64, math.MaxInt64
		} else {
			lastCreated, lastID = args.LastCreated, args.LastID
		}

		results, err := m.QS.ListLatestApprovedReviewsOfPlace(
			ctx, args.PlaceID, lastCreated, lastID, args.Limit,
		)

		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		for _, res := range results {
			rvwR := render.ReviewWithUser(res.Review, res.User)
			rvwRs = append(rvwRs, rvwR)
		}
	}

	return &Response[resource.List[resource.ReviewWithUser]]{
		Body: resource.List[resource.ReviewWithUser]{
			Object: "list",
			Data:   rvwRs,
		},
	}, nil
}
