package middleware

import (
	"net/http"
	"time"
)

func Delay(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 400)
		h.ServeHTTP(w, r)
	})
}