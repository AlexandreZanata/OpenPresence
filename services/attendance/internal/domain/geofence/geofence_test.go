package geofence

import (
	"math"
	"testing"
	"time"
)

func TestHaversineDistance_KnownPair(t *testing.T) {
	a := GpsCoordinate{Latitude: -12.5458, Longitude: -55.7061}
	b := GpsCoordinate{Latitude: -12.5448, Longitude: -55.7061}
	got := HaversineDistance(a, b)
	want := 111.19
	if math.Abs(got-want) > 1.0 {
		t.Fatalf("distance %f m, want within 1m of %f", got, want)
	}
}

func TestIsInsideCircle_CenterPoint(t *testing.T) {
	center := GpsCoordinate{Latitude: -12.5458, Longitude: -55.7061}
	if !IsInsideCircle(center, center, 100, 0) {
		t.Fatal("center must be inside circle")
	}
}

func TestIsInsideCircle_OnBoundary(t *testing.T) {
	center := GpsCoordinate{Latitude: 0, Longitude: 0}
	edge := offsetNorth(center, 100)
	if !IsInsideCircle(edge, center, 100, 0) {
		t.Fatal("point on boundary must be inside")
	}
}

func TestIsInsideCircle_JustOutside(t *testing.T) {
	center := GpsCoordinate{Latitude: 0, Longitude: 0}
	outside := offsetNorth(center, 101)
	if IsInsideCircle(outside, center, 100, 0) {
		t.Fatal("point outside radius must not be inside")
	}
}

func TestIsInsideCircle_WithDeviation(t *testing.T) {
	center := GpsCoordinate{Latitude: 0, Longitude: 0}
	near := offsetNorth(center, 110)
	if !IsInsideCircle(near, center, 100, 15) {
		t.Fatal("BR-020: within radius+deviation must pass")
	}
}

func TestIsInsidePolygon_InsideConvex(t *testing.T) {
	square := convexSquare(GpsCoordinate{Latitude: 0, Longitude: 0}, 0.001)
	inside := GpsCoordinate{Latitude: 0, Longitude: 0}
	if !IsInsidePolygon(inside, square, 0) {
		t.Fatal("center of convex square must be inside")
	}
}

func TestIsInsidePolygon_Outside(t *testing.T) {
	square := convexSquare(GpsCoordinate{Latitude: 0, Longitude: 0}, 0.001)
	outside := GpsCoordinate{Latitude: 1, Longitude: 1}
	if IsInsidePolygon(outside, square, 0) {
		t.Fatal("remote point must be outside")
	}
}

func TestIsInsidePolygon_ConcaveShape(t *testing.T) {
	concave := concaveLShape()
	inside := GpsCoordinate{Latitude: 0.0002, Longitude: 0.0002}
	if !IsInsidePolygon(inside, concave, 0) {
		t.Fatal("point in concave arm must be inside")
	}
}

func TestIsInsideAnyZone_MultipleZones(t *testing.T) {
	checker := NewChecker()
	zones := []GeofenceZone{
		{
			ID:           "z1",
			Type:         GeofenceTypeCircle,
			Center:       &GpsCoordinate{Latitude: 10, Longitude: 10},
			RadiusMeters: 50,
		},
		{
			ID:           "z2",
			Type:         GeofenceTypeCircle,
			Center:       &GpsCoordinate{Latitude: 0, Longitude: 0},
			RadiusMeters: 200,
		},
	}
	coord := GpsCoordinate{Latitude: 0, Longitude: 0}
	if !checker.IsInsideAnyZone(coord, zones) {
		t.Fatal("BR-023: inside any assigned zone must pass")
	}
}

func TestNearestZone_ReturnsClosest(t *testing.T) {
	checker := NewChecker()
	centerNear := GpsCoordinate{Latitude: 0, Longitude: 0}
	centerFar := GpsCoordinate{Latitude: 1, Longitude: 1}
	zones := []GeofenceZone{
		{ID: "far", Type: GeofenceTypeCircle, Center: &centerFar, RadiusMeters: 10},
		{ID: "near", Type: GeofenceTypeCircle, Center: &centerNear, RadiusMeters: 500},
	}
	coord := GpsCoordinate{Latitude: 0.0001, Longitude: 0.0001}
	zone, dist := checker.NearestZone(coord, zones)
	if zone == nil || zone.ID != "near" {
		t.Fatalf("expected nearest zone 'near', got %v", zone)
	}
	if dist < 0 {
		t.Fatalf("distance must be non-negative, got %f", dist)
	}
}

func TestIsInsideZone_PolygonChecker(t *testing.T) {
	checker := NewChecker()
	square := convexSquare(GpsCoordinate{Latitude: 0, Longitude: 0}, 0.001)
	zone := GeofenceZone{
		ID:      "poly",
		Type:    GeofenceTypePolygon,
		Polygon: square,
	}
	coord := GpsCoordinate{Latitude: 0, Longitude: 0}
	if !checker.IsInsideZone(coord, zone) {
		t.Fatal("checker must delegate polygon zones")
	}
}

func TestZoneValidity_BR024(t *testing.T) {
	orig := timeNow
	defer func() { timeNow = orig }()
	fixed := time.Date(2026, 6, 26, 12, 0, 0, 0, time.UTC)
	timeNow = func() time.Time { return fixed }

	past := fixed.Add(-1 * time.Hour)
	center := GpsCoordinate{Latitude: 0, Longitude: 0}
	zone := GeofenceZone{
		ID:           "temp",
		Type:         GeofenceTypeCircle,
		Center:       &center,
		RadiusMeters: 500,
		ValidUntil:   &past,
	}
	checker := NewChecker()
	if checker.IsInsideZone(center, zone) {
		t.Fatal("BR-024: inactive zone must not match")
	}
}

func offsetNorth(origin GpsCoordinate, meters float64) GpsCoordinate {
	deltaLat := meters / 111_320.0
	return GpsCoordinate{
		Latitude:  origin.Latitude + deltaLat,
		Longitude: origin.Longitude,
	}
}

func convexSquare(center GpsCoordinate, halfDeg float64) []GpsCoordinate {
	return []GpsCoordinate{
		{Latitude: center.Latitude - halfDeg, Longitude: center.Longitude - halfDeg},
		{Latitude: center.Latitude - halfDeg, Longitude: center.Longitude + halfDeg},
		{Latitude: center.Latitude + halfDeg, Longitude: center.Longitude + halfDeg},
		{Latitude: center.Latitude + halfDeg, Longitude: center.Longitude - halfDeg},
	}
}

func concaveLShape() []GpsCoordinate {
	return []GpsCoordinate{
		{Latitude: 0, Longitude: 0},
		{Latitude: 0, Longitude: 0.001},
		{Latitude: 0.0005, Longitude: 0.001},
		{Latitude: 0.0005, Longitude: 0.0005},
		{Latitude: 0.001, Longitude: 0.0005},
		{Latitude: 0.001, Longitude: 0},
	}
}
