package pathvar

import "net/http"

type PathVar string

func TeamSlug(r *http.Request) string {
	return r.PathValue("teamSlug")
}

func TournamentSlug(r *http.Request) string {
	return r.PathValue("tournamentSlug")
}

func GameSlug(r *http.Request) string {
	return r.PathValue("gameSlug")
}