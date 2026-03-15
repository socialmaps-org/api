package render

import (
	"codeberg.org/socialmaps/api/internal/model"
	"codeberg.org/socialmaps/api/internal/resource"
)

func Place(m model.Place, e model.Element) resource.Place {
	count := m.NLikes + m.NDislikes
	var likeRatio *float64
	if count != 0 {
		t := float64(m.NLikes) / float64(m.NLikes+m.NDislikes)
		likeRatio = &t
	}
	score := (m.DecNLikes + 1.0) / ((m.DecNLikes + 1.0) + (m.DecNDislikes + 1.0))

	return resource.Place{
		Object: "place",
		ID:     m.ID,
		Name:   m.Name,
		Location: resource.Location{
			Lat: m.Lat,
			Lon: m.Lon,
		},
		ReviewStats: resource.ReviewStats{
			Count:     count,
			LikeRatio: likeRatio,
			Score:     score,
		},
		OSMType: e.OsmType,
		OSMID:   e.OsmID,
		OSMTags: e.Tags,
	}
}
