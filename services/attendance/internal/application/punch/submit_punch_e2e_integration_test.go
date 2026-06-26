//go:build integration

package punch_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	apppunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/application/punch"
	domainpunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/punch"
)

func TestSubmitPunch_E2E_BR010_LowLiveness_REJECTED(t *testing.T) {
	env := newIntegrationEnvWithOpts(t, integrationOpts{
		biometric: configurableBiometricClient{
			IsLive: false, LivenessScore: 0.70,
			IsRecognized: true, RecognitionConfidence: 0.90,
		},
		clock: defaultIntegrationOpts().clock,
	})

	result, err := env.handler.Handle(context.Background(), validPunchCmd(env, nil))
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusRejected, result.Record.Status)
	require.Contains(t, result.Reasons, domainpunch.ReasonLivenessFailed)

	validCount, err := env.punchRepo.CountByStatus(
		context.Background(), env.tenantID, env.employeeID, domainpunch.PunchStatusValid,
	)
	require.NoError(t, err)
	require.Equal(t, 0, validCount)
}

func TestSubmitPunch_E2E_BR010_MockGPS_REJECTED(t *testing.T) {
	env := newIntegrationEnv(t)

	result, err := env.handler.Handle(context.Background(), validPunchCmd(env, func(cmd *apppunch.SubmitPunchCommand) {
		cmd.Location.IsMocked = true
	}))
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusRejected, result.Record.Status)
	require.Contains(t, result.Reasons, domainpunch.ReasonMockGPS)
}

func TestSubmitPunch_E2E_BR010_ClockSkew_REJECTED(t *testing.T) {
	env := newIntegrationEnv(t)

	result, err := env.handler.Handle(context.Background(), validPunchCmd(env, func(cmd *apppunch.SubmitPunchCommand) {
		cmd.DeviceTime = env.serverTime.Add(-10 * time.Minute)
	}))
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusRejected, result.Record.Status)
	require.Contains(t, result.Reasons, domainpunch.ReasonClockManipulation)
}

func TestSubmitPunch_E2E_BR010_Duplicate_REJECTED(t *testing.T) {
	env := newIntegrationEnv(t)

	_, err := env.handler.Handle(context.Background(), validPunchCmd(env, nil))
	require.NoError(t, err)

	result, err := env.handler.Handle(context.Background(), validPunchCmd(env, func(cmd *apppunch.SubmitPunchCommand) {
		cmd.Type = domainpunch.PunchTypeBreakStart
		cmd.DeviceTime = env.serverTime.Add(30 * time.Second)
	}))
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusRejected, result.Record.Status)
	require.Contains(t, result.Reasons, domainpunch.ReasonDuplicatePunch)

	validCount, err := env.punchRepo.CountByStatus(
		context.Background(), env.tenantID, env.employeeID, domainpunch.PunchStatusValid,
	)
	require.NoError(t, err)
	require.Equal(t, 1, validCount)
}

func TestSubmitPunch_E2E_BR011_OfflineExpired_DISCARDED(t *testing.T) {
	env := newIntegrationEnv(t)
	queued := env.serverTime.Add(-9 * time.Hour)

	result, err := env.handler.Handle(context.Background(), validPunchCmd(env, func(cmd *apppunch.SubmitPunchCommand) {
		cmd.IsOfflineSync = true
		cmd.OfflineQueued = &queued
	}))
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusDiscarded, result.Record.Status)
	require.Contains(t, result.Reasons, domainpunch.ReasonOfflineSyncExpired)

	discarded, err := env.punchRepo.CountByStatus(
		context.Background(), env.tenantID, env.employeeID, domainpunch.PunchStatusDiscarded,
	)
	require.NoError(t, err)
	require.Equal(t, 1, discarded)
}

func TestSubmitPunch_E2E_BR011_OfflineWithinTTL_VALID(t *testing.T) {
	env := newIntegrationEnv(t)
	queued := env.serverTime.Add(-7 * time.Hour)

	result, err := env.handler.Handle(context.Background(), validPunchCmd(env, func(cmd *apppunch.SubmitPunchCommand) {
		cmd.IsOfflineSync = true
		cmd.OfflineQueued = &queued
	}))
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusValid, result.Record.Status)
}

func TestSubmitPunch_E2E_BR014_InvalidSequence_REJECTED(t *testing.T) {
	env := newIntegrationEnv(t)

	result, err := env.handler.Handle(context.Background(), validPunchCmd(env, func(cmd *apppunch.SubmitPunchCommand) {
		cmd.Type = domainpunch.PunchTypeClockOut
	}))
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusRejected, result.Record.Status)
	require.Contains(t, result.Reasons, domainpunch.ReasonInvalidSequence)
}

func TestSubmitPunch_E2E_BR015_ServerTimeOfficial(t *testing.T) {
	env := newIntegrationEnv(t)
	deviceTime := env.serverTime.Add(2 * time.Minute)

	result, err := env.handler.Handle(context.Background(), validPunchCmd(env, func(cmd *apppunch.SubmitPunchCommand) {
		cmd.DeviceTime = deviceTime
	}))
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusValid, result.Record.Status)
	require.True(t, result.Record.PunchedAt.Equal(env.serverTime), "BR-015: punchedAt must be server time")
	require.True(t, result.Record.DeviceTime.Equal(deviceTime), "BR-015: deviceTime stored for audit")
}
