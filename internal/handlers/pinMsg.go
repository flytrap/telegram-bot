package handlers

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	tele "gopkg.in/telebot.v3"
)

// 消息置顶
func (s *HandlerManagerImp) PinMessageHandler(ctx tele.Context) error {
	localize := i18n.NewLocalizer(s.bundle, "zh-CN")
	if ctx.Chat().Private {
		adminJoinTip := localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "admin.joinTip"})
		return s.sendAutoDeleteMessage(ctx, AfterDelTime(), adminJoinTip)
	}
	if !isAdmin(ctx) {
		noPermissionTio := localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "admin.noPermissionTio"})
		return s.sendAutoDeleteMessage(ctx, AfterDelTime(), noPermissionTio)
	}
	msg := ctx.Message().ReplyTo

	ctx.Delete()
	return ctx.Bot().Pin(msg)
}
