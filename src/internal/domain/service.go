package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type GameService struct {
	storage IStorage
}

type UserService struct {
	userStorage IUser
}

func NewGameService(storage IStorage) *GameService {
	return &GameService{storage: storage}
}

func NewUserService(storage IUser) *UserService {
	return &UserService{userStorage: storage}
}

func (u *UserService) UpdateSession(login string, newUuid uuid.UUID) error {
	return u.userStorage.UpdateUserUUID(login, newUuid)
}

func (u *UserService) GetUserInfo(id uuid.UUID) (*UserDomain, error) {
	return u.userStorage.GetUserByUUID(id)
}

func (s *GameService) CreateGame(user UserDomain, gameType string) (*Game, error) {
	game := Game{
		Uuid:            uuid.New(),
		Matrix:          Field{},
		Win:             0,
		GameType:        gameType,
		Player1:         user.Login,
		Uuidplayer1:     user.Uuid,
		Player1Symbol:   "X",
		Player2Symbol:   "O",
		CurrentTurnUUID: user.Uuid,
		CreatedAt:       time.Now(),
	}

	if gameType == TypePVE {
		game.Player2 = "Computer"
		game.Uuidplayer2 = uuid.Nil
		game.Status = StatusTurnP1
	} else {
		game.Status = StatusWaiting
	}

	err := s.storage.SaveGameStorage(game)
	if err != nil {
		return nil, fmt.Errorf("ошибка сохранения игры: %w", err)
	}

	return &game, nil
}

func (s *GameService) JoinGame(gameID uuid.UUID, user UserDomain) (*Game, error) {
	game, err := s.storage.GetGameStorage(gameID)
	if err != nil {
		return nil, fmt.Errorf("игра не найдена: %w", err)
	}

	if game.Status != StatusWaiting {
		return nil, fmt.Errorf("игра уже началась или завершена")
	}

	if game.Uuidplayer1 == user.Uuid {
		return nil, fmt.Errorf("вы уже в этой игре")
	}

	game.Player2 = user.Login
	game.Uuidplayer2 = user.Uuid
	game.Status = StatusTurnP1     // Первый всегда ходит P1
	game.CurrentTurnUUID = game.Uuidplayer1 // Явно устанавливаем ход P1

	err = s.storage.SaveGameStorage(*game)
	return game, err
}

func (s *GameService) GetGame(id uuid.UUID) (*Game, error) {
	return s.storage.GetGameStorage(id)
}

func (s *GameService) NextMoveService(gameID uuid.UUID, userUUID uuid.UUID, inputMatrix Field) (*Game, error) {
	game, err := s.storage.GetGameStorage(gameID)
	if err != nil {
		return nil, fmt.Errorf("игра не найдена")
	}

	// 1. Проверка состояния игры
	if game.Status != StatusTurnP1 && game.Status != StatusTurnP2 {
		return nil, fmt.Errorf("сейчас нельзя ходить, статус игры: %s", game.Status)
	}

	// 2. Проверка, чей ход
	if game.CurrentTurnUUID != userUUID {
		return nil, fmt.Errorf("сейчас не ваш ход")
	}

	// 3. Проверка изменения матрицы (только одна ячейка)
	diffCount := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if game.Matrix[i][j] != inputMatrix[i][j] {
				// Проверяем, что изменили пустую ячейку на свой символ
				expectedVal := 1
				if game.Status == StatusTurnP2 {
					expectedVal = 2
				}
				if game.Matrix[i][j] != 0 || inputMatrix[i][j] != expectedVal {
					return nil, fmt.Errorf("недопустимый ход в ячейку [%d][%d]", i, j)
				}
				diffCount++
			}
		}
	}

	if diffCount != 1 {
		return nil, fmt.Errorf("нужно изменить ровно одну ячейку")
	}

	game.Matrix = inputMatrix

	// 4. Проверка окончания игры
	if s.updateGameState(game) {
		s.storage.SaveGameStorage(*game)
		return game, nil
	}

	// 5. Если PVE и игра не закончена — ходит бот
	if game.GameType == TypePVE && game.Status == StatusTurnP2 {
		botRow, botCol := BestMove(*game)
		game.Matrix[botRow][botCol] = 2
		s.updateGameState(game)
	}

	s.storage.SaveGameStorage(*game)
	return game, nil
}

// updateGameState проверяет победителя и обновляет статус и CurrentTurnUUID. 
// Возвращает true, если игра окончена.
func (s *GameService) updateGameState(g *Game) bool {
	win := CheckWin(*g)
	if win != 0 {
		g.Win = win
		switch win {
		case 1:
			g.Status = StatusWonP1
		case 2:
			g.Status = StatusWonP2
		case 3:
			g.Status = StatusDraw
		}
		g.CurrentTurnUUID = uuid.Nil
		return true
	}

	// Смена хода
	if g.Status == StatusTurnP1 {
		g.Status = StatusTurnP2
		g.CurrentTurnUUID = g.Uuidplayer2 // В PVE это будет uuid.Nil, что корректно для бота
	} else {
		g.Status = StatusTurnP1
		g.CurrentTurnUUID = g.Uuidplayer1
	}
	return false
}

func (u *UserService) LoginUser(user UserDomain) (uuid.UUID, error) {
	password, id, err := u.userStorage.CheckUserStorage(user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("пользователь не найден")
	}
	if password != user.Pass {
		return uuid.Nil, fmt.Errorf("неверный пароль")
	}
	return id, nil
}

func (u *UserService) RegistUser(user UserDomain) error {
	_, _, err := u.userStorage.CheckUserStorage(user)
	if err == nil {
		return fmt.Errorf("пользователь уже существует")
	}

	if user.Uuid == uuid.Nil {
		user.Uuid = uuid.New()
	}

	err = u.userStorage.SaveUserStorage(user)
	if err != nil {
		return fmt.Errorf("ошибка сохранения: %w", err)
	}
	return nil
}

func (s *GameService) DomainCheckGameSesion() ([]Game, error) {
	return s.storage.CheckGameSessions()
}

func (s *GameService) GetCompletedGames(userUUID uuid.UUID) ([]Game, error) {
	return s.storage.GetCompletedGamesByUser(userUUID)
}

func (s *GameService) GetLeaderboard(limit int) ([]LeaderboardEntry, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.storage.GetLeaderboard(limit)
}
