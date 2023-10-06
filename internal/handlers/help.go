package handlers

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	tele "gopkg.in/telebot.v3"
)

func (s *HandlerManagerImp) HelpHandler(ctx tele.Context) error {
	localize := i18n.NewLocalizer(s.bundle, "zh-CN")
	text := localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "botHelp"})

	return s.sendAutoDeleteMessage(ctx, AfterDelTime(), text)
}
