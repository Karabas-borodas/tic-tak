package datasource

import (
	"anglefar/tic-tac-toe/internal/domain"
)

func StorageToDomain(g Game) *domain.Game {
	domainBase := domain.Game{
		Uuid:            g.Uuid,
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
	return &domainBase
}

func DomainToStorage(g domain.Game) Game {
	return Game{
		Uuid:            g.Uuid,
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

// DomainToModel переводит доменную сущность в модель БД.
func DomainToModel(g domain.Game) GameStorageModel {
	return GameStorageModel{
		ID:              g.Uuid,
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

// ModelToDomain переводит модель БД в доменную сущность.
func ModelToDomain(m GameStorageModel) *domain.Game {
	return &domain.Game{
		Uuid:            m.ID,
		Matrix:          domain.Field(m.Matrix),
		Win:             m.Win,
		Status:          m.Status,
		GameType:        m.GameType,
		Player1:         m.Player1,
		Uuidplayer1:     m.Uuidplayer1,
		Player1Symbol:   m.Player1Symbol,
		Player2:         m.Player2,
		Uuidplayer2:     m.Uuidplayer2,
		Player2Symbol:   m.Player2Symbol,
		CurrentTurnUUID: m.CurrentTurnUUID,
		CreatedAt:       m.CreatedAt,
	}
}

func DomainUserToStorage(u domain.UserDomain) *UserDB {
	return &UserDB{
		Login: u.Login,
		Pass:  u.Pass,
		Uuid:  u.Uuid,
	}
}
func StorageUserToDomain(u UserDB) *domain.UserDomain {
	return &domain.UserDomain{
		Login: u.Login,
		Pass:  u.Pass,
		Uuid:  u.Uuid,
	}
}
