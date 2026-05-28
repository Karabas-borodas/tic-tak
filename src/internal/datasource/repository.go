package datasource

import (
	"anglefar/tic-tac-toe/internal/domain"
	"context"
	"fmt"

	// "os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const urldb = "postgres://myuser:mypassword@127.0.0.1:5432/ticdb"

func NewGameStorage() (*GameStorage, error) {
	cfg, err := pgxpool.ParseConfig(urldb)
	if err != nil {
		return nil, fmt.Errorf("parse pool config: %w", err)
	}

	cfg.MaxConns = 20
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.HealthCheckPeriod = time.Minute

	conn, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		// fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, fmt.Errorf("error creating database: %w", err)
	}
	if err := (conn.Ping(context.Background())); err != nil {
		conn.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	fmt.Println("✅ Successfully connected to PostgreSQL")
	gs := &GameStorage{Conn: conn}

	if err := gs.createDB(); err != nil {
		return nil, fmt.Errorf("error creating database: %w", err)
	}
	return gs, nil
}

func (g *GameStorage) createDB() error {
	const sql = `
	CREATE TABLE IF NOT EXISTS games (
		id SERIAL PRIMARY KEY,
		uuid UUID UNIQUE NOT NULL,
		matrix JSONB NOT NULL,
		win INT DEFAULT 0,
		status TEXT NOT NULL DEFAULT 'waiting',
		game_type TEXT NOT NULL DEFAULT 'PVE',
		player1 TEXT DEFAULT '',
		uuidplayer1 UUID NOT NULL,
		symbol_p1 TEXT DEFAULT 'X',
		player2 TEXT DEFAULT '',
		uuidplayer2 UUID,
		symbol_p2 TEXT DEFAULT 'O',
		current_turn_uuid UUID,
		created_at TIMESTAMPTZ DEFAULT now()
	);
	
	-- Миграция для существующих таблиц
	ALTER TABLE games ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'waiting';
	ALTER TABLE games ADD COLUMN IF NOT EXISTS game_type TEXT NOT NULL DEFAULT 'PVE';
	ALTER TABLE games ADD COLUMN IF NOT EXISTS symbol_p1 TEXT DEFAULT 'X';
	ALTER TABLE games ADD COLUMN IF NOT EXISTS symbol_p2 TEXT DEFAULT 'O';
	ALTER TABLE games ADD COLUMN IF NOT EXISTS current_turn_uuid UUID;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := g.Conn.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("create table/migration execution: %w", err)
	}
	return nil
}

// NOTE: create user

func (u *UserStorage) SaveUserStorage(user domain.UserDomain) error {
	dbUser := DomainUserToStorage(user)
	dbUser.Uuid = uuid.New()
	return u.SaveUser(*dbUser)
}

func NewUserStorage() (*UserStorage, error) {
	cfg, err := pgxpool.ParseConfig(urldb)
	if err != nil {
		return nil, fmt.Errorf("parse pool config: %w", err)
	}

	cfg.MaxConns = 20
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.HealthCheckPeriod = time.Minute

	conn, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating database: %w", err)
	}

	if err := (conn.Ping(context.Background())); err != nil {
		conn.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	fmt.Println("✅ Successfully connected to PostgreSQL")
	gs := &UserStorage{Conn: conn}

	if err := gs.createUserDB(); err != nil {
		return nil, fmt.Errorf("error creating database: %w", err)
	}

	return gs, nil
}

// NOTE: sreate user data base
func (u *UserStorage) createUserDB() error {
	const sql = `
	CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		login TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		player TEXT DEFAULT '',
		uuid UUID UNIQUE NOT NULL,
		created_at TIMESTAMPTZ DEFAULT now()
	);`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := u.Conn.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("create table execution: %w", err)
	}
	return nil
}

func (u *UserStorage) SaveUser(user UserDB) error {
	const sql = `
	INSERT INTO users (login, password, uuid) VALUES ($1, $2, $3);`
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := u.Conn.Exec(ctx, sql, user.Login, user.Pass, user.Uuid)
	if err != nil {
		return fmt.Errorf("error saving user to table: %w", err)
	}
	return nil
}

func (u *UserStorage) UpdateUserUUID(login string, newUuid uuid.UUID) error {
	const sql = `UPDATE users SET uuid = $1 WHERE login = $2;`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := u.Conn.Exec(ctx, sql, newUuid, login)
	return err
}

// NOTE: check user password db
func (u *UserStorage) CheckUserStorage(user domain.UserDomain) (string, uuid.UUID, error) {
	dbUser := DomainUserToStorage(user)
	const sql = `
	SELECT password, uuid FROM users WHERE login=$1;`
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var password string
	var id uuid.UUID
	err := u.Conn.QueryRow(ctx, sql, dbUser.Login).Scan(&password, &id)
	if err != nil {
		return "", uuid.Nil, fmt.Errorf("error querying user: %w", err)
	}
	return password, id, nil
}

func (u *UserStorage) GetUserByUUID(id uuid.UUID) (*domain.UserDomain, error) {
	const sql = `SELECT login, password, uuid FROM users WHERE uuid = $1;`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user UserDB
	err := u.Conn.QueryRow(ctx, sql, id).Scan(&user.Login, &user.Pass, &user.Uuid)
	if err != nil {
		return nil, err
	}
	return StorageUserToDomain(user), nil
}

func (g *GameStorage) SaveGameStorage(game domain.Game) error {
	gameForDB := DomainToModel(game)

	const sql = `
	INSERT INTO games (uuid, matrix, win, status, game_type, player1, uuidplayer1, symbol_p1, player2, uuidplayer2, symbol_p2, current_turn_uuid, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	ON CONFLICT (uuid) DO UPDATE SET
		matrix = EXCLUDED.matrix,
		win = EXCLUDED.win,
		status = EXCLUDED.status,
		game_type = EXCLUDED.game_type,
		player1 = EXCLUDED.player1,
		uuidplayer1 = EXCLUDED.uuidplayer1,
		symbol_p1 = EXCLUDED.symbol_p1,
		player2 = EXCLUDED.player2,
		uuidplayer2 = EXCLUDED.uuidplayer2,
		symbol_p2 = EXCLUDED.symbol_p2,
		current_turn_uuid = EXCLUDED.current_turn_uuid;`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := g.Conn.Exec(ctx, sql,
		gameForDB.ID, gameForDB.Matrix, gameForDB.Win, gameForDB.Status, gameForDB.GameType,
		gameForDB.Player1, gameForDB.Uuidplayer1, gameForDB.Player1Symbol,
		gameForDB.Player2, gameForDB.Uuidplayer2, gameForDB.Player2Symbol,
		gameForDB.CurrentTurnUUID, gameForDB.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("error saving game to table: %w", err)
	}

	return nil
}

func (g *GameStorage) GetGameStorage(id uuid.UUID) (*domain.Game, error) {
	const sql = `SELECT uuid, matrix, win, status, game_type, player1, uuidplayer1, symbol_p1, player2, uuidplayer2, symbol_p2, current_turn_uuid, created_at 
	             FROM games WHERE uuid = $1;`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var m GameStorageModel
	err := g.Conn.QueryRow(ctx, sql, id).Scan(
		&m.ID, &m.Matrix, &m.Win, &m.Status, &m.GameType,
		&m.Player1, &m.Uuidplayer1, &m.Player1Symbol,
		&m.Player2, &m.Uuidplayer2, &m.Player2Symbol,
		&m.CurrentTurnUUID, &m.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("game %s not found", id)
		}
		return nil, fmt.Errorf("error querying game from table: %w", err)
	}

	return ModelToDomain(m), nil
}

// NOTE: получение всех игр

func (g *GameStorage) CheckGameSessions() ([]domain.Game, error) {
	var sliceActiveGames []Game

	// Выбираем ВСЕ необходимые колонки, которые хотим отсканировать
	// Выбираем только те игры, которые ожидают второго игрока
	const sql = `SELECT uuid, matrix, win, status, game_type, player1, uuidplayer1, symbol_p1, player2, uuidplayer2, symbol_p2, current_turn_uuid, created_at 
	             FROM games WHERE status='waiting'`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Используем Query для SELECT
	rows, err := g.Conn.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("error querying game sessions: %w", err)
	}
	defer rows.Close() // Обязательно закрываем rows

	for rows.Next() {
		var activeGame Game
		// Количество аргументов в Scan совпадает с количеством колонок в SELECT
		err := rows.Scan(
			&activeGame.Uuid,
			&activeGame.Matrix,
			&activeGame.Win,
			&activeGame.Status,
			&activeGame.GameType,
			&activeGame.Player1,
			&activeGame.Uuidplayer1,
			&activeGame.Player1Symbol,
			&activeGame.Player2,
			&activeGame.Uuidplayer2,
			&activeGame.Player2Symbol,
			&activeGame.CurrentTurnUUID,
			&activeGame.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning game row: %w", err)
		}
		// Добавляем найденную игру в слайс
		sliceActiveGames = append(sliceActiveGames, activeGame)
	}

	// Проверяем, не возникло ли ошибок при итерации
	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration error: %w", rows.Err())
	}
	var SlieceActiveGameDomain []domain.Game
	for _, Session := range sliceActiveGames {
		SlieceActiveGameDomain = append(SlieceActiveGameDomain, *StorageToDomain(Session))
	}

	return SlieceActiveGameDomain, nil
}

func (g *GameStorage) GetCompletedGamesByUser(userUUID uuid.UUID) ([]domain.Game, error) {
	var completedGames []Game

	const sql = `
		SELECT uuid, matrix, win, status, game_type, player1, uuidplayer1, symbol_p1, player2, uuidplayer2, symbol_p2, current_turn_uuid, created_at
		FROM games
		WHERE ((status = 'won_p1' AND uuidplayer1 = $1)
		   OR (status = 'won_p2' AND uuidplayer2 = $1)
		   OR (status = 'draw' AND (uuidplayer1 = $1 OR uuidplayer2 = $1)))
		ORDER BY created_at DESC;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := g.Conn.Query(ctx, sql, userUUID)
	if err != nil {
		return nil, fmt.Errorf("error querying completed games: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cg Game
		err := rows.Scan(
			&cg.Uuid,
			&cg.Matrix,
			&cg.Win,
			&cg.Status,
			&cg.GameType,
			&cg.Player1,
			&cg.Uuidplayer1,
			&cg.Player1Symbol,
			&cg.Player2,
			&cg.Uuidplayer2,
			&cg.Player2Symbol,
			&cg.CurrentTurnUUID,
			&cg.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning completed game: %w", err)
		}
		completedGames = append(completedGames, cg)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration error: %w", rows.Err())
	}

	var completedDomain []domain.Game
	for _, cg := range completedGames {
		completedDomain = append(completedDomain, *StorageToDomain(cg))
	}

	return completedDomain, nil
}

func (g *GameStorage) GetLeaderboard(limit int) ([]domain.LeaderboardEntry, error) {
	var entries []domain.LeaderboardEntry

	const sql = `
		SELECT 
			u.uuid,
			u.login,
			COALESCE(
				CAST(COUNT(CASE WHEN (g.status = 'won_p1' AND g.uuidplayer1 = u.uuid) OR (g.status = 'won_p2' AND g.uuidplayer2 = u.uuid) THEN 1 END) AS FLOAT) 
				/ NULLIF(COUNT(CASE WHEN (g.status = 'won_p1' AND g.uuidplayer2 = u.uuid) OR (g.status = 'won_p2' AND g.uuidplayer1 = u.uuid) OR g.status = 'draw' THEN 1 END), 0),
				COUNT(CASE WHEN (g.status = 'won_p1' AND g.uuidplayer1 = u.uuid) OR (g.status = 'won_p2' AND g.uuidplayer2 = u.uuid) THEN 1 END)::FLOAT
			) AS win_ratio
		FROM users u
		LEFT JOIN games g ON u.uuid = g.uuidplayer1 OR u.uuid = g.uuidplayer2
		WHERE g.status IN ('won_p1', 'won_p2', 'draw') OR g.uuid IS NULL
		GROUP BY u.uuid, u.login
		ORDER BY win_ratio DESC
		LIMIT $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := g.Conn.Query(ctx, sql, limit)
	if err != nil {
		return nil, fmt.Errorf("error querying leaderboard: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var entry domain.LeaderboardEntry
		err := rows.Scan(&entry.UserUUID, &entry.Login, &entry.WinRatio)
		if err != nil {
			return nil, fmt.Errorf("error scanning leaderboard entry: %w", err)
		}
		entries = append(entries, entry)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration error: %w", rows.Err())
	}

	return entries, nil
}
