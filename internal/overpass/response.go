package overpass

import "time"

type Element struct {
	Type   string
	ID     uint64
	Center *struct {
		Lat float64
		Lon float64
	}
	Lat_ float64 `json:"lat"`
	Lon_ float64 `json:"lon"`
	Tags map[string]string
}

type Response struct {
	Version   float64
	Generator string
	OSM3S     struct {
		TimestampOSMBase time.Time
	}
	Elements []Element
}

func (el *Element) Name() string {
	return el.Tags["name"]
}

func (el *Element) Lat() float64 {
	if el.Center != nil {
		return el.Center.Lat
	}
	return el.Lat_
}

func (el *Element) Lon() float64 {
	if el.Center != nil {
		return el.Center.Lon
	}
	return el.Lon_
}
