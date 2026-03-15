package resource

import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

type Location struct {
	Lat float64 `json:"lat" minimum:"-180.0" maximum:"+180.0" doc:"Latitude in decimal degrees."`
	Lon float64 `json:"lon" minimum:"-90.0" maximum:"+90.0" doc:"Latitude in decimal degrees."`
}

// Inline Location instead of using refs to avoid listing it under Schemas.
func (Location) Schema(r huma.Registry) *huma.Schema {
	type raw Location
	return huma.SchemaFromType(r, reflect.TypeOf(raw{}))
}

type ReviewStats struct {
	Count     int64    `json:"count" minimum:"0" doc:"Number of **Review**s created for this **Place**."`
	LikeRatio *float64 `json:"like_ratio" minimum:"0" maximum:"1" doc:"Ratio of **Review**s that liked this **Place** to all its **Review**s."`
	Score     float64  `json:"score" exclusiveMinimum:"0" exclusiveMaximum:"1" doc:"The likeability score of this **Place** based on the ratio of **Review**s that liked it, the total number of and the recency of its **Review**s."`
}

// Inline ReviewStats instead of using refs to avoid listing it under Schemas.
func (ReviewStats) Schema(r huma.Registry) *huma.Schema {
	type raw ReviewStats
	return huma.SchemaFromType(r, reflect.TypeOf(raw{}))
}

type Place struct {
	Object      string            `json:"object,omitempty" enum:"place"`
	ID          int64             `json:"id" minimum:"1" doc:"Unique identifier for this **Place**."`
	Name        string            `json:"name,omitempty" maxLength:"256" doc:"The common, default name in the local language of this **Place**."`
	Location    Location          `json:"location,omitzero" doc:"The geolocation of this **Place** (or its centre if it's a way/relation) based on [WGS 84](https://en.wikipedia.org/wiki/WGS-84)."`
	ReviewStats ReviewStats       `json:"rating_stats,omitzero" doc:"Statistics about the **Reviews** of this **Place**."`
	OSMType     string            `json:"osm_type" doc:"OpenStreetMap [type](https://wiki.openstreetmap.org/wiki/Elements) of this **Place**."`
	OSMID       int64             `json:"osm_id" doc:"OpenStreetMap [ID](https://wiki.openstreetmap.org/wiki/Elements#id) of this **Place**."`
	OSMTags     map[string]string `json:"osm_tags" doc:"OpenStreetMap [tags](https://wiki.openstreetmap.org/wiki/Tags) of this **Place**."`
}

type PlaceStub struct {
	ID int64 `json:"id" minimum:"1"`
}

// Inline ReviewStats instead of using refs to avoid listing it under Schemas.
func (PlaceStub) Schema(r huma.Registry) *huma.Schema {
	type raw PlaceStub
	return huma.SchemaFromType(r, reflect.TypeOf(raw{}))
}
