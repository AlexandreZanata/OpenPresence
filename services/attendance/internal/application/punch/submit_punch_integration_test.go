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

func TestSubmitPunch_Integration_HappyPath_VALIDInDB(t *testing.T) {
	env := newIntegrationEnv(t)

	result, err := env.handler.Handle(context.Background(), validPunchCmd(env, nil))
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusValid, result.Record.Status)

	count, err := env.punchRepo.CountByStatus(
		context.Background(), env.tenantID, env.employeeID, domainpunch.PunchStatusValid,
	)
	require.NoError(t, err)
	require.Equal(t, 1, count)
}

func TestSubmitPunch_Integration_CrossTenant_Rejected(t *testing.T) {
	env := newIntegrationEnv(t)

	_, err := env.handler.Handle(context.Background(), apppunch.SubmitPunchCommand{
		TenantID:   env.otherTenantID,
		EmployeeID: env.employeeID,
		Type:       domainpunch.PunchTypeClockIn,
		Location:   insideLocation(),
		DeviceTime: env.serverTime,
		FrameJPEG:  []byte("frame"),
	})
	require.ErrorIs(t, err, apppunch.ErrEmployeeNotFound)

	count, err := env.punchRepo.CountByStatus(
		context.Background(), env.tenantID, env.employeeID, domainpunch.PunchStatusValid,
	)
	require.NoError(t, err)
	require.Equal(t, 0, count)
}

func TestSubmitPunch_Integration_InvalidSequence_REJECTED(t *testing.T) {
	env := newIntegrationEnv(t)

	_, err := env.handler.Handle(context.Background(), validPunchCmd(env, nil))
	require.NoError(t, err)

	result, err := env.handler.Handle(context.Background(), validPunchCmd(env, func(cmd *apppunch.SubmitPunchCommand) {
		cmd.Type = domainpunch.PunchTypeClockIn
		cmd.DeviceTime = env.serverTime.Add(2 * time.Minute)
	}))
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusRejected, result.Record.Status)

	validCount, err := env.punchRepo.CountByStatus(
		context.Background(), env.tenantID, env.employeeID, domainpunch.PunchStatusValid,
	)
	require.NoError(t, err)
	require.Equal(t, 1, validCount)
}

func TestSubmitPunch_Integration_OutOfGeofence_REJECTED(t *testing.T) {
	env := newIntegrationEnv(t)

	result, err := env.handler.Handle(context.Background(), validPunchCmd(env, func(cmd *apppunch.SubmitPunchCommand) {
		cmd.Location = domainpunch.GpsCoordinate{Latitude: -22.0, Longitude: -43.0, Accuracy: 10}
	}))
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusRejected, result.Record.Status)

	validCount, err := env.punchRepo.CountByStatus(
		context.Background(), env.tenantID, env.employeeID, domainpunch.PunchStatusValid,
	)
	require.NoError(t, err)
	require.Equal(t, 0, validCount)
}
