package app

import (
	"context"

	"github.com/google/uuid"

	apppunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/application/punch"
	infbiometric "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/infrastructure/biometric"
)

// BiometricGRPCAdapter maps infrastructure gRPC client to the punch port.
type BiometricGRPCAdapter struct {
	Client *infbiometric.Client
}

func (a BiometricGRPCAdapter) VerifyPunch(
	ctx context.Context, tenantID, employeeID uuid.UUID, frameJPEG []byte,
) (*apppunch.BiometricVerifyResult, error) {
	result, err := a.Client.VerifyPunch(ctx, tenantID, employeeID, frameJPEG)
	if err != nil {
		return nil, err
	}
	return &apppunch.BiometricVerifyResult{
		IsLive: result.IsLive, LivenessScore: result.LivenessScore,
		IsRecognized: result.IsRecognized, RecognitionConfidence: result.RecognitionConfidence,
		EmbeddingHash: result.EmbeddingHash,
	}, nil
}
