package method

import (
	"context"
	"database/sql"

	"codeberg.org/socialmaps/api/internal/render"
	"codeberg.org/socialmaps/api/internal/resource"
	"github.com/danielgtaylor/huma/v2"
)

type RetrievePlace struct {
	Common
}

type retrievePlaceArgs struct {
	PlaceID int64 `path:"place_id" minimum:"1"`
}

func (m *RetrievePlace) Execute(ctx context.Context, args *retrievePlaceArgs) (*Response[resource.Place], error) {
	plcM, err := m.QS.LoadPlace(ctx, args.PlaceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, huma.Error404NotFound("place not found")
		}
		return nil, err
	}

	plcR := render.Place(plcM)

	return &Response[resource.Place]{Body: plcR}, nil
}
