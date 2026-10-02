// SignupHandler godoc
// @Summary Register a new user
// @Description Creates a new user through Supabase Auth.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body authCredentialsRequest true "Signup credentials"
// @Success 201 {object} authResponse
// @Failure 400 {object} map[string]string
// @Failure 405 {object} map[string]string
// @Router /auth/signup [post]

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

type authUserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type authResponse struct {
	User         authUserResponse `json:"user"`
	AccessToken  string           `json:"access_token,omitempty"`
	RefreshToken string           `json:"refresh_token,omitempty"`
	TokenType    string           `json:"token_type,omitempty"`
	ExpiresIn    int              `json:"expires_in,omitempty"`
	ExpiresAt    int64            `json:"expires_at,omitempty"`
}

// SignupHandler godoc
// @Summary Register a new user
// @Description Creates a new user through Supabase Auth.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body authCredentialsRequest true "Signup credentials"
// @Success 201 {object} authResponse
// @Failure 400 {object} map[string]string
// @Failure 405 {object} map[string]string
// @Router /auth/signup [post]
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
		User: authUserResponse{
			ID:    resp.User.ID.String(),
			Email: resp.User.Email,
		},
		AccessToken:  resp.Session.AccessToken,
		RefreshToken: resp.Session.RefreshToken,
		TokenType:    resp.Session.TokenType,
		ExpiresIn:    resp.Session.ExpiresIn,
		ExpiresAt:    resp.Session.ExpiresAt,
	}

	writeAuthJSON(w, http.StatusCreated, response)
}

// LoginHandler godoc
// @Summary Log in a user
// @Description Authenticates a user through Supabase Auth and returns access and refresh tokens.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body authCredentialsRequest true "Login credentials"
// @Success 200 {object} authResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 405 {object} map[string]string
// @Router /auth/login [post]
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
		User: authUserResponse{
			ID:    resp.Session.User.ID.String(),
			Email: resp.Session.User.Email,
		},
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
