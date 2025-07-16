package middleware

import (
	"context"
	"log"
	"net/http"
	cookies_func "shared/utils/cookies"
	_"time"
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

func (c *AuthMiddleware) SetAuthMiddleware(controller http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	log.Print(r.Header)
	cookiesmas := ""
	if r.Header.Get("set-cookie") != "" {
		cookiesmas = r.Header.Get("set-cookie")
		log.Print("Get it from set-cookie")
	} else {
		cookiesmas = r.Header.Get("Cookie")
		log.Print("Get it from Cookie")
	}
		refreshToken, accessToken := cookies_func.ParseCookies(cookiesmas)

		log.Printf("Called from middleware %s\n", refreshToken)
		log.Printf("Called from middleware %s", accessToken)

		resultRefresh, _ := c.authService.ValidateToken(refreshToken)

		if !resultRefresh {
				w.WriteHeader(498)
				w.Write([]byte("Invalid token"))
				return
		}

		result, _ := c.authService.ValidateToken(accessToken)

		if !result {

			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unathorized"))
            return

			// success, newValue, err:=  c.authService.UpdateAccessToken(refreshToken)
			// log.Printf("Called from middleware access token: %t, %s, %v", success, newValue, err)
			// if err != nil {
			// 	w.WriteHeader(http.StatusForbidden)
			// 	w.Write([]byte(errMsg))
			// 	return
			// }
			// if success {
			// 	cookies_func.SetCookies(&w, "access_token", newValue, time.Now().Add(15 * time.Minute))
			// }
		}

		userID, err := c.authService.ParseUserId(refreshToken)
		if err != "" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(err))
            return
        }

    ctx := context.WithValue(r.Context(), UserIdKey, userID)

	log.Printf("CALLED FROM MDLWR %d \n", ctx.Value(UserIdKey).(uint))
	controller.ServeHTTP(w, r.WithContext(ctx))

	})
}
