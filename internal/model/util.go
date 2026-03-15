package model

func (plc OptionalPlace) IsNil() bool {
	return plc.ID == nil
}

func (plc OptionalPlace) Unwrap() Place {
	return Place{
		ID:           *plc.ID,
		Created:      *plc.Created,
		Updated:      *plc.Updated,
		Name:         *plc.Name,
		Location:     plc.Location,
		Lat:          *plc.Lat,
		Lon:          *plc.Lon,
		OsmType:      *plc.OsmType,
		OsmID:        *plc.OsmID,
		NLikes:       *plc.NLikes,
		NDislikes:    *plc.NDislikes,
		DecNLikes:    *plc.DecNLikes,
		DecNDislikes: *plc.DecNDislikes,
		DecUpdatedAt: *plc.DecUpdatedAt,
		Score:        *plc.Score,
	}

}
