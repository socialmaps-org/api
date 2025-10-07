package web

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, code int, v any) {
	// The implementation below is copied from http.Error(), imagining what
	// http.JSON() would've looked like if Go had it.

	h := w.Header()

	// Delete the Content-Length header, which might be for some other content.
	// Assuming the error string fits in the writer's buffer, we'll figure
	// out the correct Content-Length for it later.
	//
	// We don't delete Content-Encoding, because some middleware sets
	// Content-Encoding: gzip and wraps the ResponseWriter to compress on-the-fly.
	// See https://go.dev/issue/66343.
	h.Del("Content-Length")

	// There might be content type already set, but we reset it to
	// application/json for the value.
	h.Set("Content-Type", "application/json")
	w.WriteHeader(code)

	e := json.NewEncoder(w)
	e.SetEscapeHTML(false)
	e.Encode(v)
}
