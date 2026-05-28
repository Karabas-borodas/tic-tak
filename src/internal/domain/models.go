package domain

import (
	// "errors"
	"fmt"
	"log"
)

// NOTE:создание новой игры
// NOTE: функции для работы с игрой теперь в service.go

func validateField(g *Game, g2 *Game) (flag int, row int, col int) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if g.Matrix[i][j] != g2.Matrix[i][j] {
				row = i
				col = j
				flag++
			}
		}
	}
	return flag, row, col
}

// NOTE:возвращает структуру и ошибку
func NextMove(g *Game, player int, row int, col int) (game *Game, err error) {
	if g.Matrix[row][col] != 0 {
		log.Println("Этот ход уже был сделан ранее")
		err = fmt.Errorf("Этот ход уже был сделал")
		return g, err
	}

	// FIX: Сохраняем копию состояния до хода
	prevGame := *g
	game = &prevGame

	//FIX: нет валидации на входе rows cols
	if player == 1 {
		g.Matrix[row][col] = 1
	} else {
		g.Matrix[row][col] = 2
	}
	//NOTE: хз зачем эта фигня нужна по заданию
	flag, _, _ := validateField(g, game)
	if flag == 1 {
		return g, err
	} else {
		log.Println("какая то ошибка с ходом не могу сделать")
		err = fmt.Errorf("Этот ход уже был сделал")
		return g, err
	}
}

// //NOTE: gamgover
func (g *Game) IsGameOver() int {

	win := CheckWin(*g)
	return win
}
