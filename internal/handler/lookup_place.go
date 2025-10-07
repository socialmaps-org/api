package handler

import (
	"encoding/json"
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
	Handler
}

type lookupPlaceArgs struct {
	Name string  `schema:"name,required"`
	Lat  float64 `schema:"lat,required"`
	Lon  float64 `schema:"lon,required"`
}

func (h *LookupPlace) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var args lookupPlaceArgs
	if err := web.Parse(r, &args); err != nil {
		web.JSON(w, http.StatusBadRequest, err)
		return
	}
	if err := h.validate(args); err != nil {
		web.JSON(w, err.StatusCode, err)
		return
	}

	bbox := geo.NewBBox(args.Lat, args.Lon, 10)
	slog.InfoContext(r.Context(), "bbox",
		"south", bbox.South, "west", bbox.West, "north", bbox.North, "east", bbox.East,
	)

	places, err := model.ListPlacesByCoord(r.Context(), h.DB, bbox.South, bbox.West, bbox.North, bbox.East)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	dbCandidates := fun.Filter(
		places,
		func(plc *model.Place) bool {
			return name.Equivalent(args.Name, plc.Name)
		},
	)

	if len(dbCandidates) > 1 {
		w.WriteHeader(http.StatusInternalServerError)
		panic("multiple candidates in the database")
	} else if len(dbCandidates) == 1 {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(render.Place(dbCandidates[0]))
		return
	}

	slog.InfoContext(r.Context(), "db results", "results", dbCandidates)

	opRes, err := overpass.Query(
		fmt.Sprintf(
			`[out:json];nwr(%f, %f, %f, %f)[name];out center tags;`,
			bbox.South, bbox.West, bbox.North, bbox.East,
		),
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	opCandidates := fun.Filter(
		opRes.Elements,
		func(el overpass.Element) bool {
			return name.Equivalent(args.Name, el.Tags["name"])
		},
	)

	if len(opCandidates) > 1 {
		web.JSON(w, http.StatusBadRequest, &resource.Error{
			Message: "multiple elements with the same name exist",
		})
		return
	} else if len(opCandidates) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	el := opCandidates[0]
	plcM := model.CreatePlace(r.Context(), h.DB, el.Tags["name"], el.Lat(), el.Lon(), el.Type, el.ID)
	plcR := render.Place(plcM)
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(plcR)
	if err != nil {
		panic(err)
	}
}

func (h *LookupPlace) validate(args lookupPlaceArgs) *web.Error {
	if !(len(args.Name) <= 256) {
		return &web.Error{
			StatusCode: http.StatusUnprocessableEntity,
			Resource: resource.Error{
				Message: "`name` must be shorter than 257 characters",
			},
		}
	}

	if !(-90 <= args.Lat && args.Lat <= +90) {
		return &web.Error{
			StatusCode: http.StatusUnprocessableEntity,
			Resource: resource.Error{
				Message: "`lat` must be between -90 and +90 (both inclusive)",
			},
		}
	}

	if !(-180 <= args.Lon && args.Lon <= +180) {
		return &web.Error{
			StatusCode: http.StatusUnprocessableEntity,
			Resource: resource.Error{
				Message: "`lon` must be between -180 and +180 (both inclusive)",
			},
		}
	}

	return nil
}
