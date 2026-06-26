package httpapi

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// AuthFromRequest parses E2E bearer tokens: Bearer e2e.<tenantUUID>.<employeeUUID>.
func AuthFromRequest(r *http.Request) (ActorClaims, bool) {
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return ActorClaims{}, false
	}
	token := strings.TrimPrefix(header, "Bearer ")
	parts := strings.Split(token, ".")
	if len(parts) != 3 || parts[0] != "e2e" {
		return ActorClaims{}, false
	}
	tenantID, err := uuid.Parse(parts[1])
	if err != nil {
		return ActorClaims{}, false
	}
	employeeID, err := uuid.Parse(parts[2])
	if err != nil {
		return ActorClaims{}, false
	}
	return ActorClaims{TenantID: tenantID, EmployeeID: employeeID}, true
}

// BearerToken builds an E2E auth token for tests and curl scripts.
func BearerToken(tenantID, employeeID uuid.UUID) string {
	return "e2e." + tenantID.String() + "." + employeeID.String()
}
