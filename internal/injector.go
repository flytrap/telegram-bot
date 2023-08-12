package app

import (
	"github.com/flytrap/telegram-bot/internal/services"
	"github.com/google/wire"
	tele "gopkg.in/telebot.v3"
)

var InjectorSet = wire.NewSet(wire.Struct(new(Injector), "*"))

type Injector struct {
	Bot        *tele.Bot
	BotManager services.BotManager
}
