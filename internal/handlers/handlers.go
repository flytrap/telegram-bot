package handlers

import (
	"context"

	"github.com/flytrap/telegram-bot/internal/services"
	"github.com/flytrap/telegram-bot/pkg/redis"
	"github.com/google/wire"
	tele "gopkg.in/telebot.v3"
)

var HandlerSet = wire.NewSet(NewHandlerManager)

type HandlerManager interface {
	RegisterRoute() error
	CheckDeleteMessage(ctx context.Context)
}

func NewHandlerManager(bot *tele.Bot, store *redis.Store, bm services.BotManager, gs services.GroupSettingSerivce) HandlerManager {
	return &HandlerManagerImp{Bot: bot, store: store, bm: bm, gs: gs}
}

type HandlerManagerImp struct {
	Bot   *tele.Bot
	store *redis.Store
	bm    services.BotManager
	gs    services.GroupSettingSerivce
}

func (s *HandlerManagerImp) RegisterRoute() error {
	s.Bot.Handle(tele.OnAddedToGroup, s.InitHandler)
	s.Bot.Handle(tele.OnUserJoined, s.JoinInHandler)
	s.Bot.Handle(tele.OnUserLeft, s.AutoDeleteInHandler)
	s.Bot.Handle("/admin", s.AdminHandler)
	s.Bot.Handle("/help", s.HelpHandler)
	s.Bot.Handle("/del", s.DelMessageHandler)
	s.Bot.Handle("/pin", s.PinMessageHandler)
	s.Bot.Handle(tele.OnText, s.IndexHandler)
	return nil
}
