package v1

import "net/http"

func permissionDenied(w http.ResponseWriter) {
	writeJSON(w, http.StatusForbidden, ApiError{Error: "permission denied"})
}

func invalidParams(w http.ResponseWriter) {
	writeJSON(w, http.StatusForbidden, ApiError{Error: "error parsing parameter"})
}
