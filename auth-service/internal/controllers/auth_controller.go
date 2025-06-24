package controllers

import (
	"auth-service/internal/services"
	"encoding/json"
	"log"
	"net/http"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
)

type AuthController struct {
	authService *services.AuthService
	tokenAuth *jwtauth.JWTAuth
}	

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}


type RegResponse struct {
	UserID uint `json:"id"`
}


func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
		tokenAuth: jwtauth.New("HS256", []byte("secret"), nil),
	}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := c.authService.Register(req.Email, req.Password)
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

func (c *AuthController) Authorize(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	answer, err := c.authService.Authorize(req.Email, req.Password)
	if err != nil {
		log.Printf("%w", err)
	}

	log.Printf("%b", answer)
}

// func (c *AuthController) JWTAuthMiddleware(next http.Handler) http.Handler {
// 	//дописать
// }