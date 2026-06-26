//go:build integration

package punch_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	apppunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/application/punch"
	domainpunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/punch"
)

func TestSubmitPunch_E2E_RLS_CrossTenantPunchIsolation(t *testing.T) {
	env := newIntegrationEnv(t)
	ctx := context.Background()
	handlerB := env.handlerForTenant(env.otherTenantID, env.otherEmployeeID)

	_, err := handlerB.Handle(ctx, validPunchCmd(env, func(cmd *apppunch.SubmitPunchCommand) {
		cmd.TenantID = env.otherTenantID
		cmd.EmployeeID = env.otherEmployeeID
	}))
	require.NoError(t, err)

	leakCount, err := env.punchRepo.CountByStatus(
		ctx, env.tenantID, env.otherEmployeeID, domainpunch.PunchStatusValid,
	)
	require.NoError(t, err)
	require.Equal(t, 0, leakCount, "tenant A must not count tenant B punches")

	ownCount, err := env.punchRepo.CountByStatus(
		ctx, env.otherTenantID, env.otherEmployeeID, domainpunch.PunchStatusValid,
	)
	require.NoError(t, err)
	require.Equal(t, 1, ownCount, "tenant B must see its own punch")
}
