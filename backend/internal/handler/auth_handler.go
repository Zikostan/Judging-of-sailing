package handler

import (
	"net/http"

	"github.com/user/judging-of-sailing/backend/internal/model"
	"github.com/user/judging-of-sailing/backend/internal/service"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

type RegisterRequest struct {
	Email     string `json:"email" example:"judge@example.com"`
	Password  string `json:"password" example:"securepassword123"`
	Role      string `json:"role" example:"judge" enums:"judge,secretary,admin"`
	LastName  string `json:"last_name" example:"Smith"`
	FirstName string `json:"first_name" example:"John"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"judge@example.com"`
	Password string `json:"password" example:"securepassword123"`
}

// Register creates a new user account
// @Summary      Register a new user
// @Description  Creates a new user (judge, secretary, or admin) with email and password. Returns the created user without the password hash.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body  handler.RegisterRequest  true  "Registration details"
// @Success      201  {object}  model.User
// @Failure      400  {object}  map[string]string  "invalid request body"
// @Failure      409  {object}  map[string]string  "email conflict or other error"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	user, err := h.auth.Register(r.Context(), req.Email, req.Password, model.UserRole(req.Role), req.LastName, req.FirstName)
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, user)
}

// Login authenticates a user and returns a JWT token
// @Summary      Login
// @Description  Authenticates a user with email and password. Returns a JWT token and user info.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body  handler.LoginRequest  true  "Login credentials"
// @Success      200  {object}  service.AuthResponse
// @Failure      400  {object}  map[string]string  "invalid request body"
// @Failure      401  {object}  map[string]string  "invalid email or password"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := h.auth.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}