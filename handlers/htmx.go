package handlers

import "net/http"

func hxRedirect(w http.ResponseWriter, url string) {
	w.Header().Add("Hx-Redirect", url)
}