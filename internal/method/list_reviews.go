package method

import (
	"context"
	"net/http"

	"codeberg.org/socialmaps/auth/internal/model"
	"codeberg.org/socialmaps/auth/internal/render"
	"codeberg.org/socialmaps/auth/internal/resource"
	"codeberg.org/socialmaps/auth/internal/web"
)

type ListReviews struct {
	Common
}

type listReviewsArgs struct {
	PlaceID string `path:"place_id"`
}

func (m *ListReviews) Execute(ctx context.Context, args *listReviewsArgs) *web.Response {
	rvwMs := model.ListLatestReviewsOfPlace(ctx, m.DB, args.PlaceID)

	return web.NewJSONResponse(http.StatusOK, render.Reviews(rvwMs))
}

func (m *ListReviews) Validate(args *listReviewsArgs) *web.Response {
	if !model.IsValidPlaceID(args.PlaceID) {
		return web.NewJSONResponse(http.StatusBadRequest, &resource.Error{
			Message: resource.ErrMsgInvalidPlaceID,
		})
	}

	return nil
}

func (m *ListReviews) NewArgs() *listReviewsArgs {
	return &listReviewsArgs{}
}
