package geofence

import "math"

const earthRadiusMeters = 6_371_000

// HaversineDistance returns the great-circle distance in meters between two coordinates.
func HaversineDistance(a, b GpsCoordinate) float64 {
	lat1 := toRadians(a.Latitude)
	lat2 := toRadians(b.Latitude)
	dLat := toRadians(b.Latitude - a.Latitude)
	dLng := toRadians(b.Longitude - a.Longitude)

	sinDLat := math.Sin(dLat / 2)
	sinDLng := math.Sin(dLng / 2)
	h := sinDLat*sinDLat + math.Cos(lat1)*math.Cos(lat2)*sinDLng*sinDLng
	return 2 * earthRadiusMeters * math.Atan2(math.Sqrt(h), math.Sqrt(1-h))
}

func toRadians(deg float64) float64 {
	return deg * math.Pi / 180
}
