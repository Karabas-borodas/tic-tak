package datasource

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Field [3][3]int

// Game представляет внутреннюю структуру игры
type Game struct {
	Uuid            uuid.UUID
	Matrix          Field
	Win             int
	Status          string
	GameType        string
	Player1         string
	Uuidplayer1     uuid.UUID
	Player1Symbol   string
	Player2         string
	Uuidplayer2     uuid.UUID
	Player2Symbol   string
	CurrentTurnUUID uuid.UUID
	CreatedAt       time.Time
}

// GameStorageModel для маппинга в Postgres
type GameStorageModel struct {
	ID              uuid.UUID `db:"uuid"`
	Matrix          Field     `db:"matrix"`
	Win             int       `db:"win"`
	Status          string    `db:"status"`
	GameType        string    `db:"game_type"`
	Player1         string    `db:"player1"`
	Uuidplayer1     uuid.UUID `db:"uuidplayer1"`
	Player1Symbol   string    `db:"symbol_p1"`
	Player2         string    `db:"player2"`
	Uuidplayer2     uuid.UUID `db:"uuidplayer2"`
	Player2Symbol   string    `db:"symbol_p2"`
	CurrentTurnUUID uuid.UUID `db:"current_turn_uuid"`
	CreatedAt       time.Time `db:"created_at"`
}

type GameStorage struct {
	Conn *pgxpool.Pool
}

type UserDB struct {
	Login string    `db:"login"`
	Pass  string    `db:"password"`
	Uuid  uuid.UUID `db:"uuid"`
}
type UserStorage struct {
	Conn *pgxpool.Pool
}
