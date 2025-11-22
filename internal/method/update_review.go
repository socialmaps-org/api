package method

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/render"
	"codeberg.org/socialmaps/api/internal/resource"
	"codeberg.org/socialmaps/api/internal/web"
)

type UpdateReview struct {
	Common
}

type updateReviewArgs struct {
	ReviewID string `path:"review_id"`
	Liked    bool   `json:"liked"`
	Comment  string `json:"comment"`
}

// 1 hour
const maxDelayInSec = 1 * 60 * 60

func (m *UpdateReview) Execute(ctx context.Context, args *updateReviewArgs) *web.Response {
	usr := web.GetAuthUser(ctx)

	rvwM := model.LoadReview(ctx, m.DB, args.ReviewID)
	if rvwM == nil {
		return web.NewEmptyResponse(http.StatusNotFound)
	}

	if rvwM.UserID != usr.ID {
		slog.InfoContext(ctx, "CANONICAL-METHOD-LINE", "review_user", rvwM.UserID, "auth_user", usr.ID)
		return web.NewEmptyResponse(http.StatusForbidden)
	}

	if time.Now().Unix()-rvwM.Created > maxDelayInSec {
		return web.NewJSONResponse(http.StatusBadRequest, &resource.Error{
			Code:    resource.ErrorCodeTooLate,
			Message: "too late to update your review",
		})
	}

	rvwM = model.UpdateReview(ctx, m.DB, args.ReviewID, args.Liked, args.Comment)

	rvwR := render.Review(rvwM)

	return web.NewJSONResponse(http.StatusOK, rvwR)
}

func (m *UpdateReview) Validate(args *updateReviewArgs) *web.Response {
	if !model.IsValidReviewID(args.ReviewID) {
		return web.NewJSONResponse(http.StatusBadRequest, resource.Error{
			Message: resource.ErrMsgInvalidReviewID,
		})
	}

	return nil
}

func (m *UpdateReview) NewArgs() *updateReviewArgs {
	return &updateReviewArgs{}
}
