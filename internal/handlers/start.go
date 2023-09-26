package handlers

import (
	tele "gopkg.in/telebot.v3"
)

func (s *HandlerManagerImp) StartHandler(ctx tele.Context) error {
	if ctx.Chat().Private {
		return nil
	}
	selector := &tele.ReplyMarkup{}

	text := "本机器人要加入群中才能工作 \n"
	selector.URL("添加本机器人到群中", "https://t.me/test")
	return s.sendAutoDeleteMessage(ctx, AfterDelTime(), text, selector)
}
