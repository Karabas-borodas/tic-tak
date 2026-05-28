package api

import (
	domain "anglefar/tic-tac-toe/internal/domain"
	// web "anglefar/tic-tac-toe/internal/web"
)

type GameHandler struct {
	service *domain.GameService
}
