package method

import (
	"context"

	"codeberg.org/socialmaps/api/internal/model"
	"github.com/danielgtaylor/huma/v2"
)

type UnlikeReview struct {
	Common
}

type unlikeReviewArgs struct {
	ReviewID string `path:"review_id" pattern:"^rvw_[a-zA-Z0-9]+$"`
}

func (m *UnlikeReview) Execute(ctx context.Context, args *unlikeReviewArgs) (*struct{}, error) {
	usr := GetAuthUser(ctx)

	rvw := model.LoadReview(ctx, m.DB, args.ReviewID)
	if rvw == nil {
		return nil, huma.Error404NotFound("review not found")
	}

	model.UnlikeReview(ctx, m.DB, args.ReviewID, usr.ID)

	return nil, nil
}
