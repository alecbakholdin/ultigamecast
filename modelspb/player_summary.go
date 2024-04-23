package modelspb

type PlayerSummary interface {
	GetPlayerName() string
	GetPlayerId() string
	GetPoints() int
	GetGoals() int
	GetAssists() int
	GetTurns() int
	GetDrops() int
}
