package handlers

import (
	"net/http"
	"strings"

	"github.com/supabase-community/auth-go/types"
)

// PublicInfoHandler godoc
// @Summary Get public information
// @Description Returns information that does not require authentication.
// @Tags Public
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 405 {object} map[string]string
// @Router /public/info [get]
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

// ProtectedProfileHandler godoc
// @Summary Get authenticated profile
// @Description Returns safe metadata for the authenticated user.
// @Tags Protected
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 405 {object} map[string]string
// @Router /protected/profile [get]
func ProtectedProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeAuthJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Method not allowed",
		})
		return
	}

	auth, ok := GetAuthContext(r)
	if !ok {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Authenticated user context missing",
		})
		return
	}

	userResponse := auth.User

	user := userResponse.(*types.UserResponse)

	writeAuthJSON(w, http.StatusOK, map[string]interface{}{
		"id":                 user.ID,
		"email":              user.Email,
		"account_created_at": user.CreatedAt,
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
