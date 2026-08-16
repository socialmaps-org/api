package render

import (
	"golang.socialmaps.org/api/internal/model"
	"golang.socialmaps.org/api/internal/resource"
)

func Place(m model.ComputedPlace, e model.Element) resource.Place {
	return resource.Place{
		Object: "place",
		ID:     m.ID,
		Name:   m.Name,
		Location: resource.Location{
			Lat: m.Lat,
			Lon: m.Lon,
		},
		ReviewSummary: resource.ReviewSummary{
			Count: m.NReviews,
			Rating: resource.RatingSummary{
				Average: m.AvgRating,
			},
		},
		OSMType: e.OsmType,
		OSMID:   e.OsmID,
		OSMTags: e.Tags,
	}
}
