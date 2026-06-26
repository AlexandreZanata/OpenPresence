package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// PunchRequestDTO is the POST /v1/attendance/punch body (allow-listed fields).
type PunchRequestDTO struct {
	PunchType             string                 `json:"punchType"`
	DeviceTime            string                 `json:"deviceTime"`
	Location              LocationDTO            `json:"location"`
	FrameBase64           string                 `json:"frameBase64"`
	DeviceIntegrityReport DeviceIntegrityDTO     `json:"deviceIntegrityReport"`
	OfflineSync           bool                   `json:"offlineSync"`
	OfflineQueuedAt       string                 `json:"offlineQueuedAt,omitempty"`
}

// LocationDTO is GPS payload from the mobile client.
type LocationDTO struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Accuracy  float64 `json:"accuracy"`
	IsMocked  bool    `json:"isMocked"`
}

// DeviceIntegrityDTO mirrors the device report in API-CONTRACT.md.
type DeviceIntegrityDTO struct {
	IsRooted                  bool `json:"isRooted"`
	IsVpnActive               bool `json:"isVpnActive"`
	IsDeveloperOptionsEnabled bool `json:"isDeveloperOptionsEnabled"`
}

// PunchResponseDTO is returned on successful punch submission.
type PunchResponseDTO struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	PunchedAt string `json:"punchedAt"`
	Type      string `json:"type"`
}

// ErrorResponseDTO follows the API error envelope.
type ErrorResponseDTO struct {
	Error ErrorBody `json:"error"`
}

// ErrorBody holds a safe client-facing error.
type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ActorClaims identifies the authenticated employee for a punch request.
type ActorClaims struct {
	TenantID   uuid.UUID
	EmployeeID uuid.UUID
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, ErrorResponseDTO{Error: ErrorBody{Code: code, Message: message}})
}

func parseRFC3339(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}
