package controllers

import (
	"auth-service/internal/pkg/domain_models"
	"auth-service/internal/pkg/responses"
	"auth-service/internal/services"
	"encoding/json"
	_ "fmt"
	"log"
	"net/http"
	"shared/middleware"
	cookies_func "shared/utils/cookies"
	"time"
)


type AuthController struct {
	authService services.AuthService
}

//go:generate mockery --name=AuthService --dir=../services --output=./mocks --case=underscore
func NewAuthController(authService services.AuthService) AuthController {
	return AuthController{
		authService: authService,
	}
}

func (c AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req responses.RegRequest
	log.Println("DDDDDWWW")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Password == "" || req.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		errMsg := responses.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Registration failed",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}
	
	params := domain_models.RegisterParams{
		Name: req.Name,
		Email: req.Email,
		Password: req.Password,
	}
	user, tokens, err := c.authService.Register(params)
	if err != nil {
		log.Printf("Registration error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := responses.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Registration failed",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	var time_expires_refresh = time.Now().Add(30 * 24 * time.Hour)
	var time_expires_access = time.Now().Add(15 * time.Minute)

	cookies_func.SetCookies(&w, "access_token", tokens.AccessToken, time_expires_access)
	cookies_func.SetCookies(&w, "refresh_token", tokens.RefreshToken, time_expires_refresh)

	response := responses.RegResponse{
		UserID: user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}


func (c AuthController) Authorize(w http.ResponseWriter, r *http.Request) {
	var req responses.LogRequest
	log.Println("DDDDDWWW")
	log.Print(r.Header.Get("set-cookie"))

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	params := domain_models.AuthorizeParams{
		Email: req.Email,
		Password: req.Password,
	}
	log.Print(params)
	user, tokens, err := c.authService.Authorize(params)

	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMessage := "Internal Server Error"
		switch err.Error() {
			case "User not found":
				statusCode = http.StatusNotFound
				errorMessage = "User not found"
			
			case "Invalid password":
				statusCode = http.StatusBadRequest
				errorMessage = "Invalid password or email"
		}
		w.WriteHeader(statusCode)
		errMsg := responses.ErrorResponse{
			Status:  statusCode,
			Message: errorMessage,
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	response := responses.AuthResponse{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Login,
		ImgURL: user.ImgPath,
	}

	log.Printf("Called from auth_controller, tokens: %s, %s", tokens.AccessToken, tokens.RefreshToken)

	var time_expires_refresh = time.Now().Add(30 * 24 * time.Hour)
	var time_expires_access = time.Now().Add(15 * time.Minute)

	cookies_func.SetCookies(&w, "refresh_token", tokens.RefreshToken, time_expires_refresh)
	cookies_func.SetCookies(&w, "access_token", tokens.AccessToken, time_expires_access)


	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)

}

func (c AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	cookiesmas := r.Header.Get("set-cookie")
	// cookiesmas := r.Header.Get("set-cookie")

	refreshToken, _ := cookies_func.ParseCookies(cookiesmas)

	result, errMsg := c.authService.ValidateToken(refreshToken)

	if !result {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(errMsg))
			return
	}

	accesstoken, refreshtoken, err := c.authService.Refresh(refreshToken)

	if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(errMsg))
			return
	}

	var time_expires_refresh = time.Now().Add(30 * 24 * time.Hour)
	var time_expires_access = time.Now().Add(15 * time.Minute)

	cookies_func.SetCookies(&w, "refresh_token", refreshtoken, time_expires_refresh)
	cookies_func.SetCookies(&w, "access_token", accesstoken, time_expires_access)

	response := responses.RefreshResponse{
		AccessToken:  accesstoken,
		RefreshToken: refreshtoken,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)

}

func (c AuthController) Test(w http.ResponseWriter, r *http.Request) {
	// type Response struct {
	// 	Msg string `json:"msg"`
	// }
	// response:=Response {
	// 	Msg: "ok",
	// }
	userID := r.Context().Value(middleware.UserIdKey).(uint)
    log.Printf("Controller received userID: %v", userID)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	// json.NewEncoder(w).Encode(res)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Test endpoint works!",
	})
}

func (c AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	// refreshCookie, err := r.Cookie("refresh_token")
	// if err != nil {
	// 	http.Error(w, "Missing token", http.StatusBadRequest)
	// 	return
	// }
	// //или User.Id
	// success, err := c.authService.Logout(refreshCookie.Value)
	ctx := r.Context()
	userID, ok:= ctx.Value(middleware.UserIdKey).(uint)
	log.Print("UserID:", userID)
	if !ok {
	    http.Error(w, "Unauthorized", http.StatusUnauthorized)
	    return
	}
	log.Printf("Controller received userID: %v", userID)

	success, err := c.authService.Logout(userID)

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

func (c AuthController) TestCookie(w http.ResponseWriter, r *http.Request) {
	log.Printf("TestCookie controller: %s", r.Header.Get("set-cookie"))

	type Response struct {
		Cookies	string `json:"cookies"`
	}

	response := Response{
		Cookies: r.Header.Get("set-cookie"),
		// Cookies: r.Header.Get("set-cookie")
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

func (c AuthController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok:= ctx.Value(middleware.UserIdKey).(uint)
	log.Print("UserID:", userID)
	if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

	log.Printf("Controller received userID: %v", userID)
	var req responses.UpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	params := domain_models.UpdateUserParams{
		ID: userID,
	}

	if req.Email != "" {
		params.Email = &req.Email
	}

	if req.Name != "" {
		params.Name = &req.Name
	}

	if req.Password != "" {
		params.Password = &req.Password
	}

	answer, err := c.authService.UpdateUser(params)

	if err != nil {
		log.Printf("%s", err)
		return
	}

	if answer < 0 {
		w.WriteHeader(http.StatusBadRequest)
		errMsg := responses.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Can't update user info",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(answer)

}

func (c AuthController) DeleteUser(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	userID, ok:= ctx.Value(middleware.UserIdKey).(uint)
	log.Print("UserID:", userID)
	if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

	log.Printf("Controller received userID: %v", userID)
	var req responses.UpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	answer, err := c.authService.DeleteUser(userID)

	if err != nil {
		log.Printf("%s", err)
		return
	}

	if !answer {
		w.WriteHeader(http.StatusBadRequest)
		errMsg := responses.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Can't update user info",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(answer)
}
