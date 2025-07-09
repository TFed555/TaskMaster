package middleware

import (
	"context"
	"log"
	"net/http"
	"shared/utils/cookies"
	"strings"
	"time"
)

type AuthService interface{
	ValidateToken(token string) (bool, string)
	UpdateAccessToken(token string) (bool, string, error)
	ParseUserId(token string) (uint, string)
}

type AuthMiddleware struct {
	authService AuthService
}

func NewAuthMiddleware(authService AuthService) AuthMiddleware{
	return AuthMiddleware{
		authService: authService,
	}
}

type ContextKey string

const UserIdKey	ContextKey = "userID"

func (c *AuthMiddleware) SetAuthMiddleware(controller func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){

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
				accessToken = newEl[1]
			}
		}

		log.Printf("Called from middleware %s\n", refreshToken)
		log.Printf("Called from middleware %s", accessToken)


		result, errMsg := c.authService.ValidateToken(refreshToken)

		if !result {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(errMsg))
			return
		}

		result, errMsg = c.authService.ValidateToken(accessToken)

		if !result {
			success, newValue, err:=  c.authService.UpdateAccessToken(refreshToken)
			log.Printf("Called from middleware access token: %t, %s, %v", success, newValue, err)
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(errMsg))
				return
			}
			if success {
				cookies.SetCookies(&w, "access_token", newValue, time.Now().Add(15 * time.Minute))
			}
		}

		userID, err := c.authService.ParseUserId(refreshToken)
		if err != "" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(err))
            // http.Error(w, "Invalid token", http.StatusUnauthorized)

            return
        }

        ctx := context.WithValue(r.Context(), UserIdKey, userID)

		log.Printf("CALLED FROM MDLWR %d \n", ctx.Value(UserIdKey).(uint))

		controller(w, r.WithContext(ctx))
	}
}
