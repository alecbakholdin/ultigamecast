package models

type TournamentSummary struct {
	*Tournament
	Games []Game
	Data  []TournamentDatum
}
