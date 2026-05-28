# Tic-Tac-Toe API

Игра крестики-нолики с поддержкой режимов PVP (игрок против игрока) и PVE (игрок против ИИ Minimax), системой состояний и авторизацией.

## Запуск

1. Убедитесь, что PostgreSQL запущен.
2. Запустите сервер:
```bash
go run cmd/main/main.go
```
Сервер запустится на `http://localhost:3002`.

## API (curl-команды)

### 1. Регистрация и Авторизация

Для авторизации используется стандартная схема **Basic Auth**. Логин и пароль передаются в заголовке `Authorization` в виде `base64(login:password)`.

```bash
# Регистрация (JSON)
curl -X POST http://localhost:3002/register \
  -H "Content-Type: application/json" \
  -d '{"login":"player1", "password":"password123"}'

# Авторизация (Basic Auth)
# Возвращает JSON с UUID пользователя и устанавливает сессионную куку
curl -X POST http://localhost:3002/login \
  -H "Authorization: Basic $(echo -n 'player1:password123' | base64)" \
  -c cookies.txt
```

### 2. Создание игры (Требуется авторизация)

Вы можете создать игру против компьютера (PVE) или против другого игрока (PVP).

```bash
# Создать PVE игру
curl -X POST http://localhost:3002/game/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Basic $(echo -n 'player1:password123' | base64)" \
  -d '{"type":"PVE"}'

# Или используя куки:
curl -X POST http://localhost:3002/game/create \
  -H "Content-Type: application/json" \
  -d '{"type":"PVP"}' \
  -b cookies.txt
```

### 3. Поиск и присоединение к игре (PVP)

```bash
# Список игр, ожидающих игрока (Публичный)
curl -X GET http://localhost:3002/check

# Присоединиться к игре по UUID (Требуется авторизация)
curl -X POST http://localhost:3002/game/{UUID_ИГРЫ}/join \
  -b cookies.txt
```

### 4. Игровой процесс

Для совершения хода используйте UUID игры. Ходить можно только в свой ход.

```bash
# Получить текущее состояние игры (Публичный)
curl -X GET http://localhost:3002/game/{UUID_ИГРЫ}

# Сделать ход (Требуется авторизация)
# В matrix подставьте 1 (X) или 2 (O) в нужную ячейку.
curl -X POST http://localhost:3002/game/{UUID_ИГРЫ} \
  -H "Content-Type: application/json" \
  -d '{"matrix":[[1,0,0],[0,0,0],[0,0,0]]}' \
  -b cookies.txt
```

### 5. Информация о пользователях

```bash
# Получить профиль по UUID
curl -X GET http://localhost:3002/user/{UUID_ПОЛЬЗОВАТЕЛЯ}
```

## Состояния игры (`status`)

- `waiting`: Ожидание второго игрока.
- `turn_p1`: Ход первого игрока (X).
- `turn_p2`: Ход второго игрока (O).
- `won_p1`: Победил первый игрок.
- `won_p2`: Победил второй игрок.
- `draw`: Ничья.

## Значения в матрице

- `0`: Пустая ячейка.
- `1`: Крестик (X) — Игрок 1.
- `2`: Нолик (O) — Игрок 2 или Компьютер.