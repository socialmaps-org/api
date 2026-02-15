package resource

import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

type User struct {
	ID          int64  `json:"id" minimum:"1" example:"2077851"`
	DisplayName string `json:"display_name" example:"boramalper"`
}

type UserStub struct {
	ID int64 `json:"id" minimum:"1" example:"2077851"`
}

// Inline UserStub instead of using refs to avoid listing it under Schemas.
func (UserStub) Schema(r huma.Registry) *huma.Schema {
	type raw UserStub
	return huma.SchemaFromType(r, reflect.TypeOf(raw{}))
}
