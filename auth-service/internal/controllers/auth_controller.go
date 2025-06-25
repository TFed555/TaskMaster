package controllers

import (
	"auth-service/internal/services"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type AuthController struct {
	authService *services.AuthService
}	

type RegRequest struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegResponse struct {
	UserID uint `json:"id"`
}

type AuthResponse struct {
	// UserID uint `json:"id"`
	// Email string `json:"email"`
	// Login string `json:"login"`
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}



func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req RegRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := c.authService.Register(req.Login, req.Email, req.Password)
	if err != nil {
		log.Printf("Registration error: %v", err)
		http.Error(w, fmt.Sprintf("Registration failed: %v", err), http.StatusInternalServerError)
		return
	}

	response := RegResponse{
		UserID: user.ID,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

var (
    errInvalidData = errors.New("email or password is incorrect")
)

func (c *AuthController) Authorize(w http.ResponseWriter, r *http.Request) {
	var req LogRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tokens, err, exists := c.authService.Authorize(req.Email, req.Password)
	if err != nil {
		log.Printf("%s", err)
	}

	if !exists {
		http.Error(w, errInvalidData.Error(), http.StatusBadRequest)
		return
	}

	// if err.Error() == "Invalid password" {
	// 	http.Error(w, errInvalidData.Error(), http.StatusBadRequest)
	// 	return
	// }

	response := AuthResponse{
		AccessToken: tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)

}

func (c *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {

}

// func (c *AuthController) JWTAuthMiddleware(next http.Handler) http.Handler {
// 	//дописать
// }