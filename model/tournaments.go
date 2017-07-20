package model

const (
	TOURNAMENT_STATUS_ANNOUNCED = 1
	TOURNAMENT_STATUS_FINISHED  = 2
)

type Tournament struct {
	TournamentId string  `json:"tournamentId"`
	Deposit      float64 `json:"deposit"`
	Status       int8    `json:"status"`
}

type Winner struct {
	PlayerId      string  `json:"playerId"`
	Prize         float64 `json:"prize"`
	IsParticipant bool
	Backers       []string
}
