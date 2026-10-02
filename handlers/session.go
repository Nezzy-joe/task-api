package handlers

import "net/http"

// LogoutHandler godoc
// @Summary Log out the authenticated user
// @Description Revokes the authenticated user's refresh tokens through Supabase Auth.
// @Tags Authentication
// @Security BearerAuth
// @Success 204 "No Content"
// @Failure 401 {object} map[string]string
// @Failure 405 {object} map[string]string
// @Router /auth/logout [post]
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	if err := AuthClient.WithToken(auth.Token).Logout(); err != nil {
		writeAuthJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "Logout failed",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ProtectedDashboardHandler godoc
// @Summary Access protected dashboard
// @Description Demonstrates reuse of the authentication middleware on a second protected route.
// @Tags Protected
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 405 {object} map[string]string
// @Router /protected/dashboard [get]
func ProtectedDashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeAuthJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Method not allowed",
		})
		return
	}

	if _, ok := GetAuthContext(r); !ok {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Authenticated user context missing",
		})
		return
	}

	writeAuthJSON(w, http.StatusOK, map[string]string{
		"message": "Welcome to your protected dashboard.",
	})
}
