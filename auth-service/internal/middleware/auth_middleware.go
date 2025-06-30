package middleware

import (
	// "auth-service/internal/repository"
	// "auth-service/internal/models"
	cookies_func "auth-service/internal/pkg/cookies"
	"auth-service/internal/services"
	"log"
	"net/http"
	"strings"
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
		// refreshCookie, err := r.Cookie("refresh_token")
		// if err != nil {	
		// 	w.WriteHeader(http.StatusForbidden)
		// 	w.Write([]byte("Access denied"))
		// 	return
		// }
		// refreshToken:=refreshCookie.Value
		// accessCookie, err := r.Cookie("access_token")
		// if err != nil {
		// 	w.WriteHeader(http.StatusForbidden)
		// 	w.Write([]byte("Access denied"))
		// 	return
		// }
		// accessToken:=accessCookie.Value

		cookiesmas := r.Header.Get("Cookie")
		log.Print("cookies:", cookiesmas)
		// var newString string
		var refreshToken string
		var accessToken string
		newString := strings.Split(cookiesmas, "; ")
		
		for _, el := range newString {
			newEl := strings.Split(el, "=")
			if newEl[0] == "refresh_token" {
				refreshToken = newEl[1]
			}
			if newEl[0] == "access_token" {
				refreshToken = newEl[1]
			}
		}

		log.Println("Called from middleware %s", refreshToken)
		log.Println("Called from middleware %s", accessToken)


		result, errMsg := c.authService.ValidateToken(refreshToken)

		if !result {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(errMsg))
			return
		}

		result, errMsg = c.authService.ValidateToken(accessToken)

		if !result {
			success, newValue, err:=  c.authService.UpdateAccessToken(refreshToken)
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(errMsg))
				return
			}
			if success {
				cookies_func.SetCookies(&w, "access_token", newValue, time.Now().Add(15 * time.Minute))
			}
		}

		controller(w, r)
	}
}