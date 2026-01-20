package method

import (
	"context"
	"database/sql"

	"codeberg.org/socialmaps/api/internal/render"
	"codeberg.org/socialmaps/api/internal/resource"
	"github.com/danielgtaylor/huma/v2"
)

type CreateReview struct {
	Common
}

type createReviewArgs struct {
	PlaceID int64 `path:"place_id" minimum:"0"`
	Body    struct {
		Liked   bool   `json:"liked"`
		Comment string `json:"comment"`
	}
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
