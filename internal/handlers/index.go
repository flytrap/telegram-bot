package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/flytrap/telegram-bot/internal/config"
	"github.com/flytrap/telegram-bot/pkg/human"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type IndexCallbackFunc func(context.Context, string, int64, int64) ([]string, bool, error)

func (s *HandlerManagerImp) IndexHandler(ctx tele.Context) error {
	selector := &tele.ReplyMarkup{}
	category := ""            // 分类
	tag := ctx.Message().Text // 关键词
	page := int64(1)
	args := ctx.Args()
	cb := ctx.Callback()
	if cb != nil {
		args = strings.Split(cb.Data, "|")
		if len(args) > 2 {
			n, err := strconv.Atoi(args[1])
			if err != nil {
				logrus.Warn(err)
			} else {
				page = int64(n)
			}
		}
	}

	if len(args) > 0 {
		tag = args[0]
	}
	if len(args) > 1 {
		category = args[0]
		tag = args[1]
	}

	size := int64(15)
	items, hasNext, err := s.ss.QueryItems(context.Background(), category, tag, page, size)
	localize := i18n.NewLocalizer(s.bundle, "zh-CN")
	keywordNotFound := localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "keywordNotFound"})
	if err != nil || len(items) == 0 {
		logrus.Warning(err)
		return ctx.Reply(keywordNotFound)
	}
	var btnPrev tele.Btn
	var btnNext tele.Btn

	ad := s.ss.LoadAd(tag)
	rows := []tele.Row{}
	results := []string{}
	if len(ad) > 0 {
		results = append(results, ad)
	}
	for i, item := range items {
		ii := int(page-1)*int(size) + i + 1
		if config.C.Index.ItemMode == "tg_link" {
			results = append(results, human.TgGroupItemInfo(ii, item["code"].(string), int(item["type"].(float64)), item["name"].(string), int64(item["number"].(float64))))
		} else {
			rows = append(rows, selector.Row(s.detailItem(ctx, selector, ii, tag, item["code"].(string), item["name"].(string))))
		}
	}
	result := strings.Join(results, "\n")
	if page > 1 {
		btnPrev = selector.Data("⬅ prev", "prev", tag, category, fmt.Sprint(page-1))
		ctx.Bot().Handle(&btnPrev, s.IndexHandler)
	}
	if hasNext {
		btnNext = selector.Data("next ➡", "next", tag, category, fmt.Sprint(page+1))
	}
	if len(result) == 0 {
		result = localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "queryTip"})
	}
	rows = append(rows, selector.Row(btnPrev, btnNext))

	selector.Inline(rows...)
	ctx.Bot().Handle(&btnNext, s.IndexHandler)

	return ctx.EditOrReply(result, tele.ModeMarkdown, selector)
}

func (s *HandlerManagerImp) detailItem(ctx tele.Context, selector *tele.ReplyMarkup, i int, tag string, code string, name string) tele.Btn {
	btnDetail := selector.Data(fmt.Sprintf("%d: %s", i, name), "detail", tag, code)
	ctx.Bot().Handle(&btnDetail, s.detailInfo)
	return btnDetail
}

// 获取详细信息
func (s *HandlerManagerImp) detailInfo(ctx tele.Context) error {
	cb := ctx.Callback()
	localize := i18n.NewLocalizer(s.bundle, "zh-CN")
	keywordNotFound := localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "keywordNotFound"})
	if cb == nil {
		return ctx.Reply(keywordNotFound)
	}
	args := strings.Split(cb.Data, "|") // 获取回调参数
	if len(args) != 2 {
		return ctx.Reply(keywordNotFound)
	}
	tag, code := args[0], args[1]
	item, err := s.ss.GetDetail(context.Background(), code)
	if err != nil {
		return ctx.Reply(keywordNotFound)
	}
	selector := &tele.ReplyMarkup{}
	viewPrivate := localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "viewPrivate"})
	btnDetail := selector.Data(viewPrivate, "view", tag, code)
	ctx.Bot().Handle(&btnDetail, s.PrivateInfo)
	selector.Inline(selector.Row(btnDetail))
	result := human.DetailItemInfo(item["name"].(string), item["desc"].(string), item["extend"].(string), item["images"].([]interface{}), "")
	return s.sendAutoDeleteMessage(ctx, AfterDelTime(), result, tele.ModeMarkdown, selector)
}

// 获取隐私信息
func (s *HandlerManagerImp) PrivateInfo(ctx tele.Context) error {
	cb := ctx.Callback()
	localize := i18n.NewLocalizer(s.bundle, "zh-CN")
	keywordNotFound := localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "keywordNotFound"})
	if cb != nil {
		code := strings.Split(cb.Data, "|")[1]
		result, err := s.ss.GetPrivate(context.Background(), code)
		if err != nil {
			return ctx.Reply(keywordNotFound)
		}
		return s.sendAutoDeleteMessage(ctx, AfterDelTime(), result, tele.ModeMarkdown)
	}
	return ctx.Reply(keywordNotFound)
}
