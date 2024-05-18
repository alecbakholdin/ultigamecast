package handlers

import (
	"net/http"
	view_home "ultigamecast/view/home"
)

type Home struct {

}

func NewHome() *Home {
	return &Home{}
}

func (h *Home) GetHome(w http.ResponseWriter, r *http.Request) {
	view_home.HomePage().Render(r.Context(), w)
}