package handlers

import (
	"context"

	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// 初始化信息
func (s *HandlerManagerImp) InitHandler(ctx tele.Context) error {
	if len(ctx.Chat().Username) > 0 {
		err := s.gs.Init(context.Background(), ctx.Chat().Username)
		if err != nil {
			logrus.Warning(err)
			return err
		}
	}
	return s.HelpHandler(ctx)
}
