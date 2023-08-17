// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"github.com/flytrap/telegram-bot/internal/repositories"
	"github.com/flytrap/telegram-bot/internal/services"
)

// Injectors from wire.go:

func BuildInjector() (*Injector, error) {
	bot, err := InitBot()
	if err != nil {
		return nil, err
	}
	db, err := InitGormDB()
	if err != nil {
		return nil, err
	}
	groupRepository := repositories.NewGroupRepository(db)
	tagRepository := repositories.NewTagRepository(db)
	tagService := services.NewTagService(tagRepository)
	groupService := services.NewGroupService(groupRepository, tagService)
	botManager := services.NewBotManager(groupService, bot)
	tgBotServiceServer := services.NewTgBotService(groupService)
	grpcServer := InitGrpcServer(tgBotServiceServer)
	injector := &Injector{
		Bot:        bot,
		BotManager: botManager,
		GrpcServer: grpcServer,
	}
	return injector, nil
}
