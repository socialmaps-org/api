package model

func NewPlaceToComputedPlace(plc Place) ComputedPlace {
	return ComputedPlace{
		ID:       plc.ID,
		Created:  plc.Created,
		Updated:  plc.Updated,
		Name:     plc.Name,
		Location: plc.Location,
		Lat:      plc.Lat,

		Lon:       plc.Lon,
		OsmType:   plc.OsmType,
		OsmID:     plc.OsmID,
		NReviews:  0,
		AvgRating: nil,
	}
}
