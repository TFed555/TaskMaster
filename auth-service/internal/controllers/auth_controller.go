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
	authService services.AuthService
}

type RegRequest struct {
	Name     string `json:"name"`
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
	UserID uint   `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ErrorResponse struct {
	Status  uint   `json:"code"`
	Message string `json:"message"`
}

//go:generate mockery --name=AuthService --dir=../services --output=./mocks --case=underscore
func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req RegRequest
	log.Println("DDDDDWWW")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Password == "" || req.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		errMsg := ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Registration failed",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	user, tokens, err := c.authService.Register(req.Name, req.Email, req.Password)
	if err != nil {
		log.Printf("Registration error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Registration failed",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	cookies_func.SetCookies(&w, "access_token", tokens.AccessToken, time_expires_access)
	cookies_func.SetCookies(&w, "refresh_token", tokens.RefreshToken, time_expires_refresh)

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
	log.Println("DDDDDWWW")

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
		errMsg := ErrorResponse{
			Status:  http.StatusInternalServerError,
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
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Login,
	}

	cookies_func.SetCookies(&w, "refresh_token", tokens.RefreshToken, time_expires_refresh)
	cookies_func.SetCookies(&w, "access_token", tokens.AccessToken, time_expires_access)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)

}

func (c *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {

	refresh_token, err := r.Cookie("refresh_token")

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
		AccessToken:  accesstoken,
		RefreshToken: refreshtoken,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)

}

func (c *AuthController) Test(w http.ResponseWriter, r *http.Request) {
	// type Response struct {
	// 	Msg string `json:"msg"`
	// }
	// response:=Response {
	// 	Msg: "ok",
	// }
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	// json.NewEncoder(w).Encode(res)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Test endpoint works!",
	})
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}
	//или User.Id
	success, err := c.authService.Logout(refreshCookie.Value)

	if !success && err != nil {
		log.Printf("Logout failed: %v", err)
		http.Error(w, "Logout failed", http.StatusInternalServerError)
		return
	}

	cookies_func.SetCookies(&w, "refresh_token", "", time.Unix(0, 0))
	cookies_func.SetCookies(&w, "access_token", "", time.Unix(0, 0))

	type Response struct {
		Msg string `json:"msg"`
	}

	response := Response{
		Msg: "ok",
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

func (c *AuthController) TestCookie(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Header.Get("set-cookie"))

	type Response struct {
		Msg string `json:"msg"`
	}

	response := Response{
		Msg: "ok",
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}
