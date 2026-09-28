package handlers

import "net/http"

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
