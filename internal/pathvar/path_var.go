package pathvar

import "net/http"

type PathVar string

func TeamSlug(r *http.Request) string {
	return r.PathValue("teamSlug")
}

func PlayerSlug(r *http.Request) string {
	return r.PathValue("playerSlug")
}

func TournamentSlug(r *http.Request) string {
	return r.PathValue("tournamentSlug")
}

func TournamentDatumSlug(r *http.Request) string {
	return r.PathValue("datumSlug")
}

func GameSlug(r *http.Request) string {
	return r.PathValue("gameSlug")
}