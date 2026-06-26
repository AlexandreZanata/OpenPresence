package punch_test

import (
	"time"

	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/geofence"
)

const spLat, spLon = -23.5505, -46.6333

func spCenter() geofence.GpsCoordinate {
	return geofence.GpsCoordinate{Latitude: spLat, Longitude: spLon}
}

func circleZoneAt(id string, center geofence.GpsCoordinate, radiusM float64) geofence.GeofenceZone {
	c := center
	return geofence.GeofenceZone{
		ID: id, Type: geofence.GeofenceTypeCircle,
		Center: &c, RadiusMeters: radiusM, AllowedDeviation: 50,
	}
}

func circleZoneSP(id string, radiusM float64) geofence.GeofenceZone {
	return circleZoneAt(id, spCenter(), radiusM)
}

func polygonZoneSP(id string) geofence.GeofenceZone {
	const half = 0.001
	center := spCenter()
	return geofence.GeofenceZone{
		ID: id, Type: geofence.GeofenceTypePolygon,
		Polygon: []geofence.GpsCoordinate{
			{Latitude: center.Latitude - half, Longitude: center.Longitude - half},
			{Latitude: center.Latitude - half, Longitude: center.Longitude + half},
			{Latitude: center.Latitude + half, Longitude: center.Longitude + half},
			{Latitude: center.Latitude + half, Longitude: center.Longitude - half},
		},
		AllowedDeviation: 50,
	}
}

func expiredCircleZoneSP(id string) geofence.GeofenceZone {
	zone := circleZoneSP(id, 500)
	past := time.Now().Add(-2 * time.Hour)
	zone.ValidUntil = &past
	return zone
}
