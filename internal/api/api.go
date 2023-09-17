package api

import (
	"github.com/google/wire"
)

var APISet = wire.NewSet(NewTgBotApi, NewAdApi, NewCategoryApi, NewDataTagApi, NewUserApi)
