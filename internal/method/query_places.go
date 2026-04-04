package method

import (
	"context"

	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/mytime"
	"codeberg.org/socialmaps/api/internal/render"
	"codeberg.org/socialmaps/api/internal/resource"
)

type QueryPlaces struct {
	Common
}

type queryPlacesArgs struct {
	MaxLat    float64 `query:"max_lat" required:"true" minimum:"-180.0" maximum:"+180.0" doc:"Maximum latitude (i.e. north-side) of the bounding-box in which **Places** are."`
	MaxLon    float64 `query:"max_lon" required:"true" minimum:"-90.0"  maximum:"+90.0" doc:"Maximum longitude (i.e. east-side) of the bounding-box in which **Places** are."`
	MinLat    float64 `query:"min_lat" required:"true" minimum:"-180.0" maximum:"+180.0" doc:"Minimum latitude (i.e. south-side) of the bounding-box in which **Places** are."`
	MinLon    float64 `query:"min_lon" required:"true" minimum:"-90.0"  maximum:"+90.0" doc:"Minimum longitude (i.e. west-side) of the bounding-box in which **Places** are."`
	Predicate string  `query:"predicate" required:"true" doc:"A PostgreSQL-compatible SQL-standard SQL/JSON Path expression to filter **Places** in the bounding-box by their OpenStreetMap [tags](https://wiki.openstreetmap.org/wiki/Tags)." example:"$.amenity == \"restaurant\" && $.cuisine like_regex \"turkish\" && $.outdoor_seating == \"yes\""`
}

func (m *QueryPlaces) Execute(ctx context.Context, args *queryPlacesArgs) (*Response[resource.List[resource.Place]], error) {
	tuples, err := m.QS.QueryPlaces(ctx, args.MinLon, args.MinLat, args.MaxLon, args.MaxLat, args.Predicate)
	if err != nil {
		return nil, err
	}

	plcRs := make([]resource.Place, len(tuples))
	for i, tuple := range tuples {
		elm := tuple.Element
		var plcM model.Place

		if tuple.OptionalPlace.IsNil() {
			plcM, err = m.QS.CreatePlace(
				ctx,
				*elm.Name,
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

		plcRs[i] = render.Place(plcM, elm)
	}

	return &Response[resource.List[resource.Place]]{
		Body: resource.List[resource.Place]{
			Object: "list",
			Data:   plcRs,
		},
	}, nil
}
