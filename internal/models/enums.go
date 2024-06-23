package models

type EventType string

const (
	// team events
	EventTypeGoal       EventType = "Goal"
	EventTypeBlock      EventType = "Block"
	EventTypeTurn       EventType = "Turn"
	EventTypeDrop       EventType = "Drop"
	EventTypeTimeout    EventType = "Timeout"
	EventTypeTimeoutEnd EventType = "Timeout End"

	// team only (no opponent equivalent)
	EventTypeAssist       EventType = "Assist"
	EventTypeTouch        EventType = "Touch"
	EventTypeSubIn        EventType = "Sub In"
	EventTypeSubOut       EventType = "Sub Out"
	EventTypeStartingLine EventType = "Starting Line"

	// opponent events
	EventTypeOpponentGoal       EventType = "Opponent Goal"
	EventTypeOpponentBlock      EventType = "Opponent Block"
	EventTypeOpponentTurn       EventType = "Opponent Turn"
	EventTypeOpponentDrop       EventType = "Opponent Drop"
	EventTypeOpponentTimeout    EventType = "Opponent Timeout"
	EventTypeOpponentTimeoutEnd EventType = "Opponent Timeout End"

	// neutral events
	EventTypePointStart EventType = "Point Start"
	EventTypeHalf       EventType = "Half"
	EventTypeGameEnd    EventType = "Game End"
	EventTypeHalfCap    EventType = "Half Cap"
	EventTypeSoftCap    EventType = "Soft Cap"
	EventTypeHardCap    EventType = "Hard Cap"
)

type GameLiveStatus string

const (
	GameLiveStatusPrePoint           GameLiveStatus = "Pre Point"
	GameLiveStatusTeamPossession     GameLiveStatus = "Team Possession"
	GameLiveStatusOpponentPossession GameLiveStatus = "Opponent Possession"
	GameLiveStatusTeamTimeout        GameLiveStatus = "Team Timeout"
	GameLiveStatusOpponentTimeout    GameLiveStatus = "Opponent Timeout"
	GameLiveStatusHalftime           GameLiveStatus = "Halftime"
	GameLiveStatusFinal              GameLiveStatus = "Final"
)

type GameScheduleStatus string

const (
	GameScheduleStatusLive      GameScheduleStatus = "Live"
	GameScheduleStatusScheduled GameScheduleStatus = "Scheduled"
	GameScheduleStatusFinal     GameScheduleStatus = "Final"
)

var GameScheduleStatusMap = map[string]GameScheduleStatus {
	"Live": GameScheduleStatusLive,
	"Scheduled": GameScheduleStatusScheduled,
	"Final": GameScheduleStatusFinal,
}

type GameJsonField string

const (
	GameJsonFieldScheduleStatus GameJsonField = "scheduleStatus"
)
