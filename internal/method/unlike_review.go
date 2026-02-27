package method

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
)

type UnlikeReview struct {
	Common
}

type unlikeReviewArgs struct {
	ReviewID int64 `path:"review_id" minimum:"1"`
}

func (m *UnlikeReview) Execute(ctx context.Context, args *unlikeReviewArgs) (*struct{}, error) {
	usr := GetAuthUser(ctx)

	rvw, err := m.QS.LoadReview(ctx, args.ReviewID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, huma.Error404NotFound("review not found")
		}
		return nil, err
	}

	err = m.QS.UnlikeReview(ctx, rvw.ID, usr.ID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
