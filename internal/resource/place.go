package resource

import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

type Location struct {
	Lat float64 `json:"lat" minimum:"-180.0" maximum:"+180.0"`
	Lon float64 `json:"lon" minimum:"-90.0" maximum:"+90.0"`
}

// Inline Location instead of using refs to avoid listing it under Schemas.
func (Location) Schema(r huma.Registry) *huma.Schema {
	type raw Location
	return huma.SchemaFromType(r, reflect.TypeOf(raw{}))
}

type ReviewStats struct {
	Count     int64    `json:"count" minimum:"0"`
	LikeRatio *float64 `json:"like_ratio" minimum:"0" maximum:"1"`
	Score     float64  `json:"score" exclusiveMinimum:"0" exclusiveMaximum:"1"`
}

// Inline ReviewStats instead of using refs to avoid listing it under Schemas.
func (ReviewStats) Schema(r huma.Registry) *huma.Schema {
	type raw ReviewStats
	return huma.SchemaFromType(r, reflect.TypeOf(raw{}))
}

type Place struct {
	Object      string      `json:"object,omitempty" enum:"place"`
	ID          int64       `json:"id" minimum:"1"`
	Name        string      `json:"name,omitempty" maxLength:"256"`
	Location    Location    `json:"location,omitzero"`
	ReviewStats ReviewStats `json:"rating_stats,omitzero"`
}
