package method

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"codeberg.org/socialmaps/api/internal/render"
	"codeberg.org/socialmaps/api/internal/resource"
	"github.com/danielgtaylor/huma/v2"
)

type UpdateReview struct {
	Common
}

type updateReviewArgs struct {
	ReviewID int64 `path:"review_id" minimum:"0"`
	Body     struct {
		Liked   bool    `json:"liked"`
		Comment *string `json:"comment"`
	}
}

// 1 hour
const maxDelayInSec = 1 * 60 * 60

func (m *UpdateReview) Execute(ctx context.Context, args *updateReviewArgs) (*Response[resource.Review], error) {
	usr := GetAuthUser(ctx)

	rvwM, err := m.QS.LoadReview(ctx, args.ReviewID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, huma.Error404NotFound("review not found")
		}
		return nil, err
	}

	if rvwM.UserID != usr.ID {
		slog.InfoContext(ctx, "CANONICAL-METHOD-LINE", "review_user", rvwM.UserID, "auth_user", usr.ID)
		return nil, huma.Error403Forbidden("not your review")
	}

	if time.Now().Unix()-rvwM.Created > maxDelayInSec {
		return nil, huma.Error400BadRequest("too late")
	}

	var comment sql.NullString
	if args.Body.Comment != nil {
		comment = sql.NullString{String: *args.Body.Comment, Valid: true}
	} else {
		comment = sql.NullString{Valid: false}
	}

	rvwM, err = m.QS.UpdateReview(ctx, args.Body.Liked, comment, rvwM.ID)
	if err != nil {
		return nil, err
	}

	rvwR := render.Review(rvwM)

	return &Response[resource.Review]{Body: rvwR}, nil
}
