package pbmodels

type TournamentsWithGames struct {
	Tournament *Tournaments
	Games      []*Games
}
