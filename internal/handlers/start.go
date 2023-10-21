package handlers

import (
	"fmt"

	"github.com/flytrap/telegram-bot/internal/config"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	tele "gopkg.in/telebot.v3"
)

var (
	menu = &tele.ReplyMarkup{ResizeKeyboard: true}
)

func (s *HandlerManagerImp) StartHandler(ctx tele.Context) error {
	if ctx.Chat().Private {
		return nil
	}
	if menu == nil {
		initMenu()
	}
	selector := &tele.ReplyMarkup{}
	localize := i18n.NewLocalizer(s.bundle, "zh-CN")
	text := localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "startTip"})

	selector.URL(config.C.Bot.Manager.Username, fmt.Sprintf("https://t.me/%s", config.C.Bot.Manager.Username))
	return s.sendAutoDeleteMessage(ctx, AfterDelTime(), text, selector, menu)
}

func initMenu() {
	items := []tele.Row{}
	for _, item := range config.C.Bot.Menus {
		subItem := []tele.Btn{}
		for _, su := range item {
			subItem = append(subItem, menu.Text(su))
		}
		items = append(items, menu.Row(subItem...))
	}
	menu.Reply(items...)
}
