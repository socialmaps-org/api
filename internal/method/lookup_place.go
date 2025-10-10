package method

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"codeberg.org/socialmaps/auth/internal/fun"
	"codeberg.org/socialmaps/auth/internal/geo"
	"codeberg.org/socialmaps/auth/internal/model"
	"codeberg.org/socialmaps/auth/internal/name"
	"codeberg.org/socialmaps/auth/internal/overpass"
	"codeberg.org/socialmaps/auth/internal/render"
	"codeberg.org/socialmaps/auth/internal/resource"
	"codeberg.org/socialmaps/auth/internal/web"
)

type LookupPlace struct {
	Common
}

type lookupPlaceArgs struct {
	Name string  `schema:"name,required"`
	Lat  float64 `schema:"lat,required"`
	Lon  float64 `schema:"lon,required"`
}

func (m *LookupPlace) Execute(ctx context.Context, args *lookupPlaceArgs) *web.Response {
	bbox := geo.NewBBox(args.Lat, args.Lon, 10)
	slog.InfoContext(ctx, "bbox",
		"south", bbox.South, "west", bbox.West, "north", bbox.North, "east", bbox.East,
	)

	places, err := model.ListPlacesByCoord(ctx, m.DB, bbox.South, bbox.West, bbox.North, bbox.East)
	if err != nil {
		return web.NewEmptyResponse(http.StatusInternalServerError)
	}

	dbCandidates := fun.Filter(
		places,
		func(plc *model.Place) bool {
			return name.Equivalent(args.Name, plc.Name)
		},
	)

	if len(dbCandidates) > 1 {
		return web.NewEmptyResponse(http.StatusInternalServerError)
	} else if len(dbCandidates) == 1 {
		plcM := dbCandidates[0]
		return web.NewJSONResponse(http.StatusOK, render.Place(plcM))
	}

	slog.InfoContext(ctx, "db results", "results", dbCandidates)

	opRes, err := overpass.Query(
		fmt.Sprintf(
			// OpenStreetMap requires 7 decimal places for geographic coordinates.
			`[out:json];nwr(%.7f, %.7f, %.7f, %.7f)[name];out center tags;`,
			bbox.South, bbox.West, bbox.North, bbox.East,
		),
	)
	if err != nil {
		return web.NewEmptyResponse(http.StatusInternalServerError)
	}

	opCandidates := fun.Filter(
		opRes.Elements,
		func(el overpass.Element) bool {
			return name.Equivalent(args.Name, el.Tags["name"])
		},
	)

	if len(opCandidates) > 1 {
		return web.NewJSONResponse(http.StatusBadRequest, &resource.Error{
			Message: "multiple elements with the same name exist",
		})
	} else if len(opCandidates) == 0 {
		return web.NewEmptyResponse(http.StatusNotFound)
	}

	el := opCandidates[0]
	plcM := model.CreatePlace(ctx, m.DB, el.Tags["name"], el.Lat(), el.Lon(), el.Type, el.ID)
	plcR := render.Place(plcM)

	return web.NewJSONResponse(http.StatusOK, plcR)
}

func (m *LookupPlace) Validate(args *lookupPlaceArgs) *web.Response {
	if !(len(args.Name) <= 256) {
		return web.NewJSONResponse(http.StatusBadRequest, &resource.Error{
			Message: resource.ErrMsgNameTooLong,
		})
	}

	if !(-90 <= args.Lat && args.Lat <= +90) {
		return web.NewJSONResponse(http.StatusBadRequest, &resource.Error{
			Message: resource.ErrMsgLatOutOfRange,
		})
	}

	if !(-180 <= args.Lon && args.Lon <= +180) {
		return web.NewJSONResponse(http.StatusBadRequest, &resource.Error{
			Message: resource.ErrMsgLonOutOfRange,
		})
	}

	return nil
}

func (m *LookupPlace) NewArgs() *lookupPlaceArgs {
	return &lookupPlaceArgs{}
}
