//go:build integration

package punch_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	apppunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/application/punch"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/geofence"
	domainpunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/punch"
)

var enrollAngles = []string{"FRONTAL", "LEFT_15", "RIGHT_15"}

func TestSubmitPunch_E2E_Stack_BR010_HappyPath_VALIDInDB(t *testing.T) {
	env, bio := newBiometricGrpcEnv(t, nil)
	defer bio.close()

	result, err := env.handler.Handle(context.Background(), stackPunchCmd(env, bio, nil))
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusValid, result.Record.Status)

	count, err := env.punchRepo.CountByStatus(
		context.Background(), env.tenantID, env.employeeID, domainpunch.PunchStatusValid,
	)
	require.NoError(t, err)
	require.Equal(t, 1, count)
}

func TestSubmitPunch_E2E_Stack_BR014_InvalidSequence_REJECTED(t *testing.T) {
	env, bio := newBiometricGrpcEnv(t, nil)
	defer bio.close()
	ctx := context.Background()

	_, err := env.handler.Handle(ctx, stackPunchCmd(env, bio, nil))
	require.NoError(t, err)

	result, err := env.handler.Handle(ctx, stackPunchCmd(env, bio, func(cmd *apppunch.SubmitPunchCommand) {
		cmd.Type = domainpunch.PunchTypeClockIn
		cmd.DeviceTime = env.serverTime.Add(2 * time.Minute)
	}))
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusRejected, result.Record.Status)
	require.Contains(t, result.Reasons, domainpunch.ReasonInvalidSequence)

	validCount, err := env.punchRepo.CountByStatus(
		ctx, env.tenantID, env.employeeID, domainpunch.PunchStatusValid,
	)
	require.NoError(t, err)
	require.Equal(t, 1, validCount)
}

func TestSubmitPunch_E2E_Stack_BR023_AnyAssignedZone(t *testing.T) {
	far := geofence.GpsCoordinate{Latitude: 10, Longitude: 10}
	env, bio := newBiometricGrpcEnv(t, []geofence.GeofenceZone{
		circleZoneAt("zone-far", far, 50),
		circleZoneSP("zone-sp", 500),
	})
	defer bio.close()

	result, err := env.handler.Handle(context.Background(), stackPunchCmd(env, bio, nil))
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusValid, result.Record.Status)
	require.Equal(t, "zone-sp", result.Record.GeofenceID)
}

func TestSubmitPunch_E2E_Stack_OptionalEnrolledProfile(t *testing.T) {
	env, bio := newBiometricGrpcEnv(t, nil)
	defer bio.close()
	ctx := context.Background()

	for _, angle := range enrollAngles {
		enroll, err := bio.raw.EnrollFace(ctx, env.tenantID, env.employeeID, bio.validJPEG, angle)
		require.NoError(t, err)
		require.True(t, enroll.IsLive, "angle %s", angle)
		require.True(t, enroll.HasEmbedding, "angle %s", angle)
	}

	verify, err := bio.raw.VerifyPunch(ctx, env.tenantID, env.employeeID, bio.validJPEG)
	require.NoError(t, err)
	require.True(t, verify.IsLive)

	result, err := env.handler.Handle(ctx, stackPunchCmd(env, bio, nil))
	require.NoError(t, err)
	if verify.IsRecognized {
		require.Equal(t, domainpunch.PunchStatusValid, result.Record.Status)
		return
	}
	require.Equal(t, domainpunch.PunchStatusRejected, result.Record.Status)
	require.Contains(t, result.Reasons, domainpunch.ReasonFaceNotRecognized)
}

func stackPunchCmd(
	env integrationEnv,
	bio *biometricGrpcFixture,
	mutate func(*apppunch.SubmitPunchCommand),
) apppunch.SubmitPunchCommand {
	cmd := validPunchCmd(env, mutate)
	cmd.FrameJPEG = bio.validJPEG
	return cmd
}
