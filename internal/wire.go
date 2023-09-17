//go:build wireinject
// +build wireinject

package app

import (
	"github.com/flytrap/telegram-bot/internal/api"
	"github.com/flytrap/telegram-bot/internal/repositories"
	"github.com/flytrap/telegram-bot/internal/services"
	"github.com/google/wire"
)

func BuildInjector() (*Injector, error) {
	wire.Build(
		InitBot,
		InitStore,
		InitGormDB,
		repositories.RepositorySet,
		InitGrpcServer,
		// InitGateway,
		services.ServiceSet,
		api.APISet,
		// router.RouterSet,
		InjectorSet,
	)
	return new(Injector), nil
}
