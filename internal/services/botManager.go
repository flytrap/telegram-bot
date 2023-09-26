package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/flytrap/telegram-bot/internal/config"
	"github.com/flytrap/telegram-bot/internal/middleware"
	"github.com/flytrap/telegram-bot/internal/serializers"
	"github.com/flytrap/telegram-bot/pkg/human"
	"github.com/flytrap/telegram-bot/pkg/redis"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type BotManager interface {
	Start(ctx context.Context)                                                // 启动机器人
	UpdateGroupInfo(int64)                                                    // 更新群组数据
	QueryItems(context.Context, string, int64, int64) ([]string, bool, error) // 搜素信息
}

func NewBotManager(dataService DataService, im IndexMangerService, bot *tele.Bot, userService UserService, adService AdService, store *redis.Store) BotManager {
	return &BotManagerImp{dataService: dataService, IndexManager: im, userService: userService, Bot: bot, adService: adService, store: store}
}

type BotManagerImp struct {
	dataService  DataService
	IndexManager IndexMangerService
	userService  UserService
	adService    AdService
	Bot          *tele.Bot
	store        *redis.Store
}

func (s *BotManagerImp) Start(ctx context.Context) {
	s.Bot.Use(middleware.Logger(func(user map[string]interface{}) error {
		if s.userService.Check(ctx, user["UserId"].(int64)) {
			return nil
		}
		return s.userService.GetOrCreate(ctx, user)
	}, func(userId int64, username string, q string) error {
		return s.store.Xadd(ctx, "log:query", map[string]interface{}{"user_id": userId, "content": q, "type": "tg-query", "username": username})
	}))
	s.registerRoute()
	logrus.Info("启动bot")
	s.Bot.Start()
}

func (s *BotManagerImp) registerRoute() {

	// s.Bot.Handle(tele.OnText,  )

	s.Bot.Handle("/lang", func(c tele.Context) error {
		return c.Send("Lang!")
	})
}

func (s *BotManagerImp) QueryItems(ctx context.Context, text string, page int64, size int64) (items []string, hasNext bool, err error) {
	if page >= config.C.Bot.MaxPage {
		page = config.C.Bot.MaxPage // 阻止过多翻页
	}
	var n int64
	if config.C.Bot.UseCache {
		n, items, err = s.QueryCacheItems(ctx, "chinese", text, "", page, config.C.Bot.PageSize)
	} else {
		n, items, err = s.QueryDbItems(text, page, config.C.Bot.PageSize)
	}
	if err != nil {
		return nil, false, err
	}
	hasNext = (page-1)*size+int64(len(items)) < n
	ad := s.loadAd(text)
	if len(ad) > 0 {
		items = append([]string{ad, ""}, items...) // 增加广告
	}
	return items, hasNext, nil
}

func (s *BotManagerImp) loadAd(keyword string) string {
	item, err := s.adService.KeywordAd(keyword)
	if err != nil {
		return ""
	}
	showAd := ""
	if len(item.AdTag) > 0 {
		showAd = fmt.Sprintf("[%s] ", item.AdTag)
	}
	return fmt.Sprintf("%s[%s](%s)", showAd, item.Name, human.TgGroupUrl(item.Code))
}

func (s *BotManagerImp) QueryDbItems(text string, page int64, size int64) (int64, []string, error) {
	data := []*serializers.DataSerializer{}
	n, err := s.dataService.SearchTag(text, page, 15, &data)
	if err != nil {
		return 0, nil, err
	}
	items := []string{}
	for i, item := range data {
		items = append(items, item.ItemInfo(i+1))
	}
	return n, items, nil
}

func (s *BotManagerImp) QueryCacheItems(ctx context.Context, name string, text string, category string, page int64, size int64) (int64, []string, error) {
	index := s.IndexManager.IndexName(name)
	n, data, err := s.IndexManager.Query(ctx, index, text, category, page, size)
	if err != nil {
		return 0, nil, err
	}
	items := []string{}
	for i, item := range data {
		items = append(items, human.TgGroupItemInfo(int(page-1)*int(size)+i+1, item["code"].(string), int(item["type"].(float64)), item["name"].(string), int64(item["number"].(float64))))
	}
	return n, items, nil
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
			s.dataService.Delete([]string{code})
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
	activeNum := len(res["active_usernames"].([]interface{}))
	s.dataService.Update(code, int64(res["id"].(float64)), res["title"].(string), desc, n, activeNum, "", 0)
}

// 获取群人数
func (s *BotManagerImp) GetChatMembers(code string) (uint32, error) {
	s.CheckSleep()
	params := map[string]string{"chat_id": code}
	res, err := s.Bot.Raw("getChatMemberCount", params)
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
	i := (rand.Intn(10) + 1) + rand.Intn(1000)/1000
	time.Sleep(time.Second * time.Duration(i))
}
