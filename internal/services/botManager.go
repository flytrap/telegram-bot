package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/flytrap/telegram-bot/internal/config"
	"github.com/flytrap/telegram-bot/internal/middleware"
	"github.com/flytrap/telegram-bot/pkg/human"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type BotManager interface {
	Start()                // 启动机器人
	UpdateGroupInfo(int64) // 更新群组数据
}

func NewBotManager(dataService DataService, im IndexMangerService, bot *tele.Bot, userService UserService) BotManager {
	return &BotManagerImp{dataService: dataService, IndexManager: im, userService: userService, Bot: bot, Waterline: rand.Intn(15) + 5}
}

type BotManagerImp struct {
	dataService  DataService
	IndexManager IndexMangerService
	userService  UserService
	Bot          *tele.Bot
	Waterline    int // 请求间隔时间
	Tick         int // 请求计数
}

func (s *BotManagerImp) Start() {
	s.Bot.Use(middleware.Logger(func(user map[string]interface{}) error {
		return s.userService.CreateOrUpdate(user)
	}))
	s.registerRoute()
	logrus.Info("启动bot")
	s.Bot.Start()
}

func (s *BotManagerImp) registerRoute() {

	s.Bot.Handle(tele.OnText, s.SearchGroups)

	s.Bot.Handle("/lang", func(c tele.Context) error {
		return c.Send("Lang!")
	})
}

func (s *BotManagerImp) SearchGroups(ctx tele.Context) error {
	selector := &tele.ReplyMarkup{}
	tag := ""
	page := int64(1)
	cb := ctx.Callback()
	if cb != nil {
		info := cb.Data
		tag = strings.Split(info, "|")[0]
		n, err := strconv.Atoi(strings.Split(info, "|")[1])
		if err != nil {
			logrus.Warn(err)
		} else {
			page = int64(n)
		}
	} else {
		tag = ctx.Message().Text
	}
	lang := ctx.Sender().LanguageCode
	logrus.Info(lang)

	items, hasNext, err := s.QueryItems(context.Background(), tag, page, 15)
	if err != nil {
		logrus.Warning(err)
		return err
	}
	var btnPrev tele.Btn
	var btnNext tele.Btn
	if page > 1 {
		btnPrev = selector.Data("⬅ prev", "prev", tag, fmt.Sprint(page-1))
		s.Bot.Handle(&btnPrev, s.SearchGroups)
	}
	if hasNext {
		btnNext = selector.Data("next ➡", "next", tag, fmt.Sprint(page+1))
	}
	selector.Inline(selector.Row(btnPrev, btnNext))
	s.Bot.Handle(&btnNext, s.SearchGroups)
	ctx.EditOrSend(strings.Join(items, "\n"), tele.ModeMarkdown, selector)
	return nil
}

func (s *BotManagerImp) QueryItems(ctx context.Context, text string, page int64, size int64) (items []string, hasNext bool, err error) {
	if config.C.Bot.UseCache {
		items, err = s.QueryCacheItems(ctx, "chinese", text, "", page, 15)
	} else {
		items, err = s.QueryDbItems(text, page, 15)
	}
	if err != nil {
		return nil, false, err
	}
	hasNext = int64(len(items)) == size
	return items, hasNext, nil
}

func (s *BotManagerImp) QueryDbItems(text string, page int64, size int64) ([]string, error) {
	data, err := s.dataService.SearchTag(text, page, 15)
	if err != nil {
		return nil, err
	}
	items := []string{}
	for i, item := range data {
		items = append(items, item.ItemInfo(i+1))
	}
	return items, nil
}

func (s *BotManagerImp) QueryCacheItems(ctx context.Context, name string, text string, category string, page int64, size int64) ([]string, error) {
	index := s.IndexManager.IndexName(name)
	data, err := s.IndexManager.Query(ctx, index, text, category, page, size)
	if err != nil {
		return nil, err
	}
	items := []string{}
	for i, item := range data {
		items = append(items, human.TgGroupItemInfo(int(page-1)*int(size)+i+1, item["code"].(string), int(item["type"].(float64)), item["name"].(string), int64(item["num"].(float64))))
	}
	return items, nil
}

func (s *BotManagerImp) UpdateGroupInfo(num int64) {
	i := int64(1)
	res, err := s.dataService.GetNeedUpdateCode(10, i, num)
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
			s.dataService.Delete(code)
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
	s.dataService.Update(code, int64(res["id"].(float64)), res["title"].(string), desc, n)
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
