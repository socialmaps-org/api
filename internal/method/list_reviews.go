package method

import (
	"context"
	"database/sql"
	"math"

	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/render"
	"codeberg.org/socialmaps/api/internal/resource"
)

type ListReviews struct {
	Common
}

type listReviewsArgs struct {
	PlaceID int64 `path:"place_id" minimum:"0"`
	Limit   int64 `query:"limit" minimum:"1" maximum:"100" default:"20"`

	LastCreated int64 `query:"last_created" hidden:"true" dependentRequired:"LastID"`
	LastID      int64 `query:"last_id" hidden:"true"`

	FirstCreated int64 `query:"first_created" hidden:"true" dependentRequired:"FirstID"`
	FirstID      int64 `query:"first_id" hidden:"true"`
}

func (m *ListReviews) Execute(ctx context.Context, args *listReviewsArgs) (*Response[resource.List[resource.Review]], error) {
	var rvwMs []model.Review
	var err error

	if args.FirstCreated != 0 {
		rvwMs, err = m.QS.ListLatestApprovedReviewsOfPlaceReverse(
			ctx, args.PlaceID, args.FirstCreated, args.FirstID, args.Limit,
		)
	} else {
		var lastCreated, lastID int64
		if args.LastCreated == 0 {
			lastCreated, lastID = math.MaxInt64, math.MaxInt64
		} else {
			lastCreated, lastID = args.LastCreated, args.LastID
		}

		rvwMs, err = m.QS.ListLatestApprovedReviewsOfPlace(
			ctx, args.PlaceID, lastCreated, lastID, args.Limit,
		)
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &Response[resource.List[resource.Review]]{Body: render.Reviews(rvwMs)}, nil
}
