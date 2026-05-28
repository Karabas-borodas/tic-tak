package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

type mockStorage struct {
	games []Game
}

func (m *mockStorage) SaveGameStorage(game Game) error {
	for i, g := range m.games {
		if g.Uuid == game.Uuid {
			m.games[i] = game
			return nil
		}
	}
	m.games = append(m.games, game)
	return nil
}

func (m *mockStorage) GetGameStorage(id uuid.UUID) (*Game, error) {
	for _, g := range m.games {
		if g.Uuid == id {
			return &g, nil
		}
	}
	return nil, nil
}

func (m *mockStorage) CheckGameSessions() ([]Game, error) {
	var result []Game
	for _, g := range m.games {
		if g.Status == StatusWaiting {
			result = append(result, g)
		}
	}
	return result, nil
}

func (m *mockStorage) GetCompletedGamesByUser(userUUID uuid.UUID) ([]Game, error) {
	var result []Game
	for _, g := range m.games {
		isPlayer := g.Uuidplayer1 == userUUID || g.Uuidplayer2 == userUUID
		isCompleted := (g.Status == StatusWonP1 && g.Uuidplayer1 == userUUID) ||
			(g.Status == StatusWonP2 && g.Uuidplayer2 == userUUID) ||
			(g.Status == StatusDraw && isPlayer)
		if isCompleted {
			result = append(result, g)
		}
	}
	return result, nil
}

func (m *mockStorage) GetLeaderboard(limit int) ([]LeaderboardEntry, error) {
	// Simple mock return
	return []LeaderboardEntry{
		{
			UserUUID: uuid.New(),
			Login:    "player1",
			WinRatio: 2.5,
		},
	}, nil
}

func TestGetCompletedGames(t *testing.T) {
	playerUUID := uuid.New()
	completedGame1 := Game{
		Uuid:        uuid.New(),
		Status:      StatusWonP1,
		Uuidplayer1: playerUUID,
		CreatedAt:   time.Now(),
	}
	completedGame2 := Game{
		Uuid:        uuid.New(),
		Status:      StatusDraw,
		Uuidplayer1: playerUUID,
		Uuidplayer2: uuid.New(),
		CreatedAt:   time.Now(),
	}
	activeGame := Game{
		Uuid:        uuid.New(),
		Status:      StatusTurnP1,
		Uuidplayer1: playerUUID,
		CreatedAt:   time.Now(),
	}

	storage := &mockStorage{
		games: []Game{completedGame1, completedGame2, activeGame},
	}
	service := NewGameService(storage)

	games, err := service.GetCompletedGames(playerUUID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(games) != 2 {
		t.Errorf("expected 2 completed games, got %d", len(games))
	}
}

func TestGetLeaderboard(t *testing.T) {
	storage := &mockStorage{}
	service := NewGameService(storage)

	entries, err := service.GetLeaderboard(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("expected 1 entry in leaderboard, got %d", len(entries))
	}

	if entries[0].Login != "player1" {
		t.Errorf("expected player1, got %s", entries[0].Login)
	}
}
