package method

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
	"golang.socialmaps.org/api/internal/geo"
	"golang.socialmaps.org/api/internal/model"
	"golang.socialmaps.org/api/internal/mytime"
	"golang.socialmaps.org/api/internal/render"
	"golang.socialmaps.org/api/internal/resource"
)

type LookupPlace struct {
	Common
	NominatimEndpoint string
}

type lookupPlaceArgs struct {
	Name string  `query:"name" required:"true" maxLength:"256" doc:"The name of the **Place** in any language verbatim as recorded in OpenStreetMap." example:"Mausoleum at Halicarnassus"`
	Lat  float64 `query:"lat"  required:"true" minimum:"-90.0" maximum:"+90.0" doc:"Latitude of the **Place** (or its centre if it's a way/relation) in decimal degrees based on [WGS 84](https://en.wikipedia.org/wiki/WGS-84).\n\nThe location doesn't have to be exactly accurate, but must be precise enough (7 decimal places) to narrow down the search." example:"37.0377"`
	Lon  float64 `query:"lon"  required:"true" minimum:"-180.0"  maximum:"+180.0" doc:"Approximate longitude of the **Place** (or its centre if it's a way/relation) in decimal degrees based on [WGS 84](https://en.wikipedia.org/wiki/WGS-84).\n\nThe location doesn't have to be exactly accurate, but must be precise enough (7 decimal places) to narrow down the search." example:"27.4241"`
}

func (m *LookupPlace) Execute(ctx context.Context, args *lookupPlaceArgs) (*Response[resource.Place], error) {
	bbox := geo.NewBBox(args.Lat, args.Lon, 50)

	tuple, err := m.QS.LookupPlace(ctx, args.Name, bbox.West, bbox.South, bbox.East, bbox.North)

	if err == pgx.ErrNoRows {
		return nil, huma.Error404NotFound("place not found")
	}

	elm := tuple.Element
	var plcM model.ComputedPlace

	if tuple.OptionalComputedPlace.IsNil() {
		plc, err := m.QS.CreatePlace(
			ctx,
			elm.Name,
			elm.Lon,
			elm.Lat,
			elm.OsmType,
			elm.OsmID,
			mytime.Now(),
		)
		if err != nil {
			return nil, err
		}
		plcM = model.NewPlaceToComputedPlace(plc)
	} else {
		plcM = tuple.OptionalComputedPlace.Unwrap()
	}
	return &Response[resource.Place]{Body: render.Place(plcM, elm)}, nil
}
