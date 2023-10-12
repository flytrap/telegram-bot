package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/flytrap/telegram-bot/internal/config"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type IndexCallbackFunc func(context.Context, string, int64, int64) ([]string, bool, error)

func (s *HandlerManagerImp) IndexHandler(ctx tele.Context) error {
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

	size := int64(15)
	if config.C.Index.Detail {
		size = 1000
	}
	items, hasNext, err := s.ss.QueryItems(context.Background(), tag, page, size)
	localize := i18n.NewLocalizer(s.bundle, "zh-CN")
	keywordNotFound := localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "keywordNotFound"})
	if err != nil {
		logrus.Warning(err)
		return ctx.Reply(keywordNotFound)
	}
	var btnPrev tele.Btn
	var btnNext tele.Btn
	var btnDetail tele.Btn

	result := ""
	if config.C.Index.Detail && len(items) == 2 {
		viewPrivate := localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "viewPrivate"})
		another := localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "another"})
		btnDetail = selector.Data(viewPrivate, "view", tag, items[1])
		btnNext = selector.Data(another, "next", tag, "1")
		ctx.Bot().Handle(&btnDetail, s.DetailInfo)
		result = items[0]
	} else {
		result = strings.Join(items, "\n")
		if page > 1 {
			btnPrev = selector.Data("⬅ prev", "prev", tag, fmt.Sprint(page-1))
			ctx.Bot().Handle(&btnPrev, s.IndexHandler)
		}
		if hasNext {
			btnNext = selector.Data("next ➡", "next", tag, fmt.Sprint(page+1))
		}
	}
	if len(result) == 0 {
		return ctx.Reply(keywordNotFound)
	}

	selector.Inline(selector.Row(btnDetail), selector.Row(btnPrev, btnNext))
	ctx.Bot().Handle(&btnNext, s.IndexHandler)

	return ctx.EditOrReply(result, tele.ModeMarkdown, selector)
}

func (s *HandlerManagerImp) DetailInfo(ctx tele.Context) error {
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
