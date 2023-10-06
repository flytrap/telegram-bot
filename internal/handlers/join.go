package handlers

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/flytrap/telegram-bot/internal/serializers"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// 新用户加入群组
func (s *HandlerManagerImp) JoinInHandler(ctx tele.Context) error {
	ctx.Delete()
	if menu == nil {
		initMenu()
	}
	localize := i18n.NewLocalizer(s.bundle, "zh-CN")
	if ctx.Chat().Type != tele.ChatGroup && ctx.Chat().Type != tele.ChatSuperGroup {
		return s.sendAutoDeleteMessage(ctx, AfterDelTime(), localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "admin.joinTip"}), menu)
	}

	gs, _ := s.gs.GetSetting(context.Background(), ctx.Chat().Username)
	if gs.NotRobot.IsOpen {
		return s.verifyPayload(ctx, gs)
	} else {
		return s.welcomePayload(ctx, gs)
	}
}

func (s *HandlerManagerImp) verifyPayload(ctx tele.Context, gs *serializers.GroupSetting) error {
	ctx.ChatMember()
	cm, err := ctx.Bot().ChatMemberOf(ctx.Chat(), ctx.Bot().Me)
	if err != nil {
		return err
	}
	timeout := 120
	localize := i18n.NewLocalizer(s.bundle, "zh-CN")
	if cm.CanRestrictMembers {
		m, err := ctx.Bot().ChatMemberOf(ctx.Chat(), ctx.Sender())
		if err != nil {
			logrus.Warning(err)
			return err
		}
		m.CanSendMessages = false
		m.CanSendMedia = false
		m.CanSendOther = false
		m.CanAddPreviews = false
		err = ctx.Bot().Restrict(ctx.Chat(), m) // 禁言
		if err != nil {
			return err
		}
		text := localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "joinTip", TemplateData: map[string]string{"Name": getName(ctx.Sender())}})
		if gs.NotRobot.Timeout > 0 {
			timeout = gs.NotRobot.Timeout
			text += localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "joinTip2", TemplateData: map[string]int{"Timeout": timeout}})
		}
		selector := &tele.ReplyMarkup{}
		bt := selector.Data(localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "applyGroup"}), "joinGroup", fmt.Sprintf("%d", ctx.Sender().ID))
		selector.Inline(selector.Row(bt))
		ctx.Bot().Handle(&bt, func(c tele.Context) error {
			i, err := strconv.ParseInt(strings.Split(c.Callback().Data, "|")[0], 10, 64)
			if err != nil {
				return err
			}
			if c.Sender().ID == i || isAdmin(ctx) {
				s.removeVerifyStatus(ctx.Chat().ID, c.Sender().ID)
				return s.welcomePayload(ctx, gs)
			}
			return nil
		})
		return s.sendAutoDeleteUser(ctx, time.Second*time.Duration(timeout), text, selector, menu)
	} else {
		return s.sendAutoDeleteMessage(ctx, AfterDelTime(), localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "joinVerifyError"}), menu)
	}
}

func (s *HandlerManagerImp) welcomePayload(ctx tele.Context, gs *serializers.GroupSetting) error {
	if !gs.Welcome.IsOpen {
		return nil
	}
	chat := ctx.Chat()
	localize := i18n.NewLocalizer(s.bundle, "zh-CN")
	data := map[string]string{"Name": getName(ctx.Sender()), "Desc": "", "Pinned": ""}
	if gs.Welcome.Desc && len(chat.Description) > 0 {
		desc := chat.Description
		if len(gs.Welcome.Template) == 0 {
			desc = localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "descTemplate", TemplateData: map[string]string{"Desc": chat.Description}})
		}
		data["desc"] = desc
	}
	if gs.Welcome.Pinned && chat.PinnedMessage != nil && len(chat.PinnedMessage.Text) > 0 {
		pinned := fmt.Sprintf("[%s](https://t.me/%s/%d)", chat.PinnedMessage.Text, chat.Username, chat.PinnedMessage.ID)
		if len(gs.Welcome.Template) == 0 {
			data["pinned"] = localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "pinnedTemplate", TemplateData: map[string]string{"Pinned": pinned}})
		}
	}
	tp := localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "defaultTemplate"})
	if len(gs.Welcome.Template) > 0 {
		tp = gs.Welcome.Template
	}
	result, err := template.New("welcome").Parse(tp)
	if err != nil {
		return err
	}
	out := new(bytes.Buffer)
	err = result.Execute(out, data) // 渲染模版
	if err != nil {
		return err
	}
	text := ""
	if gs.Welcome.KillMe > 0 {
		text = localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "killMsg", TemplateData: map[string]int{"Timeout": gs.Welcome.KillMe}})
		return s.sendAutoDeleteMessage(ctx, time.Duration(gs.Welcome.KillMe)*time.Second, text+out.String(), menu)
	}
	return ctx.Send(text+out.String(), menu)
}
