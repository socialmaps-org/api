package nominatim

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"codeberg.org/socialmaps/api/internal/geo"
)

// https://nominatim.openstreetmap.org/search?amenity=%22cafe+izz%22&bounded=1&format=jsonv2&viewbox=-8.472,51.895,-8.470,51.896&namedetails=1&extratags=1

type Place struct {
	Type     string
	ID       int64
	Name     string
	Names    map[string]string
	Lat, Lon float64
}

type apiplace struct {
	Type        string            `json:"osm_type"`
	ID          int64             `json:"osm_id"`
	Lat         string            `json:"lat"`
	Lon         string            `json:"lon"`
	NameDetails map[string]string `json:"namedetails"`
}

func Search(ctx context.Context, endpoint string, name string, bbox geo.BBox) ([]Place, error) {
	reqctx, _ := context.WithTimeout(ctx, 1*time.Second)

	req, err := http.NewRequestWithContext(
		reqctx,
		http.MethodGet,
		endpoint+"/search?"+url.Values{
			"amenity": []string{name},
			"format":  []string{"jsonv2"},
			"bounded": []string{"1"},
			"viewbox": []string{
				fmt.Sprintf("%.7f, %.7f, %.7f, %.7f", bbox.West, bbox.South, bbox.East, bbox.North),
			},
			"namedetails": []string{"1"},
			"limit":       []string{"2"},
		}.Encode(),
		nil,
	)
	if err != nil {
		panic(err)
	}

	req.Header.Set("User-Agent", "SocialMaps.org/0.0.0 Contact Bora M. ALPER <bora@boramalper.org>")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		slog.ErrorContext(ctx, "nominatim error", "body", string(b))
		return nil, errors.New("http request failed")
	}

	var apiPlaces []apiplace
	err = json.NewDecoder(res.Body).Decode(&apiPlaces)
	if err != nil {
		return nil, err
	}

	places := make([]Place, len(apiPlaces))
	for i, ap := range apiPlaces {
		lat, err := strconv.ParseFloat(ap.Lat, 64)
		if err != nil {
			return nil, err
		}
		lon, err := strconv.ParseFloat(ap.Lon, 64)
		if err != nil {
			return nil, err
		}

		places[i] = Place{ // Because the top-level `name` object is not the `name` tag in OSM
			// but something else (English name?)
			Name:  ap.NameDetails["name"],
			Names: ap.NameDetails,
			Lat:   lat,
			Lon:   lon,
			Type:  ap.Type,
			ID:    ap.ID,
		}
	}

	return places, nil
}
