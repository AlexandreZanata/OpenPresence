//go:build integration

package punch_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	apppunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/application/punch"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/fraud"
	domainpunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/punch"
)

func TestSubmitPunch_E2E_Fraud_BR012_VPN_SUSPICIOUSInDB(t *testing.T) {
	env := newIntegrationEnv(t)

	result, err := env.handler.Handle(context.Background(), validPunchCmd(env, func(cmd *apppunch.SubmitPunchCommand) {
		cmd.DeviceReport.VPNActive = true
	}))
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusSuspicious, result.Record.Status)
	require.True(t, hasFraudFlag(result.FraudFlags, fraud.FraudTypeVPNDetected))

	suspicious, err := env.punchRepo.CountByStatus(
		context.Background(), env.tenantID, env.employeeID, domainpunch.PunchStatusSuspicious,
	)
	require.NoError(t, err)
	require.Equal(t, 1, suspicious)
}

func TestSubmitPunch_E2E_Fraud_BR012_CriticalClock_REJECTED(t *testing.T) {
	env := newIntegrationEnv(t)

	result, err := env.handler.Handle(context.Background(), validPunchCmd(env, func(cmd *apppunch.SubmitPunchCommand) {
		cmd.DeviceTime = env.serverTime.Add(-31 * time.Minute)
	}))
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusRejected, result.Record.Status)
	require.Contains(t, result.Reasons, domainpunch.ReasonClockManipulation)

	rejected, err := env.punchRepo.CountByStatus(
		context.Background(), env.tenantID, env.employeeID, domainpunch.PunchStatusRejected,
	)
	require.NoError(t, err)
	require.Equal(t, 1, rejected)
}

func TestSubmitPunch_E2E_Fraud_BR013_DeviceLockoutAfterThreeRejects(t *testing.T) {
	current := mustParseTime(integrationBaseTime)
	tracker := fraud.NewDeviceLockoutTracker()
	opts := defaultIntegrationOpts()
	opts.lockout = tracker
	opts.clock = func() time.Time { return current }
	env := newIntegrationEnvWithOpts(t, opts)

	const deviceID = "device-br013"
	reject := func() apppunch.SubmitPunchCommand {
		return validPunchCmd(env, func(cmd *apppunch.SubmitPunchCommand) {
			cmd.DeviceID = deviceID
			cmd.DeviceTime = current
			cmd.Location = domainpunch.GpsCoordinate{Latitude: -22.0, Longitude: -43.0, Accuracy: 10}
		})
	}

	for i := 0; i < 3; i++ {
		result, err := env.handler.Handle(context.Background(), reject())
		require.NoError(t, err, "attempt %d", i+1)
		require.Equal(t, domainpunch.PunchStatusRejected, result.Record.Status)
		current = current.Add(3 * time.Minute)
	}
	require.True(t, tracker.IsLocked(deviceID, current))

	_, err := env.handler.Handle(context.Background(), validPunchCmd(env, func(cmd *apppunch.SubmitPunchCommand) {
		cmd.DeviceID = deviceID
		cmd.DeviceTime = current
	}))
	require.ErrorIs(t, err, apppunch.ErrDeviceLocked)
}
