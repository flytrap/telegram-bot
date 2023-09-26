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
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

var welcomeTemplate = "本BOT代表本群所有人热烈欢迎新成员: {{.name}}\n {{.desc}}\n\n {{.pinned}}"

// 新用户加入群组
func (s *HandlerManagerImp) JoinInHandler(ctx tele.Context) error {
	ctx.Delete()
	if ctx.Chat().Type != tele.ChatGroup && ctx.Chat().Type != tele.ChatSuperGroup {
		return s.sendAutoDeleteMessage(ctx, AfterDelTime(), "请将本机器人加入您的群中再使用此命令")
	}

	gs, _ := s.gs.GetSetting(context.Background(), ctx.Chat().Username)
	if gs.NotRobot.IsOpen {
		s.welcomePayload(ctx, gs)
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
		text := "进群验证已启用\n"
		text += "您好! " + getName(ctx.Sender()) + "!\n"
		text += " 别忘了点击以下按钮获取发言权限!\n"
		if gs.NotRobot.Timeout > 0 {
			timeout = gs.NotRobot.Timeout
			text += fmt.Sprintf("如果 *%d* 秒内你没有点击以下按钮，你将被踢出群，你可以在一分钟后重新加入, 如果无法加入请重启Telegram\n", gs.NotRobot.Timeout)
			text += "注: 管理员点以下按钮也可放行"
		}
		selector := &tele.ReplyMarkup{ResizeKeyboard: true}
		bt := selector.Data("🈸"+" - 申请入群", "joinGroup", fmt.Sprintf("%d", ctx.Sender().ID))
		selector.Inline(selector.Row(bt))
		ctx.Bot().Handle(&bt, func(c tele.Context) error {
			i, err := strconv.ParseInt(strings.Split(c.Callback().Data, "|")[0], 10, 64)
			if err != nil {
				return err
			}
			if c.Sender().ID == i || isAdmin(ctx) {
				return s.removetVerifyStatus(ctx.Chat().ID, c.Sender().ID)
			}
			return nil
		})
		return s.sendAutoDeleteUser(ctx, time.Second*time.Duration(timeout), text, selector)
	} else {
		return s.sendAutoDeleteMessage(ctx, AfterDelTime(), "本机器人打开了进群验证功能，但是没有管理员权限，无法禁用用户\n")
	}
}

func (s *HandlerManagerImp) welcomePayload(ctx tele.Context, gs *serializers.GroupSetting) error {
	if !gs.Welcome.IsOpen {
		return nil
	}
	chat := ctx.Chat()
	data := map[string]string{"name": getName(ctx.Sender()), "desc": "", "pinned": ""}
	if gs.Welcome.Desc && len(chat.Description) > 0 {
		desc := chat.Description
		if len(gs.Welcome.Template) == 0 {
			desc = fmt.Sprintf("\n请遵循本群规则\n%s\n", chat.Description)
		}
		data["desc"] = desc
	}
	if gs.Welcome.Pinned && chat.PinnedMessage != nil && len(chat.PinnedMessage.Text) > 0 {
		pinned := fmt.Sprintf("[%s](https://t.me/%s/%d)", chat.PinnedMessage.Text, chat.Username, chat.PinnedMessage.ID)
		if len(gs.Welcome.Template) == 0 {
			data["pinned"] = fmt.Sprintf("\n请务必读一下置顶消息\n%s", pinned)
		}
	}
	tp := welcomeTemplate
	if len(gs.Welcome.Template) > 0 {
		tp = gs.Welcome.Template
	}
	result, err := template.New("welecome").Parse(tp)
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
		text = fmt.Sprintf("本消息将在%d秒后自毁\n", gs.Welcome.KillMe)
		return s.sendAutoDeleteMessage(ctx, time.Duration(gs.Welcome.KillMe)*time.Second, text+out.String())
	}
	return ctx.Send(text + out.String())
}
