package middleware

import (
	// "auth-service/internal/repository"
	// "auth-service/internal/models"
	cookies_func "auth-service/internal/pkg/cookies"
	"auth-service/internal/services"
	"log"
	"net/http"
	_ "strings"
	"time"

	// "strings"
	_ "time"

	_ "encoding/base64"
	_ "encoding/json"
	// "github.com/golang-jwt/jwt/v5"
	// "net/url"
)

type AuthMiddleware struct {
	authService services.AuthService
}

func NewAuthMiddleware(authService services.AuthService) *AuthMiddleware{
	return &AuthMiddleware{
		authService: authService,
	}
}

func (c *AuthMiddleware) SetAuthMiddleware(controller func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		refreshCookie, err := r.Cookie("refresh_token")
		if err != nil {	
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Access denied"))
			return
		}
		refreshToken:=refreshCookie.Value
		accessCookie, err := r.Cookie("access_token")
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Access denied"))
			return
		}
		accessToken:=accessCookie.Value

		log.Println("Called from middleware %s", refreshToken)
		log.Println("Called from middleware %s", accessToken)


		result, errMsg := c.authService.ValidateToken(refreshToken)

		if !result {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(errMsg))
			//ну и выкинуть из сессии типо
		}

		result, errMsg = c.authService.ValidateToken(accessToken)

		if !result {
			// w.WriteHeader(http.StatusForbidden)
			// w.Write([]byte(errMsg))
			success, newValue, err:=  c.authService.UpdateAccessToken(refreshToken)
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(errMsg))
			}
			if success {
				cookies_func.SetCookies(&w, "access_token", newValue, time.Now().Add(15 * time.Minute))
			}
		}

		controller(w, r)
	}


}