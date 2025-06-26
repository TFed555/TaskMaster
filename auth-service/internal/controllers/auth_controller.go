package controllers

import (
	_ "auth-service/internal/models"
	cookies_func "auth-service/internal/pkg/cookies"
	"auth-service/internal/services"
	"encoding/json"
	"errors"
	_ "fmt"
	"log"
	"net/http"
	"time"
)

var time_expires_refresh = time.Now().Add(7 * 24 * time.Hour)
var time_expires_access = time.Now().Add(15 * time.Minute)

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
	Email string `json:"email"`
	Login string `json:"login"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ErrorResponse struct {
	Status  uint  `json:"code"`
	Message string `json:"message"`
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

	user, tokens, err := c.authService.Register(req.Login, req.Email, req.Password)
	if err != nil {
		log.Printf("Registration error: %v", err)
		// http.Error(w, fmt.Sprintf("Registration failed: %v", err), http.StatusInternalServerError)
		w.WriteHeader(http.StatusBadRequest)
		errMsg:=ErrorResponse{
			Status: http.StatusInternalServerError,
			Message: "Registration failed",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	//todo: создание токенов и запись в куки
	
	cookies_func.SetCookies(&w, "refresh_token", tokens.RefreshToken, time_expires_refresh)
	cookies_func.SetCookies(&w, "access_token", tokens.AccessToken, time_expires_access)


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

	user, tokens, err, exists := c.authService.Authorize(req.Email, req.Password)
	if err != nil {
		log.Printf("%s", err)
	}

	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		errMsg:=ErrorResponse{
			Status: http.StatusInternalServerError,
			Message: "User does not exists",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	// if err.Error() == "Invalid password" {
	// 	http.Error(w, errInvalidData.Error(), http.StatusBadRequest)
	// 	return
	// }

	response := AuthResponse{
		Email: user.Email,
		Login: user.Login,
	}

	cookies_func.SetCookies(&w, "refresh_token", tokens.RefreshToken, time_expires_refresh)
	cookies_func.SetCookies(&w, "access_token", tokens.AccessToken, time_expires_access)


	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)

}

func (c *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	
	refresh_token, err:= r.Cookie("refresh_token")

	if err != nil {
		http.Error(w, "Unabled to update refresh token", http.StatusBadRequest)
		return
	}

	accesstoken, refreshtoken, err := c.authService.Refresh(refresh_token.Value)

	if err != nil {
		log.Print(err)
	}

	cookies_func.SetCookies(&w, "refresh_token", refreshtoken, time_expires_refresh)
	cookies_func.SetCookies(&w, "access_token", accesstoken, time_expires_access)


	response := RefreshResponse{
		AccessToken: accesstoken,
		RefreshToken: refreshtoken,
	}
	
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)

}

func (s *AuthController) Test(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Msg string `json:"msg"`
	}
	response:=Response {
		Msg: "ok",
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}