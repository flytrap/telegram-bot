package services

import (
	"github.com/google/wire"
)

var ServiceSet = wire.NewSet(NewDataService, NewSearchService, NewDataTagService, NewCategoryService, NewIndexMangerService, NewUserService, NewAdService, NewGroupSettingService)
