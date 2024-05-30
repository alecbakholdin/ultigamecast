package handlers

import (
	"net/http"
	view_home "ultigamecast/web/view/home"
)

type Home struct {
}

func NewHome() *Home {
	return &Home{}
}

func (h *Home) GetHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	view_home.HomePage().Render(r.Context(), w)
}
