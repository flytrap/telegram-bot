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
	category, tag, q, page, size := parseArgs(ctx)
	items, hasNext, err := s.ss.QueryItems(context.Background(), category, tag, q, page, size)
	localize := i18n.NewLocalizer(s.bundle, "zh-CN")
	keywordNotFound := localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "keywordNotFound"})
	if err != nil || len(items) == 0 {
		logrus.Warning(err)
		return ctx.Reply(keywordNotFound)
	}
	var btnPrev tele.Btn
	var btnNext tele.Btn

	ad := s.ss.LoadAd(q)
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
			rows = append(rows, selector.Row(s.detailItem(selector, ii, tag, q, item["code"].(string), item["name"].(string))))
		}
	}
	result := strings.Join(results, "\n")
	if page > 1 {
		btnPrev = selector.Data("⬅ prev", "prev", category, tag, q, fmt.Sprint(page-1))
		ctx.Bot().Handle(&btnPrev, s.IndexHandler)
	}
	if hasNext {
		btnNext = selector.Data("next ➡", "next", category, tag, q, fmt.Sprint(page+1))
	}
	if len(result) == 0 {
		result = localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "queryTip"})
	}
	rows = append(rows, selector.Row(btnPrev, btnNext))

	selector.Inline(rows...)
	ctx.Bot().Handle(&btnNext, s.IndexHandler)

	return ctx.EditOrReply(result, tele.ModeMarkdown, selector)
}

func parseArgs(ctx tele.Context) (string, string, string, int64, int64) {
	category, tag := "", "" // 分类
	q := ctx.Message().Text // 关键词
	page := int64(1)
	args := ctx.Args()
	cb := ctx.Callback()
	if cb != nil {
		args = strings.Split(cb.Data, "|")
		if len(args) > 3 {
			n, err := strconv.Atoi(args[3])
			if err != nil {
				logrus.Warn(err)
			} else {
				page = int64(n)
			}
		}
	}

	if len(args) > 2 {
		category = args[0]
		tag = args[1]
		q = args[2]
	} else if len(args) > 1 {
		tag = args[0]
		q = args[1]
	} else if len(args) == 1 {
		q = args[0]
	}
	return category, tag, q, page, config.C.Index.PageSize
}

// 列表条目
func (s *HandlerManagerImp) detailItem(selector *tele.ReplyMarkup, i int, tag string, q string, code string, name string) tele.Btn {
	btnDetail := selector.Data(fmt.Sprintf("%d: %s", i, name), "detail", tag, q, code)
	s.Bot.Handle(&btnDetail, s.detailInfo)
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
	if len(args) != 3 {
		return ctx.Reply(keywordNotFound)
	}
	tag, q, code := args[0], args[1], args[2]
	item, err := s.ss.GetDetail(context.Background(), code)
	if err != nil {
		return ctx.Reply(keywordNotFound)
	}
	selector := &tele.ReplyMarkup{}
	viewPrivate := localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "viewPrivate"})
	btnDetail := selector.Data(viewPrivate, "view", tag, q, code)
	ctx.Bot().Handle(&btnDetail, s.PrivateInfo)
	selector.Inline(selector.Row(btnDetail))
	result := human.DetailItemInfo(item["name"].(string), item["desc"].(string), item["extend"].(string), item["location"].(string), "")

	ps := tele.Album{} // 多张图合并
	for _, img := range item["images"].([]interface{}) {
		f, err := human.Base64ToIoReader(img.(string))
		if err == nil {
			ps = append(ps, &tele.Photo{File: tele.FromReader(f)})
		}
	}
	err = ctx.SendAlbum(ps)
	if err != nil {
		logrus.Warn(err)
	}

	return s.sendAutoDeleteMessage(ctx, AfterDelTime(), result, tele.ModeMarkdown, selector)
}

// 获取隐私信息
func (s *HandlerManagerImp) PrivateInfo(ctx tele.Context) error {
	cb := ctx.Callback()
	localize := i18n.NewLocalizer(s.bundle, "zh-CN")
	keywordNotFound := localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "keywordNotFound"})
	if cb != nil {
		li := strings.Split(cb.Data, "|")
		code := li[len(li)-1]
		result, err := s.ss.GetPrivate(context.Background(), code)
		if err != nil {
			return ctx.Reply(keywordNotFound)
		}
		return s.sendAutoDeleteMessage(ctx, AfterDelTime(), result, tele.ModeMarkdown)
	}
	return ctx.Reply(keywordNotFound)
}
