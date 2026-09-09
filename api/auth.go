package api

import (
	"errors"
	"net/http"

	"github.com/N30A/trakt/auth"
	"github.com/N30A/trakt/validate"
)

type authHandler struct {
	authService *auth.AuthService
}

func newAuthHandler(authService *auth.AuthService) *authHandler {
	return &authHandler{authService: authService}
}

type initRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type initResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (h *authHandler) init(w http.ResponseWriter, r *http.Request) {
	request, err := decodeJSON[initRequest](r)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	email, err := validate.Email(request.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.authService.CreateInitialUser(r.Context(), email, request.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInitialUserExists) {
			http.Error(w, "initial user already exists", http.StatusConflict)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, initResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  string(user.Role),
	})
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *authHandler) login(w http.ResponseWriter, r *http.Request) {
	request, err := decodeJSON[loginRequest](r)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	email, err := validate.Email(request.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(r.Context(), email, request.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, loginResponse{Token: token})
}
