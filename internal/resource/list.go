package resource

import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

type List[E any] struct {
	Object string `json:"object" enum:"list"`
	Data   []E    `json:"data" nullable:"false"`
}

// Inline List instead of using refs to avoid listing it under Schemas.
func (List[E]) Schema(r huma.Registry) *huma.Schema {
	type raw List[E]
	return huma.SchemaFromType(r, reflect.TypeOf(raw{}))
}
