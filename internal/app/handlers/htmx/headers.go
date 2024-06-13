package htmx

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func HxRedirect(w http.ResponseWriter, url string) {
	w.Header().Add("Hx-Redirect", url)
}

func HxRefresh(w http.ResponseWriter) {
	w.Header().Add("Hx-Refresh", "true")
}

func HxOpenModal(w http.ResponseWriter) {
	HxTrigger(w, "openmodal", "")
}

func HxCloseModal(w http.ResponseWriter) {
	HxTrigger(w, "closemodal", "")
}

func HxCreateDatepicker(w http.ResponseWriter, id string) {
	HxTriggerAfterSettle(w, "createdatepicker", id)
}

func HxClearForm(w http.ResponseWriter) {
	HxTrigger(w, "clearform", "")
}

func HxTriggerAfterSettle(w http.ResponseWriter, event, data string) {
	hxTriggerKey(w, "Hx-Trigger-After-Settle", event, data)
}

func HxTriggerAfterSwap(w http.ResponseWriter, event, data string) {
	hxTriggerKey(w, "Hx-Trigger-After-Swap", event, data)
}

func HxTrigger(w http.ResponseWriter, event, data string) {
	hxTriggerKey(w, "Hx-Trigger", event, data)
}

func hxTriggerKey(w http.ResponseWriter, key, event, data string) {
	if h := w.Header().Get(key); h != "" {
		headers := make(map[string]string)
		err := json.Unmarshal([]byte(h), &headers)
		if err != nil {
			panic(fmt.Sprintf("error unamrshaling %s header: %s", key, err))
		}
		headers[event] = data
		bytes, err := json.Marshal(headers)
		if err != nil {
			panic(fmt.Sprintf("error marshalin %s header: %s", key, err))
		}
		w.Header().Set(key, string(bytes))
	} else {
		w.Header().Set(key, fmt.Sprintf(`{"%s": "%s"}`, event, data))
	}
}

func HxLocation(w http.ResponseWriter, url string) {
	w.Header().Add("Hx-Location", url)
}

func HxRetargetSwap(w http.ResponseWriter, target, swap string) {
	w.Header().Add("Hx-Retarget", target)
	w.Header().Add("Hx-Reswap", swap)
}

func HxRetarget(w http.ResponseWriter, target string) {
	w.Header().Add("Hx-Retarget", target)
}