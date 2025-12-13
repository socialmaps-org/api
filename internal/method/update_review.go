package method

import (
	"context"
	"log/slog"
	"time"

	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/render"
	"codeberg.org/socialmaps/api/internal/resource"
	"github.com/danielgtaylor/huma/v2"
)

type UpdateReview struct {
	Common
}

type updateReviewArgs struct {
	ReviewID string `path:"review_id" pattern:"^rvw_[a-zA-Z0-9]+$"`
	Body     struct {
		Liked   bool   `json:"liked"`
		Comment string `json:"comment"`
	}
}

// 1 hour
const maxDelayInSec = 1 * 60 * 60

func (m *UpdateReview) Execute(ctx context.Context, args *updateReviewArgs) (*Response[*resource.Review], error) {
	usr := GetAuthUser(ctx)

	rvwM := model.LoadReview(ctx, m.DB, args.ReviewID)
	if rvwM == nil {
		return nil, huma.Error404NotFound("review not found")
	}

	if rvwM.UserID != usr.ID {
		slog.InfoContext(ctx, "CANONICAL-METHOD-LINE", "review_user", rvwM.UserID, "auth_user", usr.ID)
		return nil, huma.Error403Forbidden("not your review")
	}

	if time.Now().Unix()-rvwM.Created > maxDelayInSec {
		return nil, huma.Error400BadRequest("too late")
	}

	rvwM = model.UpdateReview(ctx, m.DB, args.ReviewID, args.Body.Liked, args.Body.Comment)

	rvwR := render.Review(rvwM)

	return &Response[*resource.Review]{Body: rvwR}, nil
}
