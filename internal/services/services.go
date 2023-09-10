package services

import (
	"github.com/google/wire"
)

var ServiceSet = wire.NewSet(NewDataService, NewBotManager, NewDataTagService, NewTgBotService, NewCategoryService, NewIndexMangerService, NewUserService, NewAdService)
