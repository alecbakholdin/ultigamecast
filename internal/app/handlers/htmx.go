package handlers

import "net/http"

func hxRedirect(w http.ResponseWriter, url string) {
	w.Header().Add("Hx-Redirect", url)
}

func hxRefresh(w http.ResponseWriter) {
	w.Header().Add("Hx-Refresh", "true")
}

func hxOpenModal(w http.ResponseWriter) {
	hxTrigger(w, "openmodal")
}

func hxCloseModal(w http.ResponseWriter) {
	hxTrigger(w, "closemodal")
}

func hxClearForm(w http.ResponseWriter) {
	hxTrigger(w, "clearform")
}

func hxTrigger(w http.ResponseWriter, event string) {
	w.Header().Add("Hx-Trigger", event)
}

func hxRetarget(w http.ResponseWriter, target, swap string) {
	w.Header().Add("Hx-Retarget", target)
	w.Header().Add("Hx-Reswap", swap)
}