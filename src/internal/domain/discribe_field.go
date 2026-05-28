package domain

import (
	"time"

	"github.com/google/uuid"
)

type Field [3][3]int

// Состояния игры
const (
	StatusWaiting    = "waiting"      // Ожидание второго игрока
	StatusTurnP1     = "turn_p1"      // Ход первого игрока (X)
	StatusTurnP2     = "turn_p2"      // Ход второго игрока (O)
	StatusWonP1      = "won_p1"       // Победил первый игрок
	StatusWonP2      = "won_p2"       // Победил второй игрок
	StatusDraw       = "draw"         // Ничья
	StatusTerminated = "terminated"   // Игра прервана
)

// Типы игры
const (
	TypePVP = "PVP"
	TypePVE = "PVE"
)

type Game struct {
	Uuid            uuid.UUID
	Matrix          Field
	Win             int
	Status          string
	GameType        string
	Player1         string
	Uuidplayer1     uuid.UUID
	Player1Symbol   string // "X"
	Player2         string
	Uuidplayer2     uuid.UUID
	Player2Symbol   string // "O"
	CurrentTurnUUID uuid.UUID
	CreatedAt       time.Time
}

func NewUUID() uuid.UUID {
	return uuid.New()
}

type UserDomain struct {
	Login string
	Pass  string
	Uuid  uuid.UUID
}

type LeaderboardEntry struct {
	UserUUID uuid.UUID `json:"uuid"`
	Login    string    `json:"login"`
	WinRatio float64   `json:"win_ratio"`
}
