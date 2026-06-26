//go:build integration

package punch_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	apppunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/application/punch"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/fraud"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/geofence"
	domainpunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/punch"
)

func TestSubmitPunch_E2E_Geofence_BR020_CircleInside_VALID(t *testing.T) {
	env := newIntegrationEnvWithZones(t, []geofence.GeofenceZone{circleZoneSP("circle-sp", 500)})

	result, err := env.handler.Handle(context.Background(), validPunchCmd(env, nil))
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusValid, result.Record.Status)
	require.Equal(t, "circle-sp", result.Record.GeofenceID)
}

func TestSubmitPunch_E2E_Geofence_BR020_CircleOutside_REJECTED(t *testing.T) {
	env := newIntegrationEnvWithZones(t, []geofence.GeofenceZone{circleZoneSP("circle-sp", 500)})

	result, err := env.handler.Handle(context.Background(), validPunchCmd(env, func(cmd *apppunch.SubmitPunchCommand) {
		cmd.Location = domainpunch.GpsCoordinate{Latitude: -22.0, Longitude: -43.0, Accuracy: 10}
	}))
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusRejected, result.Record.Status)
	require.Contains(t, result.Reasons, domainpunch.ReasonOutOfGeofence)
}

func TestSubmitPunch_E2E_Geofence_BR021_PolygonInside_VALID(t *testing.T) {
	env := newIntegrationEnvWithZones(t, []geofence.GeofenceZone{polygonZoneSP("poly-sp")})

	result, err := env.handler.Handle(context.Background(), validPunchCmd(env, nil))
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusValid, result.Record.Status)
	require.Equal(t, "poly-sp", result.Record.GeofenceID)
}

func TestSubmitPunch_E2E_Geofence_BR022_LowAccuracyFlag_VALID(t *testing.T) {
	env := newIntegrationEnvWithZones(t, []geofence.GeofenceZone{circleZoneSP("circle-sp", 500)})

	result, err := env.handler.Handle(context.Background(), validPunchCmd(env, func(cmd *apppunch.SubmitPunchCommand) {
		cmd.Location.Accuracy = 150 // > allowedDeviation*2 (50*2)
	}))
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusValid, result.Record.Status)
	require.True(t, hasFraudFlag(result.FraudFlags, fraud.FraudTypeGPSLowAccuracy))
}

func TestSubmitPunch_E2E_Geofence_BR023_AnyAssignedZone_VALID(t *testing.T) {
	far := geofence.GpsCoordinate{Latitude: 10, Longitude: 10}
	env := newIntegrationEnvWithZones(t, []geofence.GeofenceZone{
		circleZoneAt("zone-far", far, 50),
		circleZoneSP("zone-sp", 500),
	})

	result, err := env.handler.Handle(context.Background(), validPunchCmd(env, nil))
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusValid, result.Record.Status)
	require.Equal(t, "zone-sp", result.Record.GeofenceID)
}

func TestSubmitPunch_E2E_Geofence_BR024_ExpiredZoneIgnored_REJECTED(t *testing.T) {
	env := newIntegrationEnvWithZones(t, []geofence.GeofenceZone{expiredCircleZoneSP("expired-sp")})

	result, err := env.handler.Handle(context.Background(), validPunchCmd(env, nil))
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusRejected, result.Record.Status)
	require.Contains(t, result.Reasons, domainpunch.ReasonOutOfGeofence)
}

func hasFraudFlag(flags []fraud.FraudFlag, target fraud.FraudType) bool {
	for _, f := range flags {
		if f.Type == target {
			return true
		}
	}
	return false
}
