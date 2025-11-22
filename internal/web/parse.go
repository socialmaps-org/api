package web

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/gorilla/schema"
)

// Set a Decoder instance as a package global, because it caches metadata about
// structs, and an instance can be shared safely.
var queryDecoder = schema.NewDecoder()

func parseRequest(r *http.Request, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.Elem().Kind() != reflect.Struct {
		panic("v must be a pointer to struct")
	}

	if err := parsePath(r, v); err != nil {
		slog.InfoContext(r.Context(), "parse-request-error",
			"component", "path",
			"error", err,
		)
		return err
	}
	if err := parseQuery(r, v); err != nil {
		slog.InfoContext(r.Context(), "parse-request-error",
			"component", "query",
			"error", err,
		)
		return err
	}
	if err := parseJSON(r, v); err != nil {
		slog.InfoContext(r.Context(), "parse-request-error",
			"component", "body",
			"error", err,
		)
		return err
	}
	return nil
}

func parsePath(r *http.Request, v any) error {
	rv := reflect.ValueOf(v)
	elem := rv.Elem()
	typ := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldType := typ.Field(i)

		if pathKey := fieldType.Tag.Get("path"); pathKey != "" {
			if field.Kind() == reflect.String && field.CanSet() {
				pathValue := r.PathValue(pathKey)
				if pathValue == "" {
					return fmt.Errorf("missing path argument: %s", pathKey)
				}
				field.SetString(pathValue)
			}
		}
	}

	return nil
}

func parseQuery(r *http.Request, v any) error {
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("malformed query in the request URL")
	}
	if err := queryDecoder.Decode(v, r.Form); err != nil {
		return err
	}
	return nil
}

func parseJSON(r *http.Request, v any) error {
	rv := reflect.ValueOf(v)
	elem := rv.Elem()
	typ := elem.Type()

	// Check if any field has a json tag
	hasJSONTag := false
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if _, ok := field.Tag.Lookup("json"); ok {
			hasJSONTag = true
			break
		}
	}

	// Return early to avoid throwing an error for the inappropriate
	// Content-Type when there is no expectation of a JSON in the request body.
	if !hasJSONTag {
		return nil
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		return fmt.Errorf("content-type is not application/json")
	}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	err := d.Decode(v)
	if err != nil {
		return fmt.Errorf("malformed JSON in the request body")
	}
	return nil
}
