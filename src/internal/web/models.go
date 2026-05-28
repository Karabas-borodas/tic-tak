package web

import (
	"time"

	"github.com/google/uuid"
)

type Field [3][3]int

type GameWeb struct {
	Id              uuid.UUID `json:"uuid"`
	Matrix          Field     `json:"matrix"`
	Win             int       `json:"win"`
	Status          string    `json:"status"`
	GameType        string    `json:"game_type"`
	Player1         string    `json:"player1"`
	Uuidplayer1     uuid.UUID `json:"uuidplayer1"`
	Player1Symbol   string    `json:"symbol_p1"`
	Player2         string    `json:"player2"`
	Uuidplayer2     uuid.UUID `json:"uuidplayer2"`
	Player2Symbol   string    `json:"symbol_p2"`
	CurrentTurnUUID uuid.UUID `json:"current_turn_uuid"`
	CreatedAt       time.Time `json:"created_at"`
}

type UserWeb struct {
	Login string    `json:"login"`
	Pass  string    `json:"password"`
	Uuid  uuid.UUID `json:"uuid"`
}

type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
