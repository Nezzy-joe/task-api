package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	supabaseauth "github.com/supabase-community/auth-go"
	"github.com/supabase-community/auth-go/types"
)

// AuthClient is initialized in main.go with the Supabase project reference
// and anon key.
var AuthClient supabaseauth.Client

type authCredentialsRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	User         types.User `json:"user"`
	AccessToken  string     `json:"access_token,omitempty"`
	RefreshToken string     `json:"refresh_token,omitempty"`
	TokenType    string     `json:"token_type,omitempty"`
	ExpiresIn    int        `json:"expires_in,omitempty"`
	ExpiresAt    int64      `json:"expires_at,omitempty"`
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeAuthJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Method not allowed",
		})
		return
	}

	var req authCredentialsRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid JSON body",
		})
		return
	}

	email := strings.TrimSpace(req.Email)

	if email == "" || req.Password == "" {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Email and password are required",
		})
		return
	}

	if AuthClient == nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Authentication service is not configured",
		})
		return
	}

	resp, err := AuthClient.Signup(types.SignupRequest{
		Email:    email,
		Password: req.Password,
	})
	if err != nil {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Signup failed",
		})
		return
	}

	response := authResponse{
		User:         resp.User,
		AccessToken:  resp.Session.AccessToken,
		RefreshToken: resp.Session.RefreshToken,
		TokenType:    resp.Session.TokenType,
		ExpiresIn:    resp.Session.ExpiresIn,
		ExpiresAt:    resp.Session.ExpiresAt,
	}

	writeAuthJSON(w, http.StatusCreated, response)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeAuthJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Method not allowed",
		})
		return
	}

	var req authCredentialsRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid JSON body",
		})
		return
	}

	email := strings.TrimSpace(req.Email)

	if email == "" || req.Password == "" {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Email and password are required",
		})
		return
	}

	if AuthClient == nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Authentication service is not configured",
		})
		return
	}

	resp, err := AuthClient.SignInWithEmailPassword(email, req.Password)
	if err != nil {
		writeAuthJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "Invalid login credentials",
		})
		return
	}

	response := authResponse{
		User:         resp.Session.User,
		AccessToken:  resp.Session.AccessToken,
		RefreshToken: resp.Session.RefreshToken,
		TokenType:    resp.Session.TokenType,
		ExpiresIn:    resp.Session.ExpiresIn,
		ExpiresAt:    resp.Session.ExpiresAt,
	}

	writeAuthJSON(w, http.StatusOK, response)
}

func writeAuthJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}
