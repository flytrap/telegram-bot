// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"github.com/flytrap/telegram-bot/internal/api"
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
	dataInfoRepository := repositories.NewDataInfoRepository(db)
	dataTagRepository := repositories.NewDataTagRepository(db)
	dataTagService := services.NewDataTagService(dataTagRepository)
	categoryRepository := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepository)
	dataService := services.NewDataService(dataInfoRepository, dataTagService, categoryService)
	coreClient, err := InitStore()
	if err != nil {
		return nil, err
	}
	indexMangerService := services.NewIndexMangerService(coreClient, dataService)
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	adRepository := repositories.NewAdRepository(db)
	adService := services.NewAdService(adRepository, categoryService)
	botManager := services.NewBotManager(dataService, indexMangerService, bot, userService, adService)
	tgBotServiceServer := api.NewTgBotApi(dataService)
	adServiceServer := api.NewAdApi(adService)
	tagServiceServer := api.NewDataTagApi(dataTagService)
	categoryServiceServer := api.NewCategoryApi(categoryService)
	userServiceServer := api.NewUserApi(userService)
	grpcServer := InitGrpcServer(tgBotServiceServer, adServiceServer, tagServiceServer, categoryServiceServer, userServiceServer)
	injector := &Injector{
		Bot:          bot,
		BotManager:   botManager,
		IndexManager: indexMangerService,
		GrpcServer:   grpcServer,
	}
	return injector, nil
}
