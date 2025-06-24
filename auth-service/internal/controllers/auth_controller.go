package controllers

import (
	"auth-service/internal/services"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/gin-gonic/gin"
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
		log.Printf("%s", err)
	}

	response := RegResponse{
		UserID: user.ID,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// func (c *AuthController) JWTAuthMiddleware(next http.Handler) http.Handler {
// 	//дописать
// }