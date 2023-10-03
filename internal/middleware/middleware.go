package middleware

import (
	"github.com/flytrap/telegram-bot/internal/services"
	"github.com/flytrap/telegram-bot/pkg/redis"
	"github.com/google/wire"
	tele "gopkg.in/telebot.v3"
)

var MiddleWareSet = wire.NewSet(NewMiddleWareManager)

type MiddleWareManager interface {
	Logger() tele.MiddlewareFunc
	LinkFilter() tele.MiddlewareFunc
}

func NewMiddleWareManager(us services.UserService, store *redis.Store) MiddleWareManager {
	return &middleWareManagerImp{userService: us, store: store}
}

type middleWareManagerImp struct {
	userService services.UserService
	store       *redis.Store
}
