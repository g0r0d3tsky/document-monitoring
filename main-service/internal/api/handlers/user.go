package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	handlers "main-service/internal/api/handlers/gen"
	"main-service/internal/domain"
	"net/http"
)

type UserService interface {
	GenerateToken(ctx context.Context, login string, password string) (string, error)
	CreateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, username string) error
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	UpdateUser(ctx context.Context, username string, updatedUser *domain.User) error
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

type signInInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type signInResponse struct {
	Token string `json:"token"`
}

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := u.service.CreateUser(r.Context(), &user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

// LoginUser - Handles user login and token generation
func (u *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request, params handlers.LoginUserParams) {
	if params.Username == nil || params.Password == nil {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	token, err := u.service.GenerateToken(r.Context(), *params.Username, *params.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Login failed: %v", err), http.StatusUnauthorized)
		return
	}

	response := signInResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func LoginUserWrapper(apiHandler *APIHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params handlers.LoginUserParams

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid login data", http.StatusBadRequest)
			return
		}
		apiHandler.LoginUser(w, r, params)
	}
}

// TODO:
// LogoutUser - Handles user logout (e.g., invalidates token)
func (u *UserHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User logged out successfully"})
}

// DeleteUser - Deletes user by username
func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request, username string) {
	err := u.service.DeleteUser(r.Context(), username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}

// GetUserByName - Retrieves user details by username
func (u *UserHandler) GetUserByName(w http.ResponseWriter, r *http.Request, username string) {
	user, err := u.service.GetUserByUsername(r.Context(), username)
	if err != nil {
		http.Error(w, fmt.Sprintf("User not found: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateUser - Updates user details by username
func (u *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request, username string) {
	var updatedUser domain.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := u.service.UpdateUser(r.Context(), username, &updatedUser)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}
