//go:build integration

package authorization_test

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	appauth "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/application/authorization"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/organization"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/infrastructure/postgres"
)

func TestPunchAuthorization_E2E_ManagerApprovesNurseInSubtree(t *testing.T) {
	env := newAuthE2EEnv(t)
	nurseID := env.employeeIDs["nurse-1"]

	ok, err := env.handler.Approve(context.Background(), appauth.AuthorizePunchApprovalCommand{
		Actor: organization.ActorScope{
			Role: organization.RoleManager, TenantID: env.tenantA.String(),
			AssignedOrgNodeID: "health",
		},
		TenantID: env.tenantA, EmployeeID: nurseID,
	})
	require.NoError(t, err)
	require.True(t, ok)
}

func TestPunchAuthorization_E2E_ManagerRejectsEducationEmployee(t *testing.T) {
	env := newAuthE2EEnv(t)
	teacherID := env.employeeIDs["teacher-1"]

	ok, err := env.handler.Approve(context.Background(), appauth.AuthorizePunchApprovalCommand{
		Actor: organization.ActorScope{
			Role: organization.RoleManager, TenantID: env.tenantA.String(),
			AssignedOrgNodeID: "health",
		},
		TenantID: env.tenantA, EmployeeID: teacherID,
	})
	require.NoError(t, err)
	require.False(t, ok)
}

func TestPunchAuthorization_E2E_AuditorWriteDenied(t *testing.T) {
	env := newAuthE2EEnv(t)

	err := env.handler.AuthorizeWrite(organization.ActorScope{
		Role: organization.RoleAuditor, TenantID: env.tenantA.String(),
	})
	require.ErrorIs(t, err, appauth.ErrWriteDenied)

	canRead, err := env.auth.ReadPunch(
		organization.ActorScope{Role: organization.RoleAuditor, TenantID: env.tenantA.String()},
		organization.EmployeePlacementRef{
			EmployeeID: env.employeeIDs["nurse-1"].String(),
			TenantID:   env.tenantA.String(),
			OrgNodeID:  "nursing",
		},
	)
	require.NoError(t, err)
	require.True(t, canRead)
}

func TestPunchAuthorization_E2E_CrossTenantActorDenied(t *testing.T) {
	env := newAuthE2EEnv(t)
	nurseID := env.employeeIDs["nurse-1"]

	_, err := env.handler.Approve(context.Background(), appauth.AuthorizePunchApprovalCommand{
		Actor: organization.ActorScope{
			Role: organization.RoleSuperAdmin, TenantID: env.tenantB.String(),
		},
		TenantID: env.tenantA, EmployeeID: nurseID,
	})
	require.Error(t, err)
	require.True(t, errors.Is(err, organization.ErrCrossTenantAccess))
}

type authE2EEnv struct {
	handler     appauth.AuthorizePunchApprovalHandler
	auth        *appauth.PunchAuthorizationService
	tenantA     uuid.UUID
	tenantB     uuid.UUID
	employeeIDs map[string]uuid.UUID
}

func newAuthE2EEnv(t *testing.T) authE2EEnv {
	t.Helper()
	adminDB, appDB := startPostgres(t)
	tenantA, tenantB, employees := seedAuthEmployees(t, adminDB)
	empRepo := postgres.NewEmployeeRepository(appDB)

	treeReader := mapTreeReader{trees: map[string]*organization.OrgTree{
		tenantA.String(): municipalTree(t, tenantA.String()),
	}}
	authSvc := appauth.NewPunchAuthorizationService(treeReader)
	handler := appauth.AuthorizePunchApprovalHandler{
		Employees:  postgresEmployeeAdapter{repo: empRepo},
		Placements: fixturePlacementReader{placements: employees},
		Auth:       authSvc,
	}
	return authE2EEnv{
		handler: handler, auth: authSvc,
		tenantA: tenantA, tenantB: tenantB, employeeIDs: employees,
	}
}

type mapTreeReader struct {
	trees map[string]*organization.OrgTree
}

func (m mapTreeReader) Tree(tenantID string) (*organization.OrgTree, error) {
	tree, ok := m.trees[tenantID]
	if !ok {
		return nil, organization.ErrEmptyTree
	}
	return tree, nil
}

type postgresEmployeeAdapter struct {
	repo *postgres.EmployeeRepository
}

func (a postgresEmployeeAdapter) GetEmployee(
	ctx context.Context, tenantID, employeeID uuid.UUID,
) (*appauth.TenantEmployee, error) {
	emp, err := a.repo.GetEmployee(ctx, tenantID, employeeID)
	if err != nil || emp == nil {
		return nil, err
	}
	return &appauth.TenantEmployee{
		ID: emp.ID, TenantID: emp.TenantID, Registration: emp.Registration,
	}, nil
}

type fixturePlacementReader struct {
	placements map[string]uuid.UUID
}

func (f fixturePlacementReader) EmployeePlacement(
	_ context.Context, tenantID, employeeID uuid.UUID,
) (organization.EmployeePlacementRef, error) {
	for reg, id := range f.placements {
		if id == employeeID {
			node := "nursing"
			if reg == "teacher-1" {
				node = "school-a"
			}
			return organization.EmployeePlacementRef{
				EmployeeID: employeeID.String(),
				TenantID:   tenantID.String(),
				OrgNodeID:  node,
			}, nil
		}
	}
	return organization.EmployeePlacementRef{}, fmt.Errorf("placement not found")
}

func seedAuthEmployees(t *testing.T, admin *sqlx.DB) (tenantA, tenantB uuid.UUID, ids map[string]uuid.UUID) {
	t.Helper()
	ctx := context.Background()
	_, err := admin.ExecContext(ctx, `SET row_security = off`)
	require.NoError(t, err)

	require.NoError(t, admin.QueryRowContext(ctx, `
		INSERT INTO tenants (slug) VALUES ('tenant-a') RETURNING id`).Scan(&tenantA))
	require.NoError(t, admin.QueryRowContext(ctx, `
		INSERT INTO tenants (slug) VALUES ('tenant-b') RETURNING id`).Scan(&tenantB))

	ids = make(map[string]uuid.UUID)
	for reg := range map[string]string{"nurse-1": "nursing", "teacher-1": "school-a"} {
		var id uuid.UUID
		require.NoError(t, admin.QueryRowContext(ctx, `
			INSERT INTO employees (tenant_id, registration, status)
			VALUES ($1, $2, 'ACTIVE') RETURNING id`, tenantA, reg).Scan(&id))
		ids[reg] = id
	}
	return tenantA, tenantB, ids
}

func startPostgres(t *testing.T) (admin, app *sqlx.DB) {
	t.Helper()
	ctx := context.Background()
	container, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("openpresence"),
		tcpostgres.WithUsername("openpresence"),
		tcpostgres.WithPassword("openpresence"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(60*time.Second),
		),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = container.Terminate(ctx) })

	adminConn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	adminDB, err := sqlx.Connect("postgres", adminConn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = adminDB.Close() })

	require.NoError(t, postgres.ApplyMigrations(adminDB.DB, migrationsDir(t)))

	host, err := container.Host(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err)
	appConn := fmt.Sprintf(
		"postgres://attendance_app:attendance_app@%s:%s/openpresence?sslmode=disable",
		host, port.Port(),
	)
	appDB, err := sqlx.Connect("postgres", appConn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = appDB.Close() })
	return adminDB, appDB
}

func migrationsDir(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok)
	dir := filepath.Join(filepath.Dir(filename), "..", "..", "..", "migrations")
	abs, err := filepath.Abs(dir)
	require.NoError(t, err)
	return abs
}
