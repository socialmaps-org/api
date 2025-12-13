package method

import (
	"context"

	"codeberg.org/socialmaps/api/internal/model"
	"github.com/danielgtaylor/huma/v2"
)

type LikeReview struct {
	Common
}

type likeReviewArgs struct {
	ReviewID string `path:"review_id" pattern:"^rvw_[a-zA-Z0-9]+$"`
}

func (m *LikeReview) Execute(ctx context.Context, args *likeReviewArgs) (*struct{}, error) {
	usr := GetAuthUser(ctx)

	rvw := model.LoadReview(ctx, m.DB, args.ReviewID)
	if rvw == nil {
		return nil, huma.Error404NotFound("review not found")
	} else if rvw.UserID == usr.ID {
		// Users cannot like their own reviews
		return nil, huma.Error403Forbidden("cannot like yours")
	}

	model.LikeReview(ctx, m.DB, args.ReviewID, usr.ID)
	return nil, nil
}
