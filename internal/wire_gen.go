// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"github.com/flytrap/telegram-bot/internal/api"
	"github.com/flytrap/telegram-bot/internal/handlers"
	"github.com/flytrap/telegram-bot/internal/middleware"
	"github.com/flytrap/telegram-bot/internal/repositories"
	"github.com/flytrap/telegram-bot/internal/services"
)

// Injectors from wire.go:

func BuildInjector() (*Injector, error) {
	bot, err := InitBot()
	if err != nil {
		return nil, err
	}
	coreClient, err := InitIndex()
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
	indexMangerService := services.NewIndexMangerService(coreClient, dataService, categoryService)
	tgBotServiceServer := api.NewTgBotApi(dataService, categoryService)
	adRepository := repositories.NewAdRepository(db)
	adService := services.NewAdService(adRepository, categoryService)
	adServiceServer := api.NewAdApi(adService)
	tagServiceServer := api.NewDataTagApi(dataTagService)
	categoryServiceServer := api.NewCategoryApi(categoryService)
	userRepository := repositories.NewUserRepository(db)
	store, err := InitStore()
	if err != nil {
		return nil, err
	}
	userService := services.NewUserService(userRepository, store)
	userServiceServer := api.NewUserApi(userService)
	grpcServer := InitGrpcServer(tgBotServiceServer, adServiceServer, tagServiceServer, categoryServiceServer, userServiceServer)
	searchService := services.NewSearchService(dataService, indexMangerService, adService)
	groupSettingRepository := repositories.NewGroupSettingRepository(db)
	groupSettingService := services.NewGroupSettingService(groupSettingRepository, store)
	middleWareManager := middleware.NewMiddleWareManager(userService, store)
	bundle, err := InitBundle()
	if err != nil {
		return nil, err
	}
	handlerManager := handlers.NewHandlerManager(bot, store, searchService, groupSettingService, dataService, middleWareManager, bundle)
	injector := &Injector{
		Bot:            bot,
		IndexManager:   indexMangerService,
		GrpcServer:     grpcServer,
		HandlerManager: handlerManager,
	}
	return injector, nil
}
