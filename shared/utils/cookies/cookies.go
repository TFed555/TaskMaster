package cookies

import (
	"net/http"
	"time"
)

func SetCookies(w *http.ResponseWriter, token_name string, token_value string, exp_time time.Time) {
	http.SetCookie(*w, &http.Cookie{
		Name: token_name,
		HttpOnly: true,
		Value: token_value,
		Expires: exp_time,
	})
}
