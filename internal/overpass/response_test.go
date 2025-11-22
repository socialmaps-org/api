package overpass

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWay(t *testing.T) {
	// Arrange
	doc := []byte(`
		{
		"version": 0.6,
		"generator": "Overpass API 0.7.62.8 e802775f",
		"osm3s": {
			"timestamp_osm_base": "2025-09-12T20:55:22Z",
			"copyright": "The data included in this document is from www.openstreetmap.org. The data is made available under ODbL."
		},
		"elements": [
			{
			"type": "way",
			"id": 19830551,
			"center": {
				"lat": 53.3401843,
				"lon": -6.2710249
			},
			"tags": {
				"barrier": "fence",
				"leisure": "park",
				"name": "Saint Patrick's Park",
				"name:es": "Parque de San Patricio",
				"name:it": "Parco di San Patrizio",
				"name:ko": "세인트 패트릭 공원",
				"short_name": "St. Patrick's Park",
				"start_date": "1901",
				"wikidata": "Q113955106",
				"wikimedia_commons": "Category:St. Patrick's Park, Dublin"
			}
			}
		]
		}
	`)

	// Act
	var res Response
	err := json.Unmarshal(doc, &res)

	// Assert
	require.NoError(t, err)
}

func TestNode(t *testing.T) {
	// Arrange
	doc := []byte(`
		{
			"version": 0.6,
			"generator": "Overpass API 0.7.62.8 e802775f",
			"osm3s": {
				"timestamp_osm_base": "2025-10-07T07:33:51Z",
				"copyright": "The data included in this document is from www.openstreetmap.org. The data is made available under ODbL."
			},
			"elements": [
				{
					"type": "node",
					"id": 2084511931,
					"lat": 51.8216582,
					"lon": -8.8556381,
					"tags": {
						"amenity": "pub",
						"name": "Diamond Bar",
						"old_name": "Long's Pub",
						"wikidata": "Q22084291",
						"wikipedia": "en:The Diamond Bar"
					}
				}
			]
		}
	`)

	// Act
	var res Response
	err := json.Unmarshal(doc, &res)

	// Assert
	require.NoError(t, err)
}
