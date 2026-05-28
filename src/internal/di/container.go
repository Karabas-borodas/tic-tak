package di

import (
	"anglefar/tic-tac-toe/internal/datasource"
	"anglefar/tic-tac-toe/internal/domain"
	"anglefar/tic-tac-toe/internal/web"

	"go.uber.org/fx"
)

// Module собирает все зависимости в один fx.Option для использования в main.
var Module = fx.Options(
	// 1. Предоставляем (Provide) все наши конструкторы
	fx.Provide(
		datasource.NewUserStorage,
		fx.Annotate(
			func(s *datasource.UserStorage) domain.IUser { return s },
			fx.As(new(domain.IUser)),
		),
		domain.NewUserService,

		datasource.NewGameStorage, // Хранилище
		// Связываем конкретное хранилище с интерфейсом
		fx.Annotate(
			func(s *datasource.GameStorage) domain.IStorage { return s },
			fx.As(new(domain.IStorage)),
		),

		domain.NewGameService,      // Сервис
		web.NewUserAuthenticator, // Аутентификатор
	),

	// 2. Вызываем (Invoke) функцию, которая запустит HTTP сервер
	fx.Invoke(web.StartGame),
)
