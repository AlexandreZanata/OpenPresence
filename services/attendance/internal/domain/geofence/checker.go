package geofence

import "time"
type GeofenceChecker interface {
	IsInsideZone(coord GpsCoordinate, zone GeofenceZone) bool
	NearestZone(coord GpsCoordinate, zones []GeofenceZone) (*GeofenceZone, float64)
	IsInsideAnyZone(coord GpsCoordinate, zones []GeofenceZone) bool
}

// DefaultChecker implements GeofenceChecker with pure domain logic.
type DefaultChecker struct{}

// NewChecker returns the standard domain geofence checker.
func NewChecker() GeofenceChecker {
	return DefaultChecker{}
}

func (DefaultChecker) IsInsideZone(coord GpsCoordinate, zone GeofenceZone) bool {
	if !zoneIsActive(zone) {
		return false
	}
	switch zone.Type {
	case GeofenceTypeCircle:
		if zone.Center == nil {
			return false
		}
		return IsInsideCircle(coord, *zone.Center, zone.RadiusMeters, zone.AllowedDeviation)
	case GeofenceTypePolygon:
		return IsInsidePolygon(coord, zone.Polygon, zone.AllowedDeviation)
	default:
		return false
	}
}

func (c DefaultChecker) NearestZone(coord GpsCoordinate, zones []GeofenceZone) (*GeofenceZone, float64) {
	var nearest *GeofenceZone
	minDist := -1.0
	for i := range zones {
		zone := &zones[i]
		dist := c.distanceToZone(coord, *zone)
		if nearest == nil || dist < minDist {
			nearest = zone
			minDist = dist
		}
	}
	if nearest == nil {
		return nil, 0
	}
	return nearest, minDist
}

func (c DefaultChecker) IsInsideAnyZone(coord GpsCoordinate, zones []GeofenceZone) bool {
	for _, zone := range zones {
		if c.IsInsideZone(coord, zone) {
			return true
		}
	}
	return false
}

func (DefaultChecker) distanceToZone(coord GpsCoordinate, zone GeofenceZone) float64 {
	switch zone.Type {
	case GeofenceTypeCircle:
		if zone.Center == nil {
			return 0
		}
		dist := HaversineDistance(coord, *zone.Center)
		edge := zone.RadiusMeters + zone.AllowedDeviation
		if dist <= edge {
			return 0
		}
		return dist - edge
	case GeofenceTypePolygon:
		if IsInsidePolygon(coord, zone.Polygon, zone.AllowedDeviation) {
			return 0
		}
		return distanceToPolygonBoundary(coord, zone.Polygon)
	default:
		return 0
	}
}

func zoneIsActive(zone GeofenceZone) bool {
	now := timeNow()
	if zone.ValidFrom != nil && now.Before(*zone.ValidFrom) {
		return false
	}
	if zone.ValidUntil != nil && now.After(*zone.ValidUntil) {
		return false
	}
	return true
}

// timeNow is overridden in tests for BR-024.
var timeNow = func() time.Time {
	return time.Now()
}
