package method

import (
	"context"

	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/render"
	"codeberg.org/socialmaps/api/internal/resource"
)

type ListReviews struct {
	Common
}

type listReviewsArgs struct {
	PlaceID       string `path:"place_id" pattern:"^plc_[a-zA-Z0-9]+$"`
	Limit         uint   `query:"limit" minimum:"1" maximum:"100" default:"20"`
	StartingAfter string `query:"starting_after"`
	EndingBefore  string `query:"ending_before"`
}

func (m *ListReviews) Execute(ctx context.Context, args *listReviewsArgs) (*Response[*resource.List[*resource.Review]], error) {
	var rvwMs []*model.Review

	if args.EndingBefore != "" {
		rvwMs = model.ListLatestApprovedReviewsOfPlaceReverse(
			ctx, m.DB, args.PlaceID, args.Limit, args.EndingBefore,
		)
	} else {
		next := args.StartingAfter
		if next == "" {
			next = model.LatestID("rvw").String()
		}

		rvwMs = model.ListLatestApprovedReviewsOfPlace(
			ctx, m.DB, args.PlaceID, args.Limit, next,
		)
	}

	return &Response[*resource.List[*resource.Review]]{Body: render.Reviews(rvwMs)}, nil
}
