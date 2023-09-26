package handlers

import tele "gopkg.in/telebot.v3"

// 自动删除
func (s *HandlerManagerImp) AutoDeleteInHandler(ctx tele.Context) error {
	return ctx.Delete()
}
