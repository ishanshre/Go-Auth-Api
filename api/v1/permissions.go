package v1

import "net/http"

func permissionDenied(w http.ResponseWriter) {
	writeJSON(w, http.StatusForbidden, ApiError{Error: "permission denied"})
}
