//go:build integration

package punch_test

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	apppunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/application/punch"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/geofence"
	domainpunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/punch"
	infbiometric "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/infrastructure/biometric"
)

const livenessThresholdEnroll = 0.85

func TestSubmitPunch_E2E_BiometricGrpc_BR010_ValidPunch(t *testing.T) {
	env, bio := newBiometricGrpcEnv(t, nil)
	defer bio.close()

	result, err := env.handler.Handle(context.Background(), validPunchCmd(env, func(cmd *apppunch.SubmitPunchCommand) {
		cmd.FrameJPEG = bio.validJPEG
	}))
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusValid, result.Record.Status)

	count, err := env.punchRepo.CountByStatus(
		context.Background(), env.tenantID, env.employeeID, domainpunch.PunchStatusValid,
	)
	require.NoError(t, err)
	require.Equal(t, 1, count)
}

func TestSubmitPunch_E2E_BiometricGrpc_BR010_LowLiveness_REJECTED(t *testing.T) {
	env, bio := newBiometricGrpcEnv(t, nil)
	defer bio.close()

	result, err := env.handler.Handle(context.Background(), validPunchCmd(env, func(cmd *apppunch.SubmitPunchCommand) {
		cmd.FrameJPEG = bio.lowLivenessJPEG
	}))
	require.NoError(t, err)
	require.Equal(t, domainpunch.PunchStatusRejected, result.Record.Status)
	require.Contains(t, result.Reasons, domainpunch.ReasonLivenessFailed)

	count, err := env.punchRepo.CountByStatus(
		context.Background(), env.tenantID, env.employeeID, domainpunch.PunchStatusValid,
	)
	require.NoError(t, err)
	require.Equal(t, 0, count)
}

func TestBiometricGrpc_E2E_BR002_EnrollLivenessRejected(t *testing.T) {
	env, bio := newBiometricGrpcEnv(t, nil)
	defer bio.close()

	result, err := bio.raw.EnrollFace(
		context.Background(), env.tenantID, env.employeeID, bio.lowLivenessJPEG, "FRONTAL",
	)
	require.NoError(t, err)
	require.False(t, result.IsLive)
	require.Less(t, result.LivenessScore, livenessThresholdEnroll)
	require.False(t, result.HasEmbedding)
}

type biometricGrpcFixture struct {
	raw             *infbiometric.Client
	server          *exec.Cmd
	validJPEG       []byte
	lowLivenessJPEG []byte
}

func (f *biometricGrpcFixture) close() {
	if f.raw != nil {
		_ = f.raw.Close()
	}
	if f.server != nil && f.server.Process != nil {
		_ = f.server.Process.Kill()
		_, _ = f.server.Process.Wait()
	}
}

func newBiometricGrpcEnv(t *testing.T, zones []geofence.GeofenceZone) (integrationEnv, *biometricGrpcFixture) {
	return newBiometricGrpcEnvWithClock(t, zones, nil)
}

func newBiometricGrpcEnvWithClock(
	t *testing.T,
	zones []geofence.GeofenceZone,
	clockFn func() time.Time,
) (integrationEnv, *biometricGrpcFixture) {
	t.Helper()
	validJPEG, lowLivenessJPEG := loadBiometricFixtures(t)
	bio := startBiometricServer(t)

	opts := defaultIntegrationOpts()
	opts.biometric = punchBiometricAdapter{client: bio.raw}
	if zones != nil {
		opts.zones = zones
	}
	if clockFn != nil {
		opts.clock = clockFn
	}
	env := newIntegrationEnvWithOpts(t, opts)
	bio.validJPEG = validJPEG
	bio.lowLivenessJPEG = lowLivenessJPEG
	return env, bio
}

type punchBiometricAdapter struct {
	client *infbiometric.Client
}

func (a punchBiometricAdapter) VerifyPunch(
	ctx context.Context, tenantID, employeeID uuid.UUID, frameJPEG []byte,
) (*apppunch.BiometricVerifyResult, error) {
	result, err := a.client.VerifyPunch(ctx, tenantID, employeeID, frameJPEG)
	if err != nil {
		return nil, err
	}
	return &apppunch.BiometricVerifyResult{
		IsLive: result.IsLive, LivenessScore: result.LivenessScore,
		IsRecognized: result.IsRecognized, RecognitionConfidence: result.RecognitionConfidence,
		EmbeddingHash: result.EmbeddingHash,
	}, nil
}

func loadBiometricFixtures(t *testing.T) (valid, lowLiveness []byte) {
	t.Helper()
	dir := biometricFixturesDir(t)
	valid, err := os.ReadFile(filepath.Join(dir, "valid_128.jpg"))
	require.NoError(t, err)
	lowLiveness, err = os.ReadFile(filepath.Join(dir, "low_liveness_128.jpg"))
	require.NoError(t, err)
	return valid, lowLiveness
}

func biometricFixturesDir(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok)
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(filename), "..", "..", "..", "..", ".."))
	return filepath.Join(repoRoot, "services", "biometric", "tests", "fixtures")
}

func startBiometricServer(t *testing.T) *biometricGrpcFixture {
	t.Helper()
	grpcPort := freeTCPPort(t)
	httpPort := freeTCPPort(t)
	svcDir := filepath.Join(repoRoot(t), "services", "biometric")

	cmd := exec.Command("cargo", "run", "--quiet", "--bin", "biometric-server")
	cmd.Dir = svcDir
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

	waitForTCP(t, grpcPort, 30*time.Second)
	client, err := infbiometric.NewClient(fmt.Sprintf("127.0.0.1:%d", grpcPort))
	require.NoError(t, err)

	return &biometricGrpcFixture{raw: client, server: cmd}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok)
	return filepath.Clean(filepath.Join(filepath.Dir(filename), "..", "..", "..", "..", ".."))
}

func freeTCPPort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := ln.Addr().(*net.TCPAddr).Port
	require.NoError(t, ln.Close())
	return port
}

func waitForTCP(t *testing.T, port int, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 200*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("biometric gRPC server not ready on %s", addr)
}
