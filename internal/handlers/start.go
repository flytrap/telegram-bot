package handlers

import (
	"fmt"

	"github.com/flytrap/telegram-bot/internal/config"
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

	text := "本机器人要加入群中才能工作 \n"
	text += "Telegram语言包设置： \n【 [简体中文](tg://setlanguage?lang=zhcncc) | [繁體中文](tg://setlanguage?lang=zh-hant-beta) | [English](tg://setlanguage?lang=en) 】 \n"
	selector.URL("添加本机器人到群中", fmt.Sprintf("https://t.me/%s", config.C.Bot.Manager.Username))
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
