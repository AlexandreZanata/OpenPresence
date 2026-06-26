//go:build integration

package httpapi_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
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

	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/app"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/punch"
	infbiometric "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/infrastructure/biometric"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/infrastructure/postgres"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/interfaces/httpapi"
)

const uc001BaseTime = "2026-06-26T09:00:00Z"

func TestPunchAPI_UC001_E2E_FullValidFlow(t *testing.T) {
	env := newAPIEnv(t)
	body := validPunchBody(env.validJPEG, insideLocation(), uc001BaseTime, false, "")

	resp := env.post(t, body, env.tenantID, env.employeeID)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var payload httpapi.PunchResponseDTO
	decodeResp(t, resp, &payload)
	require.Equal(t, "VALID", payload.Status)
	require.Equal(t, "CLOCK_IN", payload.Type)
	require.Equal(t, uc001BaseTime, payload.PunchedAt)

	count, err := env.stack.PunchRepo.CountByStatus(
		context.Background(), env.tenantID, env.employeeID, punch.PunchStatusValid,
	)
	require.NoError(t, err)
	require.Equal(t, 1, count)
}

func TestPunchAPI_UC001_E2E_Unauthorized(t *testing.T) {
	env := newAPIEnv(t)
	req, err := http.NewRequest(http.MethodPost, env.server.URL+"/v1/attendance/punch", bytes.NewReader([]byte("{}")))
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestPunchAPI_UC001_E2E_OutOfGeofence_NoValidPersist(t *testing.T) {
	env := newAPIEnv(t)
	body := validPunchBody(env.validJPEG, outLocation(), uc001BaseTime, false, "")

	resp := env.post(t, body, env.tenantID, env.employeeID)
	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	count, err := env.stack.PunchRepo.CountByStatus(
		context.Background(), env.tenantID, env.employeeID, punch.PunchStatusValid,
	)
	require.NoError(t, err)
	require.Equal(t, 0, count)
}

func TestPunchAPI_UC001_E2E_OfflineSync_VALID(t *testing.T) {
	env := newAPIEnv(t)
	queued := "2026-06-26T08:00:00Z"
	body := validPunchBody(env.validJPEG, insideLocation(), uc001BaseTime, true, queued)

	resp := env.post(t, body, env.tenantID, env.employeeID)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var payload httpapi.PunchResponseDTO
	decodeResp(t, resp, &payload)
	require.Equal(t, "VALID", payload.Status)
}

type apiEnv struct {
	server     *httptest.Server
	stack      app.PunchStack
	tenantID   uuid.UUID
	employeeID uuid.UUID
	validJPEG  []byte
}

func (e apiEnv) post(t *testing.T, body []byte, tenantID, employeeID uuid.UUID) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, e.server.URL+"/v1/attendance/punch", bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+httpapi.BearerToken(tenantID, employeeID))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func newAPIEnv(t *testing.T) apiEnv {
	t.Helper()
	adminDB, appDB := startPostgres(t)
	tenantID, employeeID := seedEmployee(t, adminDB)
	validJPEG := loadFixtureJPEG(t)

	bio := startBiometricServer(t)
	bioClient, err := infbiometric.NewClient(bio.grpcAddr)
	require.NoError(t, err)
	t.Cleanup(func() { _ = bioClient.Close() })

	baseTime := mustParseTime(uc001BaseTime)
	stack := app.NewPunchStack(app.PunchStackConfig{
		DB:        appDB,
		Biometric: app.BiometricGRPCAdapter{Client: bioClient},
		Clock:     func() time.Time { return baseTime },
	})
	handler := &httpapi.PunchHandler{Submit: stack.Handler}
	server := httptest.NewServer(httpapi.NewMux(handler))
	t.Cleanup(server.Close)

	return apiEnv{
		server: server, stack: stack,
		tenantID: tenantID, employeeID: employeeID, validJPEG: validJPEG,
	}
}

