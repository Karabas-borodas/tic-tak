package domain

import (
	// discribefield "anglefar/tic-tac-toe/internal/domain/discribe_field"
	"testing"
)

func TestBestMoveWithEmptyCell(t *testing.T) {
	game := Game{
		Matrix: Field{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		Uuid: NewUUID(),
	}

	row, col := BestMove(game)
	if row < 0 || row > 3 || col < 0 || col > 3 {
		t.Errorf("ожидалось (0, 0) для заполненной доски, получено (%d, %d)", row, col)

	}
}
