package handlers

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/flytrap/telegram-bot/internal/config"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

var (
	groupChannel      *tele.Chat
	recommendTagIndex int
)

func (s *HandlerManagerImp) initChannel() error {
	if len(config.C.Index.Recommend.Channel) == 0 {
		return nil
	}
	c, err := s.Bot.ChatByUsername(fmt.Sprintf("@%s", config.C.Index.Recommend.Channel))
	if err != nil {
		logrus.Warning(err)
		return err
	}
	groupChannel = c
	recommendTagIndex = 0
	return nil
}

// 给channel发送消息
func (s *HandlerManagerImp) SendDataToChannel() {
	if groupChannel == nil {
		s.initChannel()
	}
	if groupChannel == nil {
		logrus.Warning("not found channel")
		return
	}
	tag := ""
	for {
		if len(config.C.Index.Recommend.Tags) > 0 {
			tag = config.C.Index.Recommend.Tags[recommendTagIndex%len(config.C.Index.Recommend.Tags)]
			recommendTagIndex += 1
		}
		items, _, err := s.ss.QueryItems(context.Background(), config.C.Index.Recommend.Category, tag, config.C.Index.Recommend.Q, 1, 1000)
		if err != nil || len(items) == 0 {
			logrus.Warning(err)
			time.Sleep(time.Hour * 3) // 3个小时发送一次
			continue
		}
		item := items[rand.Intn(len(items)-1)]
		selector := &tele.ReplyMarkup{}
		result, ps, err := s.getDetailInfo(item["code"].(string), selector)
		if err != nil {
			time.Sleep(time.Hour * 3) // 3个小时发送一次
			continue
		}
		if len(*ps) > 0 {
			s.Bot.Send(groupChannel, *ps)
		}
		s.Bot.Send(groupChannel, result, selector)
		time.Sleep(time.Hour * 3) // 3个小时发送一次
	}
}
