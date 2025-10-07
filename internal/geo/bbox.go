package geo

import "math"

type BBox struct {
	South, West, North, East float64
}

const (
	// earthRadiusM is the approximate radius of the Earth in metres.
	earthRadiusM = 6_371_000.0
	// degToRad is the conversion factor from degrees to radians
	degToRad = math.Pi / 180.0
	// radToDeg is the conversion factor from radians to degrees
	radToDeg = 180.0 / math.Pi
)

func NewBBox(lat, lon, radiusM float64) BBox {
	// Convert latitude to radians for the cosine function
	latRad := lat * degToRad

	// Calculate the change in latitude and longitude
	deltaLat := radiusM / earthRadiusM * radToDeg
	deltaLon := radiusM / (earthRadiusM * math.Cos(latRad)) * radToDeg

	return BBox{
		South: lat - deltaLat,
		West:  lon - deltaLon,
		North: lat + deltaLat,
		East:  lon + deltaLon,
	}
}
