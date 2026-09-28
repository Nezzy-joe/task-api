package handlers

import (
	"context"
	"net/http"
)

type authContextKey string

const authKey authContextKey = "authenticated_user"

type AuthContext struct {
	Token string
	User  interface{}
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := extractBearerToken(r.Header.Get("Authorization"))
		if !ok {
			writeAuthJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "Access token required",
			})
			return
		}

		if AuthClient == nil {
			writeAuthJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "Authentication service is not configured",
			})
			return
		}

		userResponse, err := AuthClient.WithToken(token).GetUser()
		if err != nil {
			writeAuthJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "Invalid or expired token",
			})
			return
		}

		auth := AuthContext{
			Token: token,
			User:  userResponse,
		}

		ctx := context.WithValue(r.Context(), authKey, auth)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetAuthContext(r *http.Request) (AuthContext, bool) {
	auth, ok := r.Context().Value(authKey).(AuthContext)
	return auth, ok
}
