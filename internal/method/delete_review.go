package method

import (
	"context"
	"database/sql"

	"github.com/danielgtaylor/huma/v2"
)

type DeleteReview struct {
	Common
}

type deleteReviewArgs struct {
	ReviewID int64 `path:"review_id" minimum:"0"`
}

func (m *DeleteReview) Execute(ctx context.Context, args *deleteReviewArgs) (*struct{}, error) {
	usr := GetAuthUser(ctx)

	rvw, err := m.QS.LoadReview(ctx, args.ReviewID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, huma.Error404NotFound("review not found")
		}
		return nil, err
	}

	if rvw.UserID != usr.ID {
		// Users cannot delete others' reviews
		return nil, huma.Error403Forbidden("not your review")
	}

	err = m.QS.DeleteReview(ctx, args.ReviewID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
