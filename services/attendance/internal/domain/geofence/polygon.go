package geofence

import "math"

// IsInsidePolygon reports whether coord is inside polygon or within deviationM of its edges (BR-021).
func IsInsidePolygon(coord GpsCoordinate, polygon []GpsCoordinate, deviationM float64) bool {
	if len(polygon) < 3 {
		return false
	}
	if pointInPolygon(coord, polygon) {
		return true
	}
	if deviationM <= 0 {
		return false
	}
	return distanceToPolygonBoundary(coord, polygon) <= deviationM
}

func pointInPolygon(coord GpsCoordinate, polygon []GpsCoordinate) bool {
	inside := false
	j := len(polygon) - 1
	for i := range polygon {
		xi, yi := polygon[i].Longitude, polygon[i].Latitude
		xj, yj := polygon[j].Longitude, polygon[j].Latitude
		if rayCrossesSegment(coord.Longitude, coord.Latitude, xi, yi, xj, yj) {
			inside = !inside
		}
		j = i
	}
	return inside
}

func rayCrossesSegment(px, py, x1, y1, x2, y2 float64) bool {
	if (y1 > py) == (y2 > py) {
		return false
	}
	xIntersect := (x2-x1)*(py-y1)/(y2-y1) + x1
	return px < xIntersect
}

func distanceToPolygonBoundary(coord GpsCoordinate, polygon []GpsCoordinate) float64 {
	if len(polygon) < 2 {
		return math.MaxFloat64
	}
	minDist := math.MaxFloat64
	for i := range polygon {
		j := (i + 1) % len(polygon)
		d := distanceToSegment(coord, polygon[i], polygon[j])
		if d < minDist {
			minDist = d
		}
	}
	return minDist
}

func distanceToSegment(p, a, b GpsCoordinate) float64 {
	// Project in lat/lng space, then measure with Haversine (adequate for small fences).
	ax, ay := a.Longitude, a.Latitude
	bx, by := b.Longitude, b.Latitude
	px, py := p.Longitude, p.Latitude

	dx, dy := bx-ax, by-ay
	if dx == 0 && dy == 0 {
		return HaversineDistance(p, a)
	}
	t := ((px-ax)*dx + (py-ay)*dy) / (dx*dx + dy*dy)
	t = math.Max(0, math.Min(1, t))
	closest := GpsCoordinate{
		Latitude:  ay + t*dy,
		Longitude: ax + t*dx,
	}
	return HaversineDistance(p, closest)
}
