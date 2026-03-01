package method

import (
	"context"
	"errors"
	"log/slog"

	"codeberg.org/socialmaps/api/internal/geo"
	"codeberg.org/socialmaps/api/internal/mytime"
	"codeberg.org/socialmaps/api/internal/render"
	"codeberg.org/socialmaps/api/internal/resource"
	"github.com/jackc/pgx/v5"
)

type LookupPlace struct {
	Common
	NominatimEndpoint string
}

type lookupPlaceArgs struct {
	Name string  `query:"name" required:"true" maxLength:"256" doc:"The name of the **Place** in any language verbatim as recorded in OpenStreetMap."`
	Lat  float64 `query:"lat"  required:"true" minimum:"-180.0" maximum:"+180.0" doc:"Latitude of the **Place** (or its centre if it's a way/relation) in decimal degrees based on [WGS 84](https://en.wikipedia.org/wiki/WGS-84).\n\nThe location doesn't have to be exactly accurate, but must be precise enough (7 decimal places) to narrow down the search."`
	Lon  float64 `query:"lon"  required:"true" minimum:"-90.0"  maximum:"+90.0" doc:"Approximate longitude of the **Place** (or its centre if it's a way/relation) in decimal degrees based on [WGS 84](https://en.wikipedia.org/wiki/WGS-84).\n\nThe location doesn't have to be exactly accurate, but must be precise enough (7 decimal places) to narrow down the search."`
}

func (m *LookupPlace) Execute(ctx context.Context, args *lookupPlaceArgs) (*Response[resource.Place], error) {
	bbox := geo.NewBBox(args.Lat, args.Lon, 50)

	places, err := m.QS.LookupPlaces(ctx, bbox.South, bbox.North, bbox.West, bbox.East, args.Name)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	if len(places) > 1 {
		return nil, errors.New("internal server error")
	} else if len(places) == 1 {
		plcM := places[0]
		return &Response[resource.Place]{Body: render.Place(plcM)}, nil
	}

	elements, err := m.QS.LookupElements(ctx, bbox.South, bbox.North, bbox.West, bbox.East, args.Name)
	if err != nil {
		return nil, err
	}

	if len(elements) > 1 {
		return nil, errors.New("internal server error")
	} else if len(elements) == 0 {
		return nil, errors.New("not found")
	}

	elm := elements[0]

	slog.InfoContext(ctx, "create-place",
		"name", elm.Name,
		"osm_type", elm.OsmType,
		"osm_id", elm.OsmID,
	)

	plcM, err := m.QS.CreatePlace(ctx, *elm.Name, elm.Lon, elm.Lat, elm.OsmType, elm.OsmID, mytime.Now())
	if err != nil {
		return nil, err
	}

	return &Response[resource.Place]{Body: render.Place(plcM)}, nil
}
