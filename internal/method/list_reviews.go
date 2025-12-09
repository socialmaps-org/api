package method

import (
	"context"
	"net/http"

	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/render"
	"codeberg.org/socialmaps/api/internal/resource"
	"codeberg.org/socialmaps/api/internal/web"
)

type ListReviews struct {
	Common
}

type listReviewsArgs struct {
	PlaceID       string `path:"place_id"`
	Limit         uint   `schema:"limit,default:10"`
	StartingAfter string `schema:"starting_after"`
	EndingBefore  string `schema:"ending_before"`
}

func (m *ListReviews) Execute(ctx context.Context, args *listReviewsArgs) *web.Response {
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

	return web.NewJSONResponse(http.StatusOK, render.Reviews(rvwMs))
}

func (m *ListReviews) Validate(args *listReviewsArgs) *web.Response {
	if !model.IsValidPlaceID(args.PlaceID) {
		return web.NewJSONResponse(http.StatusBadRequest, &resource.Error{
			Message: resource.ErrMsgInvalidPlaceID,
		})
	}

	if args.Limit > 100 {
		return web.NewJSONResponse(http.StatusBadRequest, &resource.Error{
			Message: resource.ErrMsgLimitTooBig,
		})
	} else if args.Limit == 0 {
		return web.NewJSONResponse(http.StatusBadRequest, &resource.Error{
			Message: resource.ErrMsgLimitZero,
		})
	}

	if args.EndingBefore != "" && args.StartingAfter != "" {
		return web.NewJSONResponse(http.StatusBadRequest, &resource.Error{
			Message: resource.ErrMsgBeforeAfterBothPresent,
		})
	}

	return nil
}

func (m *ListReviews) NewArgs() *listReviewsArgs {
	return &listReviewsArgs{}
}
