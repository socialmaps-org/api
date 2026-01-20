package resource

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type ReviewStats struct {
	Count     int64    `json:"count"`
	LikeRatio *float64 `json:"like_ratio"`
	Score     float64  `json:"score"`
}

type Place struct {
	Object      string      `json:"object,omitempty"`
	ID          int64       `json:"id"`
	Name        string      `json:"name,omitempty"`
	Location    Location    `json:"location,omitzero"`
	ReviewStats ReviewStats `json:"rating_stats,omitzero"`
}
