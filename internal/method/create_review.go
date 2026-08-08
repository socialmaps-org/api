package method

import (
	"context"
	"reflect"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
	"golang.socialmaps.org/api/internal/mytime"
	"golang.socialmaps.org/api/internal/render"
	"golang.socialmaps.org/api/internal/resource"
)

type CreateReview struct {
	Common
}

type createReviewBodyArg struct {
	Liked   bool   `json:"liked" doc:"Whether the user liked this **Place** or not." example:"true"`
	Comment string `json:"comment" doc:"The comment written by the user about this **Place**, if written. Otherwise can be an empty string." example:"It’s one of the Seven Wonders of the Ancient World!"`

	ReviewedAt *int64 `json:"reviewed_at,omitzero" doc:"The [UNIX timestamp](https://en.wikipedia.org/wiki/Unix_time) of when the **Place** was originally reviewed at, if different from now (such as while importing **Review**s from another platform). This cannot be in the future."`
}

type createReviewArgs struct {
	PlaceID int64 `path:"place_id" minimum:"0" doc:"Unique identifier for the **Place** the user is creating a **Review** for."`
	Body    createReviewBodyArg
}

// Inline createReviewBodyArg instead of using refs to avoid listing it under Schemas.
func (createReviewBodyArg) Schema(r huma.Registry) *huma.Schema {
	type raw createReviewBodyArg
	return huma.SchemaFromType(r, reflect.TypeOf(raw{}))
}

func (m *CreateReview) Execute(ctx context.Context, args *createReviewArgs) (*Response[resource.Review], error) {
	usr := GetAuthUser(ctx)

	tuple, err := m.QS.LoadPlace(ctx, args.PlaceID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, huma.Error404NotFound("place not found")
		}
		return nil, err
	}
	plc := tuple.Place

	now := mytime.Now()
	var reviewedAt time.Time
	if args.Body.ReviewedAt != nil {
		reviewedAt = time.Unix(*args.Body.ReviewedAt, 0)
	} else {
		reviewedAt = now
	}

	if reviewedAt.After(now) {
		return nil, huma.Error400BadRequest(
			"reviewed_at cannot be in the future",
			&huma.ErrorDetail{
				Message:  "reviewed_at cannot be in the future",
				Location: "body.reviewed_at",
				Value:    args.Body.ReviewedAt,
			},
		)
	}

	rvwM, err := m.QS.CreateReview(ctx, plc.ID, usr.ID, args.Body.Liked, &args.Body.Comment, now, reviewedAt)
	if err != nil {
		return nil, err
	}

	rvwR := render.Review(rvwM)

	return &Response[resource.Review]{Body: rvwR}, nil
}
