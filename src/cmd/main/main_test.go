package main_test

import (
	domain "anglefar/tic-tac-toe/internal/domain"
	"testing"
	// "anglefar/tic-tac-toe/internal/uuid"
)

func TestMainGameState(t *testing.T) {
	game := domain.Game{
		Matrix: domain.Field{
			{1, 2, 1},
			{1, 2, 1},
			{2, 1, 2},
		},
		Uuid: domain.NewUUID(),
	}
	row, col := domain.BestMove(game)
	// Доска заполнена — нет свободных клеток, ожидаем (-1, -1).
	if row != -1 || col != -1 {
		t.Errorf("ожидалось (-1, -1) для заполненной доски, получено (%d, %d)", row, col)
	}
}

// TestBestMoveWithEmptyCell проверяет BestMove при наличии свободной клетки.
func TestBestMoveWithEmptyCell(t *testing.T) {
	game := domain.Game{
		Matrix: domain.Field{
			{1, 2, 1},
			{1, 0, 1},
			{2, 1, 2},
		},
		Uuid: domain.NewUUID(),
	}
	row, col := domain.BestMove(game)
	if row < 0 || row > 2 || col < 0 || col > 2 {
		t.Errorf("ожидались координаты 0–2, получено (%d, %d)", row, col)
	}
	if game.Matrix[row][col] != 0 {
		t.Errorf("клетка (%d, %d) занята", row, col)
	}
}
