package web

import (
	"encoding/json"
	"fmt"
	"maps"
	"mime"
	"net/http"
)

type Response struct {
	StatusCode int
	Body       any
	Headers    http.Header
}

func NewResponse(statusCode int, body any, headers ...string) *Response {
	var res Response

	res.StatusCode = statusCode
	res.Body = body

	lh := len(headers)
	if lh%2 != 0 {
		panic("headers array must have an even number of elements")
	}

	res.Headers = make(http.Header, lh/2)
	for i := 0; i < lh/2; i++ {
		res.Headers.Set(headers[i], headers[i+1])
	}

	return &res
}

func NewJSONResponse(statusCode int, body any, headers ...string) *Response {
	allHeaders := append(headers, "Content-Type", "application/json")
	return NewResponse(statusCode, body, allHeaders...)
}

func NewEmptyResponse(statusCode int) *Response {
	return NewResponse(statusCode, nil)
}

func (r *Response) send(w http.ResponseWriter) {
	h := w.Header()
	maps.Copy(h, r.Headers)
	w.WriteHeader(r.StatusCode)

	if r.StatusCode == http.StatusNoContent || r.Body == nil {
		return
	}

	mediatype, params, err := mime.ParseMediaType(h.Get("Content-Type"))
	if err != nil {
		panic(err)
	}

	switch mediatype {
	case `application/json`:
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "\t")
		enc.Encode(r.Body)
	case `text/plain`:
		if charset := params["charset"]; charset == "" {
			panic(fmt.Errorf(`charset must be set to "utf-8" for text/plain`))
		} else if charset != "utf-8" {
			panic(fmt.Errorf("unsupported charset for text/plain: %s", params["charset"]))
		}
		s := r.Body.(string)
		w.Write([]byte(s))
	default:
		panic(fmt.Errorf("unprocessable content-type: %s", mediatype))
	}
}
