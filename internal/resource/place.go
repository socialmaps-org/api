package resource

import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

type Location struct {
	Lat float64 `json:"lat" minimum:"-180.0" maximum:"+180.0" doc:"Latitude in decimal degrees." example:"37.03778174462128"`
	Lon float64 `json:"lon" minimum:"-90.0" maximum:"+90.0" doc:"Latitude in decimal degrees." example:"27.42417151281665"`
}

// Inline Location instead of using refs to avoid listing it under Schemas.
func (Location) Schema(r huma.Registry) *huma.Schema {
	type raw Location
	return huma.SchemaFromType(r, reflect.TypeOf(raw{}))
}

type RatingSummary struct {
	Average *float64 `json:"average" minimum:"1" maximum:"5" doc:"Average **rating** of all **Reviews** created for this **Place**." example:"4.8"`
}

// Inline RatingSummary instead of using refs to avoid listing it under Schemas.
func (RatingSummary) Schema(r huma.Registry) *huma.Schema {
	type raw RatingSummary
	return huma.SchemaFromType(r, reflect.TypeOf(raw{}))
}

type ReviewSummary struct {
	Count  int64         `json:"count" minimum:"0" doc:"Number of **Review**s created for this **Place**." example:"24"`
	Rating RatingSummary `json:"rating"`
}

// Inline ReviewSummary instead of using refs to avoid listing it under Schemas.
func (ReviewSummary) Schema(r huma.Registry) *huma.Schema {
	type raw ReviewSummary
	return huma.SchemaFromType(r, reflect.TypeOf(raw{}))
}

type Place struct {
	Object        string            `json:"object,omitempty" enum:"place"`
	ID            int64             `json:"id" minimum:"1" doc:"Unique identifier for this **Place**." example:"1"`
	Name          string            `json:"name,omitempty" maxLength:"256" doc:"The common, default name in the local language of this **Place**." example:"Halikarnas Mozolesi"`
	Location      Location          `json:"location" doc:"The geolocation of this **Place** (or its centre if it's a way/relation) based on [WGS 84](https://en.wikipedia.org/wiki/WGS-84)."`
	ReviewSummary ReviewSummary     `json:"review_summary" doc:"Summary of the **Reviews** of this **Place**."`
	OSMType       string            `json:"osm_type" doc:"OpenStreetMap [type](https://wiki.openstreetmap.org/wiki/Elements) of this **Place**." example:"W"`
	OSMID         int64             `json:"osm_id" doc:"OpenStreetMap [ID](https://wiki.openstreetmap.org/wiki/Elements#id) of this **Place**." example:"739181962"`
	OSMTags       map[string]string `json:"osm_tags" doc:"OpenStreetMap [tags](https://wiki.openstreetmap.org/wiki/Tags) of this **Place**." example:"{\"tourism\": \"attraction\", \"historic\": \"yes\", \"name\": \"Halikarnas Mozolesi\", \"name:en\": \"Mausoleum at Halicarnassus\", \"opening_hours\": \"Tu-Su 08:30-17:00\"}"`
}

type PlaceStub struct {
	ID int64 `json:"id" minimum:"1" example:"1"`
}

// Inline ReviewSummary instead of using refs to avoid listing it under Schemas.
func (PlaceStub) Schema(r huma.Registry) *huma.Schema {
	type raw PlaceStub
	return huma.SchemaFromType(r, reflect.TypeOf(raw{}))
}
