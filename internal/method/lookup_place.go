package method

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"codeberg.org/socialmaps/api/internal/fun"
	"codeberg.org/socialmaps/api/internal/geo"
	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/name"
	"codeberg.org/socialmaps/api/internal/nominatim"
	"codeberg.org/socialmaps/api/internal/render"
	"codeberg.org/socialmaps/api/internal/resource"
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
	slog.InfoContext(ctx, "bbox",
		"south", bbox.South, "west", bbox.West, "north", bbox.North, "east", bbox.East,
	)

	places, err := m.QS.ListPlacesByCoord(ctx, bbox.South, bbox.North, bbox.West, bbox.East)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	dbCandidates := fun.Filter(
		places,
		func(plc model.Place) bool {
			// TODO: we should save all names and check for equivalency against
			// each of them like we do for nominatim results to support localisation
			return name.Equivalent(args.Name, plc.Name)
		},
	)

	slog.InfoContext(ctx, "db results", "results", dbCandidates)

	if len(dbCandidates) > 1 {
		return nil, errors.New("internal server error")
	} else if len(dbCandidates) == 1 {
		plcM := dbCandidates[0]
		return &Response[resource.Place]{Body: render.Place(plcM)}, nil
	}

	nomPlaces, err := nominatim.Search(ctx, m.NominatimEndpoint, args.Name, bbox)
	if err != nil {
		return nil, err
	}

	nomCandidates := fun.Filter(
		nomPlaces,
		func(np nominatim.Place) bool {
			for _, nam := range np.Names {
				if name.Equivalent(nam, args.Name) {
					return true
				}
			}
			return false
		},
	)

	if len(nomCandidates) > 1 {
		return nil, errors.New("internal server error")
	} else if len(nomCandidates) == 0 {
		return nil, errors.New("not found")
	}

	plcN := nomCandidates[0]

	plcM, err := m.QS.CreatePlace(ctx, plcN.Name, plcN.Lat, plcN.Lon, plcN.Type, plcN.ID)
	if err != nil {
		return nil, err
	}

	return &Response[resource.Place]{Body: render.Place(plcM)}, nil
}
