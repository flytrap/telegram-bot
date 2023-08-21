package services

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type BotManager interface {
	Start()              // 启动机器人
	UpdateGroupInfo(int) // 更新群组数据
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

	s.Bot.Handle(tele.OnText, s.SearchGroup)

	s.Bot.Handle("/hello", func(c tele.Context) error {
		return c.Send("Hello!")
	})
}

func (s *BotManagerImp) SearchGroup(ctx tele.Context) error {
	selector := &tele.ReplyMarkup{}
	tag := ""
	page := 1
	cb := ctx.Callback()
	if cb != nil {
		info := cb.Data
		tag = strings.Split(info, "|")[0]
		n, err := strconv.Atoi(strings.Split(info, "|")[1])
		if err != nil {
			logrus.Warn(err)
		} else {
			page = n
		}
	} else {
		tag = ctx.Message().Text
	}
	data, err := s.Gs.SearchTag(tag, 15, page)
	if err != nil {
		return err
	}
	items := []string{}
	for i, item := range data {
		items = append(items, item.ItemInfo(i+1))
	}
	var btnPrev tele.Btn
	var btnNext tele.Btn
	if page > 1 {
		btnPrev = selector.Data("⬅ prev", "prev", tag, fmt.Sprint(page-1))
		s.Bot.Handle(&btnPrev, s.SearchGroup)
	}
	if len(data) == 15 {
		btnNext = selector.Data("next ➡", "next", tag, fmt.Sprint(page+1))
	}
	selector.Inline(selector.Row(btnPrev, btnNext))
	s.Bot.Handle(&btnNext, s.SearchGroup)
	ctx.EditOrSend(strings.Join(items, "\n"), tele.ModeMarkdown, selector)
	return nil
}

func (s *BotManagerImp) UpdateGroupInfo(num int) {
	i := 1
	res, err := s.Gs.GetNeedUpdateCode(10, num, i)
	if err != nil {
		logrus.Warn(err)
		return
	}
	for _, item := range res {
		s.updateInfo(item)
	}
}

func (s *BotManagerImp) updateInfo(code string) {
	n, err := s.GetChatMembers(fmt.Sprintf("@%s", code))
	if err != nil {
		logrus.Warn(err)
		if strings.Contains(err.Error(), "chat not found") {
			s.Gs.Delete(code)
		} else if strings.Contains(err.Error(), "retry after") {
			time.Sleep(time.Second * 999)
		}
		return
	}

	res, err := s.GetChatInfo(fmt.Sprintf("@%s", code))
	if err != nil {
		logrus.Warn(err)
		return
	}
	desc := ""
	if _, ok := res["description"]; ok {
		desc = res["description"].(string)
	}
	s.Gs.Update(code, int64(res["id"].(float64)), res["title"].(string), desc, n)
}

// 获取群人数
func (s *BotManagerImp) GetChatMembers(code string) (uint32, error) {
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
		return uint32(result["result"].(float64)), nil
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
