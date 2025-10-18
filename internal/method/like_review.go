package method

import (
	"context"
	"net/http"

	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/resource"
	"codeberg.org/socialmaps/api/internal/web"
)

type LikeReview struct {
	Common
}

type likeReviewArgs struct {
	ReviewID string `path:"review_id"`
}

func (m *LikeReview) Execute(ctx context.Context, args *likeReviewArgs) *web.Response {
	userID := "usr_foo"

	model.LikeReview(ctx, m.DB, userID, args.ReviewID)

	return web.NewEmptyResponse(http.StatusNoContent)
}

func (m *LikeReview) Validate(args *likeReviewArgs) *web.Response {
	if !model.IsValidReviewID(args.ReviewID) {
		return web.NewJSONResponse(http.StatusBadRequest, resource.Error{
			Message: resource.ErrMsgInvalidReviewID,
		})
	}

	return nil
}

func (m *LikeReview) NewArgs() *likeReviewArgs {
	return &likeReviewArgs{}
}
