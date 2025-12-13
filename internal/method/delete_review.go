package method

import (
	"context"

	"codeberg.org/socialmaps/api/internal/model"
	"github.com/danielgtaylor/huma/v2"
)

type DeleteReview struct {
	Common
}

type deleteReviewArgs struct {
	ReviewID string `path:"review_id" pattern:"^rvw_[a-zA-Z0-9]+$"`
}

func (m *DeleteReview) Execute(ctx context.Context, args *deleteReviewArgs) (*struct{}, error) {
	usr := GetAuthUser(ctx)

	rvw := model.LoadReview(ctx, m.DB, args.ReviewID)
	if rvw == nil {
		return nil, huma.Error404NotFound("review not found")
	} else if rvw.UserID != usr.ID {
		// Users cannot delete others' reviews
		return nil, huma.Error403Forbidden("not your review")
	}

	model.DeleteReview(ctx, m.DB, args.ReviewID)
	return nil, nil
}
