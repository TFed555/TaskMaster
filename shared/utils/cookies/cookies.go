package cookies

import (
	"log"
	"net/http"
	"strings"
	"time"
)

func ParseCookies(cookiesmas string) (rToken string, aToken string) {
	log.Print("cookies:", cookiesmas)
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

	return refreshToken, accessToken
} 

func SetCookies(w *http.ResponseWriter, token_name string, token_value string, exp_time time.Time) {
	http.SetCookie(*w, &http.Cookie{
		Name: token_name,
		HttpOnly: true,
		Value: token_value,
		Expires: exp_time,
	})
}
