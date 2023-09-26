package handlers

import (
	"time"

	tele "gopkg.in/telebot.v3"
)

// 消息置顶
func (s *HandlerManagerImp) PinMessageHandler(ctx tele.Context) error {
	if ctx.Chat().Private {
		return s.sendAutoDeleteMessage(ctx, time.Minute, "请将本机器人加入您的群中再使用此命令")
	}
	if !isAdmin(ctx) {
		return s.sendAutoDeleteMessage(ctx, time.Minute, "只有本群管理才可以使用此命令")
	}
	msg := ctx.Message().ReplyTo

	ctx.Delete()
	return ctx.Bot().Pin(msg)
}
