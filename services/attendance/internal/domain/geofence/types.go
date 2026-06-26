package geofence

import "time"

// GeofenceType classifies zone geometry (see docs/GLOSSARY.md).
type GeofenceType string

const (
	GeofenceTypeCircle  GeofenceType = "CIRCLE"
	GeofenceTypePolygon GeofenceType = "POLYGON"
)

// GpsCoordinate is a WGS-84 point (see docs/GLOSSARY.md).
type GpsCoordinate struct {
	Latitude  float64
	Longitude float64
}

// GeofenceZone is a geographic boundary where punch is allowed.
type GeofenceZone struct {
	ID               string
	Type             GeofenceType
	Center           *GpsCoordinate
	RadiusMeters     float64
	Polygon          []GpsCoordinate
	AllowedDeviation float64
	ValidFrom        *time.Time
	ValidUntil       *time.Time
}
