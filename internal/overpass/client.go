package overpass

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

const DEFAULT_ENDPOINT = "https://overpass-api.de/api/interpreter"

func Query(q string) (*Response, error) {
	formData := url.Values{
		"data": []string{q},
	}

	httpRes, err := http.DefaultClient.Post(
		DEFAULT_ENDPOINT,
		"application/x-www-form-urlencoded; charset=UTF-8",
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		slog.Info("CANONICAL-OVERPASS-LINE", "status", "error", "error", err.Error(), "query", q)
		return nil, err
	}

	var res Response
	if err = json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
		slog.Info("CANONICAL-OVERPASS-LINE", "status", "error", "error", err.Error(), "query", q)
		return nil, err
	}

	slog.Info("CANONICAL-OVERPASS-LINE", "status", "success", "n_elements", len(res.Elements), "query", q)
	return &res, nil
}

func Retrieve(name string, lat, lon float64) (*Response, error) {
	q := fmt.Sprintf(`[out:json];nwr["name"](around:10,%f,%f);out center tags;`, lat, lon)
	return Query(q)
}
