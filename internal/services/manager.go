package services

import (
	"strings"

	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type BotManager interface {
	Start()
}

func NewBotManager(gs GroupService, bot *tele.Bot) BotManager {
	return &BotManagerImp{Gs: gs, Bot: bot}
}

type BotManagerImp struct {
	Gs  GroupService
	Bot *tele.Bot
}

func (s *BotManagerImp) Start() {
	s.registRoute()
	logrus.Info("启动bot")
	s.Bot.Start()
}

func (s *BotManagerImp) registRoute() {
	s.Bot.Handle(tele.OnText, func(ctx tele.Context) error {
		tag := ctx.Message().Text
		data, err := s.Gs.SearchTag(tag, 10, 1)
		if err != nil {
			return err
		}
		items := []string{}
		for i, item := range data {
			items = append(items, item.ItemInfo(i+1))
		}
		ctx.Send(strings.Join(items, "\n"), tele.ModeMarkdown)
		return nil
	})
}
