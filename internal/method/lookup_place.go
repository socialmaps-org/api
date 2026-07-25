package method

import (
	"context"
	"errors"
	"fmt"

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
	Name string  `query:"name" required:"true" maxLength:"256" doc:"The name of the **Place** in any language verbatim as recorded in OpenStreetMap."`
	Lat  float64 `query:"lat"  required:"true" minimum:"-180.0" maximum:"+180.0" doc:"Latitude of the **Place** (or its centre if it's a way/relation) in decimal degrees based on [WGS 84](https://en.wikipedia.org/wiki/WGS-84).\n\nThe location doesn't have to be exactly accurate, but must be precise enough (7 decimal places) to narrow down the search."`
	Lon  float64 `query:"lon"  required:"true" minimum:"-90.0"  maximum:"+90.0" doc:"Approximate longitude of the **Place** (or its centre if it's a way/relation) in decimal degrees based on [WGS 84](https://en.wikipedia.org/wiki/WGS-84).\n\nThe location doesn't have to be exactly accurate, but must be precise enough (7 decimal places) to narrow down the search."`
}

func (m *LookupPlace) Execute(ctx context.Context, args *lookupPlaceArgs) (*Response[resource.Place], error) {
	bbox := geo.NewBBox(args.Lat, args.Lon, 50)

	tuples, err := m.QS.QueryPlaces(ctx, bbox.South, bbox.North, bbox.West, bbox.East, fmt.Sprintf(`$.name == "%s"`, args.Name))
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	if len(tuples) > 1 {
		return nil, errors.New("internal server error")
	}

	tuple := tuples[0]
	elm := tuple.Element
	var plcM model.Place

	if tuple.OptionalPlace.IsNil() {
		plcM, err = m.QS.CreatePlace(
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
	} else {
		plcM = tuple.OptionalPlace.Unwrap()
	}
	return &Response[resource.Place]{Body: render.Place(plcM, elm)}, nil
}
