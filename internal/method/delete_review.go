package method

import (
	"context"
	"net/http"

	"codeberg.org/socialmaps/auth/internal/model"
	"codeberg.org/socialmaps/auth/internal/resource"
	"codeberg.org/socialmaps/auth/internal/web"
)

type DeleteReview struct {
	Common
}

type deleteReviewArgs struct {
	ReviewID string `path:"review_id"`
}

func (m *DeleteReview) Execute(ctx context.Context, args *deleteReviewArgs) *web.Response {
	model.DeleteReview(ctx, m.DB, args.ReviewID)
	return web.NewResponse(http.StatusNoContent, nil)
}

func (m *DeleteReview) Validate(args *deleteReviewArgs) *web.Response {
	if !model.IsValidReviewID(args.ReviewID) {
		return web.NewJSONResponse(http.StatusBadRequest, resource.Error{
			Message: resource.ErrMsgInvalidReviewID,
		})
	}

	return nil
}

func (m *DeleteReview) NewArgs() *deleteReviewArgs {
	return &deleteReviewArgs{}
}
