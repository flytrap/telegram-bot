//go:build wireinject
// +build wireinject

package app

import (
	"github.com/flytrap/telegram-bot/internal/api"
	"github.com/flytrap/telegram-bot/internal/handlers"
	"github.com/flytrap/telegram-bot/internal/middleware"
	"github.com/flytrap/telegram-bot/internal/repositories"
	"github.com/flytrap/telegram-bot/internal/services"
	"github.com/google/wire"
)

func BuildInjector() (*Injector, error) {
	wire.Build(
		InitBundle,
		InitBot,
		InitIndex,
		InitStore,
		InitGormDB,
		middleware.MiddleWareSet,
		repositories.RepositorySet,
		handlers.HandlerSet,
		InitGrpcServer,
		// InitGateway,
		services.ServiceSet,
		api.APISet,
		// router.RouterSet,
		InjectorSet,
	)
	return new(Injector), nil
}
