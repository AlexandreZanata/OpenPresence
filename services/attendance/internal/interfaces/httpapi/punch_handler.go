package httpapi

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	apppunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/application/punch"
	"github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/fraud"
	domainpunch "github.com/AlexandreZanata/OpenPresence/services/attendance/internal/domain/punch"
)

// PunchHandler exposes POST /v1/attendance/punch.
type PunchHandler struct {
	Submit apppunch.SubmitPunchHandler
}

// ServeHTTP implements UC-001 clock-in over HTTP.
func (h *PunchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}
	actor, ok := AuthFromRequest(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing or invalid bearer token")
		return
	}
	var body PunchRequestDTO
	if err := decodeJSON(r.Body, &body); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid JSON body")
		return
	}
	cmd, err := toCommand(actor, body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}
	result, err := h.Submit.Handle(r.Context(), cmd)
	if err != nil {
		h.writeSubmitError(w, err)
		return
	}
	h.writeResult(w, result)
}

func (h *PunchHandler) writeSubmitError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apppunch.ErrEmployeeNotFound):
		writeError(w, http.StatusForbidden, "FORBIDDEN", "employee not found")
	case errors.Is(err, apppunch.ErrDeviceLocked):
		writeError(w, http.StatusTooManyRequests, "DEVICE_LOCKED", "device locked")
	default:
		writeError(w, http.StatusInternalServerError, "INTERNAL", "punch submission failed")
	}
}

func (h *PunchHandler) writeResult(w http.ResponseWriter, result *apppunch.SubmitPunchResult) {
	resp := PunchResponseDTO{
		ID:        result.Record.ID,
		Status:    string(result.Record.Status),
		PunchedAt: result.Record.PunchedAt.UTC().Format("2006-01-02T15:04:05Z"),
		Type:      string(result.Record.Type),
	}
	switch result.Record.Status {
	case domainpunch.PunchStatusValid:
		writeJSON(w, http.StatusCreated, resp)
	case domainpunch.PunchStatusPending, domainpunch.PunchStatusDiscarded, domainpunch.PunchStatusSuspicious:
		writeJSON(w, http.StatusAccepted, resp)
	default:
		writeError(w, http.StatusUnprocessableEntity, "PUNCH_REJECTED", rejectionMessage(result))
	}
}

func rejectionMessage(result *apppunch.SubmitPunchResult) string {
	if len(result.Reasons) > 0 {
		return string(result.Reasons[0])
	}
	return "punch rejected"
}

func toCommand(actor ActorClaims, body PunchRequestDTO) (apppunch.SubmitPunchCommand, error) {
	deviceTime, err := parseRFC3339(body.DeviceTime)
	if err != nil {
		return apppunch.SubmitPunchCommand{}, err
	}
	frame, err := base64.StdEncoding.DecodeString(body.FrameBase64)
	if err != nil {
		return apppunch.SubmitPunchCommand{}, err
	}
	cmd := apppunch.SubmitPunchCommand{
		TenantID:   actor.TenantID,
		EmployeeID: actor.EmployeeID,
		Type:       domainpunch.PunchType(body.PunchType),
		Location: domainpunch.GpsCoordinate{
			Latitude:  body.Location.Latitude,
			Longitude: body.Location.Longitude,
			Accuracy:  body.Location.Accuracy,
			IsMocked:  body.Location.IsMocked,
		},
		DeviceTime: deviceTime,
		FrameJPEG:  frame,
		DeviceReport: fraud.DeviceIntegrityReport{
			IsRooted:  body.DeviceIntegrityReport.IsRooted,
			VPNActive: body.DeviceIntegrityReport.IsVpnActive,
		},
		IsOfflineSync: body.OfflineSync,
	}
	if body.OfflineSync && body.OfflineQueuedAt != "" {
		queued, err := parseRFC3339(body.OfflineQueuedAt)
		if err != nil {
			return apppunch.SubmitPunchCommand{}, err
		}
		cmd.OfflineQueued = &queued
	}
	return cmd, nil
}

func decodeJSON(r io.Reader, dst any) error {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}
