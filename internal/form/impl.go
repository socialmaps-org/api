package form

import (
	"net/http"

	"github.com/gorilla/schema"
)

// Set a Decoder instance as a package global, because it caches metadata about
// structs, and an instance can be shared safely.
var decoder = schema.NewDecoder()

func Unmarshal(r *http.Request, v any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	return decoder.Decode(v, r.Form)
}
