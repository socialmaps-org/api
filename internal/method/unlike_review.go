package method

import (
	"context"
	"net/http"

	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/resource"
	"codeberg.org/socialmaps/api/internal/web"
)

type UnlikeReview struct {
	Common
}

type unlikeReviewArgs struct {
	ReviewID string `path:"review_id"`
}

func (m *UnlikeReview) Execute(ctx context.Context, args *unlikeReviewArgs) *web.Response {
	userID := "usr_foo"

	model.UnlikeReview(ctx, m.DB, userID, args.ReviewID)

	return web.NewEmptyResponse(http.StatusNoContent)
}

func (m *UnlikeReview) Validate(args *unlikeReviewArgs) *web.Response {
	if !model.IsValidReviewID(args.ReviewID) {
		return web.NewJSONResponse(http.StatusBadRequest, resource.Error{
			Message: resource.ErrMsgInvalidReviewID,
		})
	}

	return nil
}

func (m *UnlikeReview) NewArgs() *unlikeReviewArgs {
	return &unlikeReviewArgs{}
}
