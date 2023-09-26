package handlers

import (
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// 消息删除
func (s *HandlerManagerImp) DelMessageHandler(ctx tele.Context) error {
	ctx.Delete()
	if ctx.Chat().Private {
		return ctx.Send("请将本机器人加入您的群中再使用此命令")
	}
	if !isAdmin(ctx) {
		return ctx.Send("只有本群管理才可以使用此命令")
	}
	msg := ctx.Message().ReplyTo

	ctx.Delete()
	logrus.Info("delete: ", msg)
	return ctx.Bot().Delete(msg)
}
