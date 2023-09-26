package services

import (
	"github.com/google/wire"
)

var ServiceSet = wire.NewSet(NewDataService, NewBotManager, NewDataTagService, NewCategoryService, NewIndexMangerService, NewUserService, NewAdService, NewGroupSettingService)
