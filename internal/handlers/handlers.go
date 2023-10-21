package handlers

import (
	"context"

	"github.com/flytrap/telegram-bot/internal/config"
	"github.com/flytrap/telegram-bot/internal/middleware"
	"github.com/flytrap/telegram-bot/internal/services"
	"github.com/flytrap/telegram-bot/pkg/redis"
	"github.com/google/wire"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

var HandlerSet = wire.NewSet(NewHandlerManager)

type HandlerManager interface {
	Start(ctx context.Context, openIndex bool) // 启动机器人
	registerRoute(bool) error                  // 注册路由
	CheckDeleteMessage(ctx context.Context)    // 检查需要删除的消息
	UpdateGroupInfo(int64) error               // 更新群组数据
}

func NewHandlerManager(bot *tele.Bot, store *redis.Store, ss services.SearchService, gs services.GroupSettingService, dataService services.DataService, cs services.CategoryService, m middleware.MiddleWareManager, bundle *i18n.Bundle) HandlerManager {
	return &HandlerManagerImp{Bot: bot, store: store, ss: ss, gs: gs, dataService: dataService, cs: cs, m: m, bundle: bundle}
}

type HandlerManagerImp struct {
	Bot         *tele.Bot
	store       *redis.Store
	dataService services.DataService
	ss          services.SearchService
	cs          services.CategoryService
	gs          services.GroupSettingService
	m           middleware.MiddleWareManager
	bundle      *i18n.Bundle
}

func (s *HandlerManagerImp) Start(ctx context.Context, openIndex bool) {
	s.Bot.Use(s.m.LinkFilter())
	s.Bot.Use(s.m.Logger())
	logrus.Info("启动bot")
	s.registerRoute(openIndex)
	s.Bot.Start()
}

func (s *HandlerManagerImp) registerRoute(openIndex bool) error {
	s.Bot.Handle("/start", s.StartHandler)
	if config.C.Bot.OpenManager {
		s.Bot.Handle(tele.OnAddedToGroup, s.InitHandler)
		s.Bot.Handle(tele.OnUserJoined, s.JoinInHandler)
		s.Bot.Handle(tele.OnUserLeft, s.AutoDeleteInHandler)
		s.Bot.Handle("/admin", s.AdminHandler)
		s.Bot.Handle("/help", s.HelpHandler)
	}
	if openIndex {
		s.Bot.Handle(config.C.Index.Command.Category, s.CategoryHandler)
		s.Bot.Handle(config.C.Index.Command.CategoryHelp, s.CategoryHelpHandler)
		s.Bot.Handle(config.C.Index.Command.CategoryTag, s.CategoryTagHandler)
		s.Bot.Handle(config.C.Index.Command.CategorySearch, s.CategoryQHandler)
		s.Bot.Handle(tele.OnText, s.IndexHandler)
	}
	return nil
}
