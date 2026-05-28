package domain

func BestMove(g Game) (Row int, Col int) {
	bestScore := -1000
	bestRow := -1
	bestCol := -1

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if g.Matrix[i][j] == 0 {
				g.Matrix[i][j] = 2
				score := minMax(g, false, 0)
				g.Matrix[i][j] = 0

				if score > bestScore {
					bestScore = score
					bestRow, bestCol = i, j
				}
			}
		}
	}
	return bestRow, bestCol
}

// CheckWin возвращает:
//
//	1 — победил игрок 1
//	2 — победил игрок 2 (AI)
//	3 — ничья (все клетки заполнены, победителя нет)
//	0 — игра продолжается
func CheckWin(g Game) int {
	// Rows.
	for i := 0; i < 3; i++ {
		if g.Matrix[i][0] != 0 &&
			g.Matrix[i][0] == g.Matrix[i][1] &&
			g.Matrix[i][1] == g.Matrix[i][2] {
			return g.Matrix[i][0]
		}
	}

	// Columns.
	for j := 0; j < 3; j++ {
		if g.Matrix[0][j] != 0 &&
			g.Matrix[0][j] == g.Matrix[1][j] &&
			g.Matrix[1][j] == g.Matrix[2][j] {
			return g.Matrix[0][j]
		}
	}

	// Diagonals.
	if g.Matrix[0][0] != 0 &&
		g.Matrix[0][0] == g.Matrix[1][1] &&
		g.Matrix[1][1] == g.Matrix[2][2] {
		return g.Matrix[0][0]
	}
	if g.Matrix[0][2] != 0 &&
		g.Matrix[0][2] == g.Matrix[1][1] &&
		g.Matrix[1][1] == g.Matrix[2][0] {
		return g.Matrix[0][2]
	}

	// Ничья: все клетки заполнены, победителя нет.
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if g.Matrix[i][j] == 0 {
				return 0 // Есть свободные клетки — игра продолжается.
			}
		}
	}
	return 3 // Ничья.
}

func minMax(g Game, isMaximizing bool, depth int) int {
	winner := CheckWin(g)
	if winner == 2 {
		return 10 - depth
	}
	if winner == 1 {
		return depth - 10
	}
	if winner == 3 {
		return 0 // Ничья.
	}

	if isMaximizing {
		bestScore := -1000
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if g.Matrix[i][j] == 0 {
					g.Matrix[i][j] = 2
					score := minMax(g, false, depth+1)
					g.Matrix[i][j] = 0
					if score > bestScore {
						bestScore = score
					}
				}
			}
		}
		return bestScore
	}

	bestScore := 1000
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if g.Matrix[i][j] == 0 {
				g.Matrix[i][j] = 1
				score := minMax(g, true, depth+1)
				g.Matrix[i][j] = 0
				if score < bestScore {
					bestScore = score
				}
			}
		}
	}
	return bestScore
}
