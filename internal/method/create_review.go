package method

import (
	"context"

	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/render"
	"codeberg.org/socialmaps/api/internal/resource"
	"github.com/danielgtaylor/huma/v2"
)

type CreateReview struct {
	Common
}

type createReviewArgs struct {
	PlaceID string `path:"place_id" pattern:"^plc_[a-zA-Z0-9]+$"`
	Body    struct {
		Liked   bool   `json:"liked"`
		Comment string `json:"comment"`
	}
}

func (m *CreateReview) Execute(ctx context.Context, args *createReviewArgs) (*Response[*resource.Review], error) {
	usr := GetAuthUser(ctx)

	plc := model.LoadPlaceByID(ctx, m.DB, args.PlaceID)
	if plc == nil {
		return nil, huma.Error404NotFound("place not found")
	}

	rvwM := model.CreateReview(ctx, m.DB, args.PlaceID, usr.ID, args.Body.Liked, args.Body.Comment)

	rvwR := render.Review(rvwM)

	return &Response[*resource.Review]{Body: rvwR}, nil
}
