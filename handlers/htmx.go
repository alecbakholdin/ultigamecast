package handlers

import "net/http"

func hxRedirect(w http.ResponseWriter, url string) {
	w.Header().Add("Hx-Redirect", url)
}

func hxRefresh(w http.ResponseWriter) {
	w.Header().Add("Hx-Refresh", "true")
}