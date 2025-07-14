package middleware

import (
	"context"
	"log"
	"net/http"
	cookies_func "shared/utils/cookies"
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
	cookiesmas := ""
	if r.Header.Get("set-cookie") != "" {
		cookiesmas = r.Header.Get("set-cookie")
	} else {
		cookiesmas = r.Header.Get("Cookie")
	}
	_, accessToken := cookies_func.ParseCookies(cookiesmas)

	log.Printf("Called from middleware %s", accessToken)

	result, errMsg := c.authService.ValidateToken(accessToken)

	if !result {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(errMsg))
    	// http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	userID, err := c.authService.ParseUserId(accessToken)
	if err != "" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err))
        return
    }

    ctx := context.WithValue(r.Context(), UserIdKey, userID)

	log.Printf("CALLED FROM MDLWR %d \n", ctx.Value(UserIdKey).(uint))
	controller(w, r.WithContext(ctx))

	}
}
