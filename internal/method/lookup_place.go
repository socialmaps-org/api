package method

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"codeberg.org/socialmaps/api/internal/fun"
	"codeberg.org/socialmaps/api/internal/geo"
	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/name"
	"codeberg.org/socialmaps/api/internal/overpass"
	"codeberg.org/socialmaps/api/internal/render"
	"codeberg.org/socialmaps/api/internal/resource"
)

type LookupPlace struct {
	Common
	OverpassEndpoint string
}

type lookupPlaceArgs struct {
	Name string  `query:"name" required:"true" maxLength:"256"`
	Lat  float64 `query:"lat"  required:"true" minimum:"-180.0" maximum:"+180.0"`
	Lon  float64 `query:"lon"  required:"true" minimum:"-90.0"  maximum:"+90.0"`
}

func (m *LookupPlace) Execute(ctx context.Context, args *lookupPlaceArgs) (*Response[*resource.Place], error) {
	bbox := geo.NewBBox(args.Lat, args.Lon, 10)
	slog.InfoContext(ctx, "bbox",
		"south", bbox.South, "west", bbox.West, "north", bbox.North, "east", bbox.East,
	)

	places, err := model.ListPlacesByCoord(ctx, m.DB, bbox.South, bbox.West, bbox.North, bbox.East)
	if err != nil {
		return nil, err
	}

	dbCandidates := fun.Filter(
		places,
		func(plc *model.Place) bool {
			return name.Equivalent(args.Name, plc.Name)
		},
	)

	if len(dbCandidates) > 1 {
		return nil, errors.New("internal server error")
	} else if len(dbCandidates) == 1 {
		plcM := dbCandidates[0]
		return &Response[*resource.Place]{Body: render.Place(plcM)}, nil
	}

	slog.InfoContext(ctx, "db results", "results", dbCandidates)

	opRes, err := overpass.Query(
		m.OverpassEndpoint,
		fmt.Sprintf(
			// OpenStreetMap requires 7 decimal places for geographic coordinates.
			`[out:json];nwr(%.7f, %.7f, %.7f, %.7f)[name];out center tags;`,
			bbox.South, bbox.West, bbox.North, bbox.East,
		),
	)
	if err != nil {
		return nil, err
	}

	opCandidates := fun.Filter(
		opRes.Elements,
		func(el overpass.Element) bool {
			return name.Equivalent(args.Name, el.Tags["name"])
		},
	)

	if len(opCandidates) > 1 {
		return nil, errors.New("multiple elements with the same name exist")
	} else if len(opCandidates) == 0 {
		return nil, errors.New("not found")
	}

	el := opCandidates[0]
	plcM := model.CreatePlace(ctx, m.DB, el.Tags["name"], el.Lat(), el.Lon(), el.Type, el.ID)

	return &Response[*resource.Place]{Body: render.Place(plcM)}, nil
}
