package geofence

// IsInsideCircle reports whether coord is within radiusM plus deviationM of center (BR-020).
func IsInsideCircle(coord, center GpsCoordinate, radiusM, deviationM float64) bool {
	dist := HaversineDistance(coord, center)
	return dist <= radiusM+deviationM
}