func validPunchBody(jpeg []byte, loc map[string]any, deviceTime string, offline bool, queued string) []byte {
	payload := map[string]any{
		"punchType":  "CLOCK_IN",
		"deviceTime": deviceTime,
		"location":   loc,
		"frameBase64": base64.StdEncoding.EncodeToString(jpeg),
		"deviceIntegrityReport": map[string]bool{
			"isRooted": false, "isVpnActive": false, "isDeveloperOptionsEnabled": false,
		},
		"offlineSync": offline,
	}
	if queued != "" {
		payload["offlineQueuedAt"] = queued
	}
	b, _ := json.Marshal(payload)
	return b
}

func insideLocation() map[string]any {
	return map[string]any{"latitude": -23.5505, "longitude": -46.6333, "accuracy": 10.0, "isMocked": false}
}

func outLocation() map[string]any {
	return map[string]any{"latitude": -22.0, "longitude": -43.0, "accuracy": 10.0, "isMocked": false}
}

func decodeResp(t *testing.T, resp *http.Response, dst any) {
	t.Helper()
	defer resp.Body.Close()
	require.NoError(t, json.NewDecoder(resp.Body).Decode(dst))
}

func loadFixtureJPEG(t *testing.T) []byte {
	t.Helper()
	path := filepath.Join(repoRoot(t), "services", "biometric", "tests", "fixtures", "valid_128.jpg")
	b, err := os.ReadFile(path)
	require.NoError(t, err)
	return b
}

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "..", "..", ".."))
}

type bioServer struct {
	grpcAddr string
	cmd      *exec.Cmd
}

func startBiometricServer(t *testing.T) bioServer {
	t.Helper()
	grpcPort := freePort(t)
	httpPort := freePort(t)
	cmd := exec.Command("cargo", "run", "--quiet", "--bin", "biometric-server")
	cmd.Dir = filepath.Join(repoRoot(t), "services", "biometric")
	cmd.Env = append(os.Environ(),
		"BIOMETRIC_USE_STUB=true",
		fmt.Sprintf("BIOMETRIC_GRPC_ADDR=127.0.0.1:%d", grpcPort),
		fmt.Sprintf("BIOMETRIC_HTTP_ADDR=127.0.0.1:%d", httpPort),
		"RUST_LOG=warn",
	)
	require.NoError(t, cmd.Start())
	t.Cleanup(func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
			_, _ = cmd.Process.Wait()
		}
	})
	waitTCP(t, grpcPort)
	return bioServer{grpcAddr: fmt.Sprintf("127.0.0.1:%d", grpcPort), cmd: cmd}
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

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432/tcp")
	appConn := fmt.Sprintf(
		"postgres://attendance_app:attendance_app@%s:%s/openpresence?sslmode=disable",
		host, port.Port(),
	)
	appDB, err := sqlx.Connect("postgres", appConn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = appDB.Close() })
	return adminDB, appDB
}

func seedEmployee(t *testing.T, admin *sqlx.DB) (tenantID, employeeID uuid.UUID) {
	t.Helper()
	ctx := context.Background()
	_, err := admin.ExecContext(ctx, `SET row_security = off`)
	require.NoError(t, err)
	require.NoError(t, admin.QueryRowContext(ctx, `
		INSERT INTO tenants (slug) VALUES ('uc001') RETURNING id`).Scan(&tenantID))
	require.NoError(t, admin.QueryRowContext(ctx, `
		INSERT INTO employees (tenant_id, registration, status)
		VALUES ($1, 'EMP-UC001', 'ACTIVE') RETURNING id`, tenantID).Scan(&employeeID))
	return tenantID, employeeID
}

func migrationsDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)
	abs, err := filepath.Abs(filepath.Join(filepath.Dir(file), "..", "..", "..", "migrations"))
	require.NoError(t, err)
	return abs
}

func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := ln.Addr().(*net.TCPAddr).Port
	require.NoError(t, ln.Close())
	return port
}

func waitTCP(t *testing.T, port int) {
	t.Helper()
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 200*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("port %s not ready", addr)
}

func mustParseTime(value string) time.Time {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		panic(err)
	}
	return t.UTC()
}
