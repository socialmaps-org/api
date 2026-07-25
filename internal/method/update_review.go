package method

import (
	"context"
	"log/slog"
	"net/http"
	"reflect"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
	"golang.socialmaps.org/api/internal/mytime"
	"golang.socialmaps.org/api/internal/render"
	"golang.socialmaps.org/api/internal/resource"
)

type UpdateReview struct {
	Common
}

type updateReviewBodyArg struct {
	Liked   bool   `json:"liked" doc:"Whether the user liked this **Place** or not."`
	Comment string `json:"comment" doc:"The comment written by the user about this **Place**, if written. Otherwise can be an empty string."`
}

type updateReviewArgs struct {
	ReviewID int64 `path:"review_id" minimum:"1" doc:"Unique identifier for the **Review** the user is updating."`
	Body     updateReviewBodyArg
}

// Inline updateReviewBodyArg instead of using refs to avoid listing it under Schemas.
func (updateReviewBodyArg) Schema(r huma.Registry) *huma.Schema {
	type raw updateReviewBodyArg
	return huma.SchemaFromType(r, reflect.TypeOf(raw{}))
}

const maxDelay = 1 * time.Hour

func (m *UpdateReview) Execute(ctx context.Context, args *updateReviewArgs) (*DynamicResponse[resource.Review], error) {
	usr := GetAuthUser(ctx)

	now := mytime.Now()

	rvwM, err := m.QS.LoadReview(ctx, args.ReviewID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, huma.Error404NotFound("review not found")
		}
		return nil, err
	}

	if rvwM.UserID != usr.ID {
		slog.InfoContext(ctx, "CANONICAL-METHOD-LINE", "review_user", rvwM.UserID, "auth_user", usr.ID)
		return nil, huma.Error403Forbidden("not your review")
	}

	if now.Sub(rvwM.Created) > maxDelay {
		return nil, huma.Error400BadRequest("too late")
	}

	var status int
	if (args.Body.Comment == "" && rvwM.Comment == nil) || *rvwM.Comment == args.Body.Comment {
		status = http.StatusOK
	} else {
		status = http.StatusAccepted
	}

	var comment *string
	if args.Body.Comment != "" {
		comment = &args.Body.Comment
	} else {
		comment = nil
	}

	rvwM, err = m.QS.UpdateReview(ctx, args.Body.Liked, comment, rvwM.ID, now)
	if err != nil {
		return nil, err
	}

	rvwR := render.Review(rvwM)

	return &DynamicResponse[resource.Review]{Body: rvwR, Status: status}, nil
}
