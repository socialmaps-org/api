package method

import (
	"context"
	"database/sql"
	"reflect"

	"codeberg.org/socialmaps/api/internal/render"
	"codeberg.org/socialmaps/api/internal/resource"
	"github.com/danielgtaylor/huma/v2"
)

type CreateReview struct {
	Common
}

type createReviewBodyArg struct {
	Liked   bool   `json:"liked" doc:"Whether the user liked this **Place** or not."`
	Comment string `json:"comment" doc:"The comment written by the user about this **Place**, if written. Otherwise can be an empty string."`
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

	plc, err := m.QS.LoadPlace(ctx, args.PlaceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, huma.Error404NotFound("place not found")
		}
		return nil, err
	}

	rvwM, err := m.QS.CreateReview(ctx, plc.ID, usr.ID, args.Body.Liked, sql.NullString{String: args.Body.Comment, Valid: true})
	if err != nil {
		return nil, err
	}

	rvwR := render.Review(rvwM)

	return &Response[resource.Review]{Body: rvwR}, nil
}
