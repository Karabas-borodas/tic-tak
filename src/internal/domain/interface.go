package domain

import (
	// "anglefar/tic-tac-toe/internal/datasource"
	// "fmt"

	"github.com/google/uuid"
)

type IServices interface {
	NewGame() (*Game, error)
	NextMoveService(id uuid.UUID, player int, g *Game) (*Game, error)
	IsGameOver() int
}

type IStorage interface {
	SaveGameStorage(game Game) error
	GetGameStorage(id uuid.UUID) (*Game, error)
	CheckGameSessions() ([]Game, error)
	GetCompletedGamesByUser(userUUID uuid.UUID) ([]Game, error)
	GetLeaderboard(limit int) ([]LeaderboardEntry, error)
}

type IUser interface {
	SaveUserStorage(user UserDomain) error
	CheckUserStorage(user UserDomain) (string, uuid.UUID, error)
	UpdateUserUUID(login string, newUuid uuid.UUID) error
	GetUserByUUID(id uuid.UUID) (*UserDomain, error)
}
