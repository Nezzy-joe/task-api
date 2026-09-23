package handlers

import (
	"net/http"
	"strings"
)

func PublicInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeAuthJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Method not allowed",
		})
		return
	}

	writeAuthJSON(w, http.StatusOK, map[string]string{
		"message": "Welcome stranger! This info is public.",
	})
}

func ProtectedProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeAuthJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Method not allowed",
		})
		return
	}

	token, ok := extractBearerToken(r.Header.Get("Authorization"))
	if !ok {
		writeAuthJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "Access token required",
		})
		return
	}

	// Stage 2 only checks that a bearer token was presented.
	// Token verification happens in Stage 3.
	_ = token

	writeAuthJSON(w, http.StatusOK, map[string]string{
		"message": "Protected profile reached. Token verification comes in Stage 3.",
	})
}

func extractBearerToken(authorization string) (string, bool) {
	parts := strings.Fields(authorization)

	if len(parts) != 2 {
		return "", false
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}

	if parts[1] == "" {
		return "", false
	}

	return parts[1], true
}
