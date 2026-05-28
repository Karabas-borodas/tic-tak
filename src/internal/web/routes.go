package web

import (
	domain "anglefar/tic-tac-toe/internal/domain"

	"context"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/fx"
)

const httpPort = ":3002"

//go:embed index.html
var indexHTML []byte

// StartGame теперь принимает fx.Lifecycle, чтобы fx знал, как управлять сервером
func StartGame(lc fx.Lifecycle, service *domain.GameService, user *domain.UserService, auth *UserAuthenticator) {
	mux := http.NewServeMux()

	// GET / — отдаём HTML-страницу
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
	})

	// POST /register — регистрация пользователя
	mux.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		var req SignUpRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "некорректный JSON", http.StatusBadRequest)
			return
		}
		if req.Login == "" || req.Password == "" {
			http.Error(w, "Логин и пароль обязательны", http.StatusBadRequest)
			return
		}
		
		domainUser := &domain.UserDomain{
			Login: req.Login,
			Pass:  req.Password,
		}
		err := user.RegistUser(*domainUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "пользователь зарегистрирован")
	})

	// POST /login — авторизация (Basic Auth Base64)
	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
			http.Error(w, "требуется Authorization header", http.StatusUnauthorized)
			return
		}

		payload, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
		if err != nil {
			http.Error(w, "ошибка декодирования base64", http.StatusBadRequest)
			return
		}

		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "неверный формат заголовка", http.StatusBadRequest)
			return
		}

		login, pass := pair[0], pair[1]
		userUUID, err := user.LoginUser(domain.UserDomain{Login: login, Pass: pass})
		if err != nil {
			http.Error(w, "неверные учетные данные", http.StatusUnauthorized)
			return
		}

		// Сохраняем UUID в куки (для авторизации через middleware)
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    userUUID.String(),
			Path:     "/",
			HttpOnly: true,
			MaxAge:   3600,
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"uuid": userUUID.String()})
	})

	// Защищенные эндпоинты
	mux.HandleFunc("POST /game/create", auth.Middleware(func(w http.ResponseWriter, r *http.Request) {
		userUUID := r.Context().Value(UserContextKey).(uuid.UUID)
		u, _ := user.GetUserInfo(userUUID)

		var req struct {
			Type string `json:"type"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			req.Type = domain.TypePVE
		}

		game, err := service.CreateGame(*u, req.Type)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DomainToWeb(*game))
	}))

	mux.HandleFunc("POST /game/{uuid}/join", auth.Middleware(func(w http.ResponseWriter, r *http.Request) {
		userUUID := r.Context().Value(UserContextKey).(uuid.UUID)
		u, _ := user.GetUserInfo(userUUID)

		gameID, err := uuid.Parse(r.PathValue("uuid"))
		if err != nil {
			http.Error(w, "некорректный UUID игры", http.StatusBadRequest)
			return
		}

		game, err := service.JoinGame(gameID, *u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DomainToWeb(*game))
	}))

	mux.HandleFunc("POST /game/{uuid}", auth.Middleware(func(w http.ResponseWriter, r *http.Request) {
		userUUID := r.Context().Value(UserContextKey).(uuid.UUID)
		
		gameID, err := uuid.Parse(r.PathValue("uuid"))
		if err != nil {
			http.Error(w, "некорректный UUID игры", http.StatusBadRequest)
			return
		}

		var webGame GameWeb
		if err := json.NewDecoder(r.Body).Decode(&webGame); err != nil {
			http.Error(w, "некорректный JSON", http.StatusBadRequest)
			return
		}

		resultGame, err := service.NextMoveService(gameID, userUUID, domain.Field(webGame.Matrix))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DomainToWeb(*resultGame))
	}))

	// Публичные эндпоинты (только чтение)
	mux.HandleFunc("GET /game/{uuid}", func(w http.ResponseWriter, r *http.Request) {
		gameID, err := uuid.Parse(r.PathValue("uuid"))
		if err != nil {
			http.Error(w, "некорректный UUID игры", http.StatusBadRequest)
			return
		}
		game, err := service.GetGame(gameID)
		if err != nil {
			http.Error(w, "игра не найдена", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DomainToWeb(*game))
	})

	mux.HandleFunc("GET /check", func(w http.ResponseWriter, r *http.Request) {
		sliceSessions, err := service.DomainCheckGameSesion()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sliceSessions)
	})

	historyHandler := auth.Middleware(func(w http.ResponseWriter, r *http.Request) {
		userUUID := r.Context().Value(UserContextKey).(uuid.UUID)
		completedGames, err := service.GetCompletedGames(userUUID)
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
	mux.HandleFunc("GET /games/history", historyHandler)
	mux.HandleFunc("GET /game/history", historyHandler)
	mux.HandleFunc("GET /history", historyHandler)

	leaderboardHandler := auth.Middleware(func(w http.ResponseWriter, r *http.Request) {
		limitStr := r.PathValue("n")
		if limitStr == "" {
			limitStr = r.URL.Query().Get("n")
		}
		if limitStr == "" {
			limitStr = r.URL.Query().Get("limit")
		}

		limit := 10
		if limitStr != "" {
			var parsedLimit int
			if _, err := fmt.Sscanf(limitStr, "%d", &parsedLimit); err == nil && parsedLimit > 0 {
				limit = parsedLimit
			}
		}

		entries, err := service.GetLeaderboard(limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entries)
	})
	mux.HandleFunc("GET /leaderboard", leaderboardHandler)
	mux.HandleFunc("GET /leaderboard/{n}", leaderboardHandler)

	mux.HandleFunc("GET /user/{uuid}", func(w http.ResponseWriter, r *http.Request) {
		userUUID, err := uuid.Parse(r.PathValue("uuid"))
		if err != nil {
			http.Error(w, "некорректный UUID пользователя", http.StatusBadRequest)
			return
		}
		u, err := user.GetUserInfo(userUUID)
		if err != nil {
			http.Error(w, "пользователь не найден", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DomainToWebUser(*u))
	})

	srv := &http.Server{
		Addr:    httpPort,
		Handler: mux,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("Запуск сервера на %s\n", srv.Addr)
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("Ошибка сервера: %s\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Остановка сервера...")
			return srv.Shutdown(ctx)
		},
	})
}
