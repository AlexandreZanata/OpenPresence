package httpapi

import (
	"net/http"
)

// NewMux registers attendance HTTP routes for UC-001.
func NewMux(punch *PunchHandler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/v1/attendance/punch", punch)
	mux.HandleFunc("/health/live", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	return mux
}
