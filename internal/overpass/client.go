package overpass

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	_ "golang.org/x/crypto/x509roots/fallback"
)

func Query(endpoint, q string) (*Response, error) {
	formData := url.Values{
		"data": []string{q},
	}

	httpRes, err := http.DefaultClient.Post(
		endpoint,
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
