package method

import (
	"context"
	"net/http"

	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/resource"
	"codeberg.org/socialmaps/api/internal/web"
)

type DeleteReview struct {
	Common
}

type deleteReviewArgs struct {
	ReviewID string `path:"review_id"`
}

func (m *DeleteReview) Execute(ctx context.Context, args *deleteReviewArgs) *web.Response {
	usr := web.GetAuthUser(ctx)

	rvw := model.LoadReview(ctx, m.DB, args.ReviewID)
	if rvw == nil {
		return web.NewResponse(http.StatusNotFound, nil)
	} else if rvw.UserID != usr.ID {
		// Users cannot delete others' reviews
		return web.NewResponse(http.StatusForbidden, nil)
	}

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
