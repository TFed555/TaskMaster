package middleware

import (
	// "auth-service/internal/repository"
	// "auth-service/internal/models"
	"auth-service/internal/services"
	"log"
	"net/http"
	_"strings"

	// "strings"
	_ "time"

	_ "encoding/base64"
	_ "encoding/json"
	// "github.com/golang-jwt/jwt/v5"
	// "net/url"
)

// func AuthMiddleware(controller func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
// 	handlerfunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
// 		header := r.Header.Get("Authorization")
// 		if header == "" {
// 			w.Write([]byte("access denied"))
// 		}
// 		log.Printf("called from authmiddleware %s", strings.Split(header, " "))
// 		token := strings.Split(header, " ")[1]
// 		payload := strings.Split(token, ".")[1]
// 		dst := make([]byte, base64.RawURLEncoding.DecodedLen(len(payload)))
// 		n, err := base64.RawURLEncoding.Decode(dst, []byte(payload))
// 		payloadJson := dst[:n]
// 		if err != nil {
// 		log.Println("decode error:", err)
// 		return
// 		}

// 		type AccessToken struct {
// 			Sub uint `json:"sub"`
// 			Exp int64 `json:"exp"`
// 		}
// 		var accesstoken AccessToken
// 		err = json.Unmarshal(payloadJson, &accesstoken)
// 		if err != nil {
// 			log.Println("Can't unmarshal")
// 			return
// 		}
// 		log.Println(time.Unix(accesstoken.Exp, 0))
// 		if accesstoken.Exp < time.Now().Unix() {
// 			log.Print("Time ran out")
// 			return
// 		}
// 	})

// 	return handlerfunc
// }

func AuthMiddleware(controller func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		refreshCookie, err := r.Cookie("refresh_token")
		if err != nil {	
			w.Write([]byte("Access denied"))
			return
		}
		refreshToken:=refreshCookie.Value
		accessCookie, err := r.Cookie("access_token")
		if err != nil {
			w.Write([]byte("Access denied"))
			return
		}
		accessToken:=accessCookie.Value

		// payloadRefresh := strings.Split(refreshToken, " ")[1]
		// refreshToken = strings.Split(payloadRefresh, ".")[1]

		// payloadAccess := strings.Split(accessToken, " ")[1]
		// accessToken = strings.Split(payloadAccess, ".")[1]

		log.Println("Called from middleware %s", refreshToken)
		log.Println("Called from middleware %s", accessToken)
		var authService *services.AuthService

		result, errMsg := authService.ValidateToken(refreshToken)

		if !result {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(errMsg))
			//ну и выкинуть из сессии типо
		}

		result, errMsg = authService.ValidateToken(accessToken)

		if !result {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(errMsg))
		}

		controller(w, r)
	}


}