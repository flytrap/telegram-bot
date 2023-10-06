package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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

	items, hasNext, err := s.ss.QueryItems(context.Background(), tag, page, 15)
	localize := i18n.NewLocalizer(s.bundle, "zh-CN")
	keywordNotFound := localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "keywordNotFound"})
	if err != nil {
		logrus.Warning(err)
		return ctx.Reply(keywordNotFound)
	}
	var btnPrev tele.Btn
	var btnNext tele.Btn
	if page > 1 {
		btnPrev = selector.Data("⬅ prev", "prev", tag, fmt.Sprint(page-1))
		ctx.Bot().Handle(&btnPrev, s.IndexHandler)
	}
	if hasNext {
		btnNext = selector.Data("next ➡", "next", tag, fmt.Sprint(page+1))
	}
	selector.Inline(selector.Row(btnPrev, btnNext))
	ctx.Bot().Handle(&btnNext, s.IndexHandler)
	result := strings.Join(items, "\n")
	if len(result) == 0 {
		return ctx.Reply(keywordNotFound)
	}
	return ctx.EditOrReply(result, tele.ModeMarkdown, selector)
}
