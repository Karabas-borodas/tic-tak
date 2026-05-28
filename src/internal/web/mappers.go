package web

import (
	domain "anglefar/tic-tac-toe/internal/domain"
)

func WebToDomain(g GameWeb) *domain.Game {
	return &domain.Game{
		Uuid:            g.Id,
		Matrix:          domain.Field(g.Matrix),
		Win:             g.Win,
		Status:          g.Status,
		GameType:        g.GameType,
		Player1:         g.Player1,
		Uuidplayer1:     g.Uuidplayer1,
		Player1Symbol:   g.Player1Symbol,
		Player2:         g.Player2,
		Uuidplayer2:     g.Uuidplayer2,
		Player2Symbol:   g.Player2Symbol,
		CurrentTurnUUID: g.CurrentTurnUUID,
		CreatedAt:       g.CreatedAt,
	}
}

func DomainToWeb(g domain.Game) GameWeb {
	return GameWeb{
		Id:              g.Uuid,
		Matrix:          Field(g.Matrix),
		Win:             g.Win,
		Status:          g.Status,
		GameType:        g.GameType,
		Player1:         g.Player1,
		Uuidplayer1:     g.Uuidplayer1,
		Player1Symbol:   g.Player1Symbol,
		Player2:         g.Player2,
		Uuidplayer2:     g.Uuidplayer2,
		Player2Symbol:   g.Player2Symbol,
		CurrentTurnUUID: g.CurrentTurnUUID,
		CreatedAt:       g.CreatedAt,
	}
}

func WebToDomainUser(u UserWeb) *domain.UserDomain {
	return &domain.UserDomain{
		Login: u.Login,
		Pass:  u.Pass,
		Uuid:  u.Uuid,
	}
}
func DomainToWebUser(u domain.UserDomain) UserWeb {
	return UserWeb{
		Login: u.Login,
		Pass:  u.Pass,
		Uuid:  u.Uuid,
	}
}
