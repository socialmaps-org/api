package method

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
	"golang.socialmaps.org/api/internal/render"
	"golang.socialmaps.org/api/internal/resource"
)

type RetrievePlace struct {
	Common
}

type retrievePlaceArgs struct {
	PlaceID int64 `path:"place_id" minimum:"1" doc:"Unique identifier for the **Place** you are retrieving."`
}

func (m *RetrievePlace) Execute(ctx context.Context, args *retrievePlaceArgs) (*Response[resource.Place], error) {
	tuple, err := m.QS.LoadPlace(ctx, args.PlaceID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, huma.Error404NotFound("place not found")
		}
		return nil, err
	}

	plcR := render.Place(tuple.ComputedPlace, tuple.Element)

	return &Response[resource.Place]{Body: plcR}, nil
}
