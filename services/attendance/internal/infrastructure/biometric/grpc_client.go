package biometric

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	biometricpb "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/infrastructure/biometric/pb"
)

// VerifyResult is the outcome of a biometric VerifyPunch gRPC call.
type VerifyResult struct {
	IsLive                bool
	LivenessScore         float64
	RecognitionConfidence float64
	IsRecognized          bool
	EmbeddingHash         string
}

// EnrollResult is the outcome of a biometric EnrollFace gRPC call.
type EnrollResult struct {
	IsLive        bool
	LivenessScore float64
	QualityScore  float64
	HasEmbedding  bool
}

// Client calls the Rust BiometricService over gRPC.
type Client struct {
	conn   *grpc.ClientConn
	stub   biometricpb.BiometricServiceClient
	target string
}

// NewClient dials target (e.g. "127.0.0.1:19090") without TLS.
func NewClient(target string) (*Client, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("biometric grpc dial: %w", err)
	}
	return &Client{
		conn:   conn,
		stub:   biometricpb.NewBiometricServiceClient(conn),
		target: target,
	}, nil
}

// Close releases the gRPC connection.
func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

// VerifyPunch calls BiometricService.VerifyPunch for a punch frame.
func (c *Client) VerifyPunch(
	ctx context.Context,
	tenantID, employeeID uuid.UUID,
	frameJPEG []byte,
) (*VerifyResult, error) {
	resp, err := c.stub.VerifyPunch(ctx, &biometricpb.VerifyPunchRequest{
		FrameJpeg:  frameJPEG,
		EmployeeId: employeeID.String(),
		TenantId:   tenantID.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("verify punch rpc: %w", err)
	}
	hash := hashEmbedding(resp.GetEmbedding())
	if hash == "" {
		hash = hashFrame(frameJPEG)
	}
	return &VerifyResult{
		IsLive:                resp.GetIsLive(),
		LivenessScore:         float64(resp.GetLivenessScore()),
		RecognitionConfidence: float64(resp.GetRecognitionConfidence()),
		IsRecognized:          resp.GetIsRecognized(),
		EmbeddingHash:         hash,
	}, nil
}

// EnrollFace calls BiometricService.EnrollFace for enrollment flows.
func (c *Client) EnrollFace(
	ctx context.Context,
	tenantID, employeeID uuid.UUID,
	frameJPEG []byte,
	angle string,
) (*EnrollResult, error) {
	resp, err := c.stub.EnrollFace(ctx, &biometricpb.EnrollFaceRequest{
		FrameJpeg:  frameJPEG,
		EmployeeId: employeeID.String(),
		TenantId:   tenantID.String(),
		Angle:      angle,
	})
	if err != nil {
		return nil, fmt.Errorf("enroll face rpc: %w", err)
	}
	return &EnrollResult{
		IsLive:        resp.GetIsLive(),
		LivenessScore: float64(resp.GetLivenessScore()),
		QualityScore:  float64(resp.GetQualityScore()),
		HasEmbedding:  len(resp.GetEmbedding()) > 0,
	}, nil
}

func hashEmbedding(embedding []byte) string {
	if len(embedding) == 0 {
		return ""
	}
	sum := sha256.Sum256(embedding)
	return hex.EncodeToString(sum[:16])
}

func hashFrame(frame []byte) string {
	sum := sha256.Sum256(frame)
	return hex.EncodeToString(sum[:16])
}
