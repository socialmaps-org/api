package web

import (
	"encoding/json"
	"net/http"
	"reflect"

	"codeberg.org/socialmaps/auth/internal/resource"
	"github.com/gorilla/schema"
)

// Set a Decoder instance as a package global, because it caches metadata about
// structs, and an instance can be shared safely.
var queryDecoder = schema.NewDecoder()

func Parse(r *http.Request, v any) *Error {
	if err := parsePath(r, v); err != nil {
		return err
	}
	if err := parseQuery(r, v); err != nil {
		return err
	}
	if err := parseJSON(r, v); err != nil {
		return err
	}
	return nil
}

func parsePath(r *http.Request, v any) *Error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.Elem().Kind() != reflect.Struct {
		panic("v must be a pointer to struct")
	}

	elem := rv.Elem()
	typ := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldType := typ.Field(i)

		if pathKey := fieldType.Tag.Get("path"); pathKey != "" {
			if field.Kind() == reflect.String && field.CanSet() {
				field.SetString(r.PathValue(pathKey))
			}
		}
	}

	return nil
}

func parseQuery(r *http.Request, v any) *Error {
	if err := r.ParseForm(); err != nil {
		return &Error{
			StatusCode: http.StatusUnprocessableEntity,
			Resource:   resource.Error{Message: "cannot parse query params"},
		}
	}
	if err := queryDecoder.Decode(v, r.Form); err != nil {
		return &Error{
			StatusCode: http.StatusBadRequest,
			Resource:   resource.Error{Message: "invalid query params"},
		}
	}
	return nil
}

func parseJSON(r *http.Request, v any) *Error {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		return &Error{
			StatusCode: http.StatusUnsupportedMediaType,
			Resource:   resource.Error{Message: "content-type is not application/json"},
		}
	}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	err := d.Decode(v)
	if err != nil {
		return &Error{
			StatusCode: http.StatusUnprocessableEntity,
			Resource:   resource.Error{Message: "cannot parse JSON body"},
		}
	}
	return nil
}
