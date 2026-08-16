package model

func (plc OptionalComputedPlace) IsNil() bool {
	return plc.ID == nil
}

func (plc OptionalComputedPlace) Unwrap() ComputedPlace {
	return ComputedPlace{
		ID:        *plc.ID,
		Created:   *plc.Created,
		Updated:   *plc.Updated,
		Name:      *plc.Name,
		Location:  plc.Location,
		Lat:       *plc.Lat,
		Lon:       *plc.Lon,
		OsmType:   *plc.OsmType,
		OsmID:     *plc.OsmID,
		NReviews:  *plc.NReviews,
		AvgRating: plc.AvgRating,
	}
}
