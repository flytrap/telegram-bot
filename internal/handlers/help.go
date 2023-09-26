package handlers

import (
	tele "gopkg.in/telebot.v3"
)

func (s *HandlerManagerImp) HelpHandler(ctx tele.Context) error {
	// ctx.Message().Chat.ID
	text := "欢迎使用群管机器人！！\n\n"
	text += "以下是本BOT的独家命令： \n\n /help - 输出帮助" + "\n" + "/admin - 机器人设置" + "\n" + "/ban - 用此命令回复消息踢发消息的人" + "\n" + "\n"

	return s.sendAutoDeleteMessage(ctx, AfterDelTime(), text)
}
