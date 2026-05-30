package web

import (
	"anglefar/tic-tac-toe/internal/domain"
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type contextKey string

const UserContextKey contextKey = "user_uuid"

type UserAuthenticator struct {
	userService *domain.UserService
}

func NewUserAuthenticator(userService *domain.UserService) *UserAuthenticator {
	return &UserAuthenticator{userService: userService}
}

// Middleware для защиты эндпоинтов
func (a *UserAuthenticator) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем токен из куки или заголовка
		var userUUID uuid.UUID
		var err error

		cookie, err := r.Cookie("session_token")
		if err == nil {
			userUUID, err = uuid.Parse(cookie.Value)
		} else {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				if strings.HasPrefix(authHeader, "Bearer ") {
					userUUID, err = a.validateBearerAuth(authHeader)
				} else {
					userUUID, err = a.validateBasicAuth(authHeader)
				}
			} else {
				err = http.ErrNoCookie
			}
		}

		if err != nil {
			http.Error(w, "неавторизован", http.StatusUnauthorized)
			return
		}

		// Проверяем существование пользователя
		_, err = a.userService.GetUserInfo(userUUID)
		if err != nil {
			http.Error(w, "неавторизован", http.StatusUnauthorized)
			return
		}

		// Добавляем UUID в контекст
		ctx := context.WithValue(r.Context(), UserContextKey, userUUID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (a *UserAuthenticator) validateBasicAuth(authHeader string) (uuid.UUID, error) {
	if !strings.HasPrefix(authHeader, "Basic ") {
		return uuid.Nil, http.ErrNoCookie
	}

	payload, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
	if err != nil {
		return uuid.Nil, err
	}

	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		return uuid.Nil, http.ErrNoCookie
	}

	login, pass := pair[0], pair[1]
	if login == "" || pass == "" {
		return uuid.Nil, http.ErrNoCookie
	}

	userUUID, err := a.userService.LoginUser(domain.UserDomain{Login: login, Pass: pass})
	if err != nil {
		return uuid.Nil, err
	}

	return userUUID, nil
}

func (a *UserAuthenticator) validateBearerAuth(authHeader string) (uuid.UUID, error) {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return uuid.Nil, http.ErrNoCookie
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	userUUID, err := uuid.Parse(token)
	if err != nil {
		return uuid.Nil, err
	}

	return userUUID, nil
}
