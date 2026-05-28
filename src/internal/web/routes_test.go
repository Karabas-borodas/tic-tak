package web

import (
	"anglefar/tic-tac-toe/internal/domain"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

type mockStorage struct {
	games []domain.Game
}

func (m *mockStorage) SaveGameStorage(game domain.Game) error {
	m.games = append(m.games, game)
	return nil
}

func (m *mockStorage) GetGameStorage(id uuid.UUID) (*domain.Game, error) {
	for _, g := range m.games {
		if g.Uuid == id {
			return &g, nil
		}
	}
	return nil, nil
}

func (m *mockStorage) CheckGameSessions() ([]domain.Game, error) {
	return m.games, nil
}

func (m *mockStorage) GetCompletedGamesByUser(userUUID uuid.UUID) ([]domain.Game, error) {
	return m.games, nil
}

func (m *mockStorage) GetLeaderboard(limit int) ([]domain.LeaderboardEntry, error) {
	return []domain.LeaderboardEntry{
		{
			UserUUID: uuid.New(),
			Login:    "top_player",
			WinRatio: 3.5,
		},
	}, nil
}

type mockUserStorage struct{}

func (m *mockUserStorage) SaveUserStorage(user domain.UserDomain) error { return nil }
func (m *mockUserStorage) CheckUserStorage(user domain.UserDomain) (string, uuid.UUID, error) {
	return "pass", user.Uuid, nil
}
func (m *mockUserStorage) UpdateUserUUID(login string, newUuid uuid.UUID) error { return nil }
func (m *mockUserStorage) GetUserByUUID(id uuid.UUID) (*domain.UserDomain, error) {
	return &domain.UserDomain{Login: "user1", Pass: "pass", Uuid: id}, nil
}

func TestHistoryEndpoint(t *testing.T) {
	storage := &mockStorage{
		games: []domain.Game{
			{
				Uuid:        uuid.New(),
				Status:      domain.StatusWonP1,
				Uuidplayer1: uuid.New(),
				CreatedAt:   time.Now(),
			},
		},
	}
	userStorage := &mockUserStorage{}

	gameService := domain.NewGameService(storage)
	userService := domain.NewUserService(userStorage)
	auth := NewUserAuthenticator(userService)

	req := httptest.NewRequest("GET", "/games/history", nil)
	w := httptest.NewRecorder()

	userUUID := uuid.New()
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: userUUID.String(),
	})
	ctx := context.WithValue(req.Context(), UserContextKey, userUUID)
	req = req.WithContext(ctx)

	handler := auth.Middleware(func(w http.ResponseWriter, r *http.Request) {
		completedGames, err := gameService.GetCompletedGames(userUUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var webGames []GameWeb
		for _, g := range completedGames {
			webGames = append(webGames, DomainToWeb(g))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(webGames)
	})

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response []GameWeb
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response) != 1 {
		t.Errorf("expected 1 completed game, got %d", len(response))
	}
}

func TestLeaderboardEndpoint(t *testing.T) {
	storage := &mockStorage{}
	userStorage := &mockUserStorage{}

	gameService := domain.NewGameService(storage)
	userService := domain.NewUserService(userStorage)
	auth := NewUserAuthenticator(userService)

	req := httptest.NewRequest("GET", "/leaderboard?n=5", nil)
	w := httptest.NewRecorder()

	userUUID := uuid.New()
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: userUUID.String(),
	})
	ctx := context.WithValue(req.Context(), UserContextKey, userUUID)
	req = req.WithContext(ctx)

	handler := auth.Middleware(func(w http.ResponseWriter, r *http.Request) {
		limitStr := r.URL.Query().Get("n")
		limit := 10
		if limitStr == "5" {
			limit = 5
		}

		entries, err := gameService.GetLeaderboard(limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entries)
	})

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response []domain.LeaderboardEntry
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response) != 1 {
		t.Errorf("expected 1 entry, got %d", len(response))
	}

	if response[0].Login != "top_player" {
		t.Errorf("expected login top_player, got %s", response[0].Login)
	}
}
