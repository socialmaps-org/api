package method

import (
	"context"

	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/render"
	"codeberg.org/socialmaps/api/internal/resource"
	"github.com/danielgtaylor/huma/v2"
)

type RetrievePlace struct {
	Common
}

type retrievePlaceArgs struct {
	PlaceID string `path:"place_id" pattern:"^plc_[a-zA-Z0-9]+$"`
}

func (m *RetrievePlace) Execute(ctx context.Context, args *retrievePlaceArgs) (*Response[*resource.Place], error) {
	plcM := model.LoadPlaceByID(ctx, m.DB, args.PlaceID)
	if plcM == nil {
		return nil, huma.Error404NotFound("place not found")
	}

	plcR := render.Place(plcM)

	return &Response[*resource.Place]{Body: plcR}, nil
}
