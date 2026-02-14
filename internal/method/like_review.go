package method

import (
	"context"
	"database/sql"

	"github.com/danielgtaylor/huma/v2"
)

type LikeReview struct {
	Common
}

type likeReviewArgs struct {
	ReviewID int64 `path:"review_id" minimum:"1"`
}

func (m *LikeReview) Execute(ctx context.Context, args *likeReviewArgs) (*struct{}, error) {
	usr := GetAuthUser(ctx)

	rvw, err := m.QS.LoadReview(ctx, args.ReviewID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, huma.Error404NotFound("review not found")
		}
		return nil, err
	}

	if rvw.UserID == usr.ID {
		// Users cannot like their own reviews
		return nil, huma.Error403Forbidden("cannot like yours")
	}

	err = m.QS.LikeReview(ctx, args.ReviewID, usr.ID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
