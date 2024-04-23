package handlers

import "ultigamecast/repository"

type Games struct {
	TeamRepo *repository.Team
	TouramentRepo *repository.Tournament

}