package services

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type BotManager interface {
	Start()           // 启动机器人
	UpdateGroupInfo() // 更新群组数据
}

func NewBotManager(gs GroupService, bot *tele.Bot) BotManager {
	return &BotManagerImp{Gs: gs, Bot: bot, Waterline: rand.Intn(15) + 5}
}

type BotManagerImp struct {
	Gs        GroupService
	Bot       *tele.Bot
	Waterline int // 请求间隔时间
	Tick      int // 请求计数
}

func (s *BotManagerImp) Start() {
	s.registRoute()
	logrus.Info("启动bot")
	s.Bot.Start()
}

func (s *BotManagerImp) registRoute() {
	s.Bot.Handle(tele.OnText, func(ctx tele.Context) error {
		tag := ctx.Message().Text
		data, err := s.Gs.SearchTag(tag, 10, 1)
		if err != nil {
			return err
		}
		items := []string{}
		for i, item := range data {
			items = append(items, item.ItemInfo(i+1))
		}
		ctx.Send(strings.Join(items, "\n"), tele.ModeMarkdown)
		return nil
	})

	s.Bot.Handle("/hello", func(c tele.Context) error {
		return c.Send("Hello!")
	})
}

func (s *BotManagerImp) UpdateGroupInfo() {
	i := 1
	res, err := s.Gs.GetNeedUpdateCode(10, 100, i)
	if err != nil {
		logrus.Warn(err)
		return
	}
	for _, item := range res {
		n, err := s.GetChatMembers(fmt.Sprintf("@%s", item))
		if err != nil {
			logrus.Warn(err)
			if strings.Contains(err.Error(), "chat not found") {
				s.Gs.Delete(item)
			}
			continue
		}

		res, err := s.GetChatInfo(fmt.Sprintf("@%s", item))
		if err != nil {
			logrus.Warn(err)
			continue
		}
		desc := ""
		if _, ok := res["description"]; ok {
			desc = res["description"].(string)
		}
		s.Gs.Update(item, int(res["id"].(float64)), res["title"].(string), desc, n)
	}
	if len(res) == 100 {
		s.UpdateGroupInfo()
	}
}

// 获取群人数
func (s *BotManagerImp) GetChatMembers(code string) (int, error) {
	s.CheckSleep()
	params := map[string]string{"chat_id": code}
	res, err := s.Bot.Raw("getChatMemberCount", params)
	s.Tick += 1
	if err != nil {
		return 0, err
	}
	result := map[string]interface{}{}
	if err := json.Unmarshal(res, &result); err != nil {
		return 0, err
	}
	if result["ok"].(bool) {
		return int(result["result"].(float64)), nil
	}
	return 0, nil
}

// 获取群信息
func (s *BotManagerImp) GetChatInfo(code string) (map[string]interface{}, error) {
	s.CheckSleep()
	params := map[string]string{"chat_id": code}
	res, err := s.Bot.Raw("getChat", params)
	if err != nil {
		return nil, err
	}
	result := map[string]interface{}{}
	if err := json.Unmarshal(res, &result); err != nil {
		return nil, err
	}
	return result["result"].(map[string]interface{}), nil
}

// 检查是否需要暂停一下
func (s *BotManagerImp) CheckSleep() {
	s.Tick += 1
	if s.Tick >= s.Waterline {
		s.Tick = 0
		s.Waterline = rand.Intn(15) + 5
		time.Sleep(time.Second * time.Duration(s.Waterline))
	}
}
