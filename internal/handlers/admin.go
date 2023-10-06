package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/flytrap/telegram-bot/internal/serializers"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	tele "gopkg.in/telebot.v3"
)

func (s *HandlerManagerImp) AdminHandler(ctx tele.Context) error {
	ctx.DeleteAfter(AfterDelTime())
	localize := i18n.NewLocalizer(s.bundle, "zh-CN")
	if ctx.Chat().Type != tele.ChatGroup && ctx.Chat().Type != tele.ChatSuperGroup {
		adminJoinTip := localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "admin.joinTip"})
		return s.sendAutoDeleteMessage(ctx, AfterDelTime(), adminJoinTip)
	}
	if !isAdmin(ctx) {
		noPermissionTio := localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "admin.noPermissionTio"})
		return s.sendAutoDeleteMessage(ctx, AfterDelTime(), noPermissionTio)
	}
	gs, err := s.gs.GetSetting(context.Background(), ctx.Chat().Username)
	if err != nil {
		return nil
	}
	req := ""
	if len(ctx.Args()) > 0 && len(ctx.Args()[0]) > 0 {
		req = ctx.Args()[0]
	}
	if ctx.Callback() != nil {
		req = ctx.Callback().Unique
	}
	helpTip := localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "admin.help"})
	if len(req) > 0 {
		switch req {
		case "notrobot":
			return s.notrobot(ctx, gs)
		case "welcome":
			return s.welcome(ctx, gs)
		default:
			return s.sendAutoDeleteMessage(ctx, AfterDelTime(), helpTip, tele.ModeMarkdown)
		}
	}
	bts := &tele.ReplyMarkup{}
	welcomeTip := localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "admin.welcome"})
	notrobotTip := localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "admin.notrobot"})
	welcomeBt := bts.Data(welcomeTip, "welcome", "welcome")
	notrbotBt := bts.Data(notrobotTip, "notrobot", "notrobot")
	s.Bot.Handle(&welcomeBt, s.AdminHandler)
	s.Bot.Handle(&notrbotBt, s.AdminHandler)
	bts.Inline(bts.Row(welcomeBt, notrbotBt))
	return s.sendAutoDeleteMessage(ctx, AfterDelTime(), helpTip, tele.ModeMarkdown, bts)
}

func (s *HandlerManagerImp) notrobot(ctx tele.Context, gs *serializers.GroupSetting) error {
	args := ctx.Args()
	if ctx.Callback() != nil {
		args = strings.Split(ctx.Callback().Data, "|")
	}
	text := ""
	timeout := -1
	isOpen := gs.NotRobot.IsOpen
	localize := i18n.NewLocalizer(s.bundle, "zh-CN")
	if len(args) > 2 && args[1] == "status" {
		if args[2] == "on" {
			gs.NotRobot.IsOpen = true
			isOpen = true
			text = localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "admin.notrobotOn"})
		} else if args[2] == "off" {
			gs.NotRobot.IsOpen = false
			isOpen = false
			text = localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "admin.notrobotOff"})
		}
		s.gs.SetNotRobot(context.Background(), ctx.Chat().Username, gs.NotRobot)
	} else if len(args) > 2 && args[1] == "timeout" {
		timeout, err := strconv.Atoi(args[2])
		if err != nil {
			return s.sendAutoDeleteMessage(ctx, AfterDelTime(), localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "admin.timeoutError"}))
		}
		gs.NotRobot.Timeout = timeout
		text = localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "admin.timeoutSuccess", TemplateData: map[string]string{"Timeout": args[2]}})
		if timeout == 0 {
			text += localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "admin.timeoutZero"})
		}
		s.gs.SetNotRobot(context.Background(), ctx.Chat().Username, gs.NotRobot)
	} else {
		text = localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "admin.helpDetail"})
	}
	if len(text) > 0 {
		bts := s.setButtonsForAdminNotRobot(isOpen, timeout)
		return s.sendAutoDeleteMessage(ctx, AfterDelTime(), text, bts)
	}
	return nil
}

func (s *HandlerManagerImp) setButtonsForAdminNotRobot(isOpen bool, timeout int) *tele.ReplyMarkup {
	selector := tele.ReplyMarkup{}
	localize := i18n.NewLocalizer(s.bundle, "zh-CN")
	status := localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "admin.off"})
	statusTotal := "on"
	if isOpen {
		status = localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "admin.on"})
		statusTotal = "off"
	}
	selector.Row(selector.Data(localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "admin.groupJoinButton"})+status, "joinGroup", "notrobot:status:"+statusTotal))
	if !isOpen {
		return &selector
	}
	ts := strconv.FormatInt(int64(timeout), 10)
	bts := []tele.Btn{selector.Data(localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "admin.banTimeoutS", TemplateData: map[string]string{"Timeout": ts}}), "notrobot", "notrobot", "timeout", ts)}
	if timeout != 120 {
		bts = append(bts, selector.Data(localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "admin.banTimeoutS", TemplateData: map[string]string{"Timeout": "120"}}), "notrobot", "notrobot", "timeout", "120"))
	}
	if timeout != 3600 {
		bts = append(bts, selector.Data(localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "admin.banTimeoutH", TemplateData: map[string]string{"Timeout": "1"}}), "notrobot", "notrobot", "timeout", "3600"))
	}
	if timeout != 3600*24 {
		bts = append(bts, selector.Data(localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "admin.banTimeoutD", TemplateData: map[string]string{"Timeout": "1"}}), "notrobot", "notrobot", "timeout", "86400"))
	}
	rows := []tele.Row{}
	for _, bt := range bts {
		s.Bot.Handle(&bt, s.AdminHandler)
		rows = append(rows, selector.Row(bt))
	}
	selector.Inline(rows...)
	return &selector
}

func (s *HandlerManagerImp) welcome(ctx tele.Context, gs *serializers.GroupSetting) error {
	args := ctx.Args()
	if ctx.Callback() != nil {
		args = strings.Split(ctx.Callback().Data, "|")
	}
	localize := i18n.NewLocalizer(s.bundle, "zh-CN")
	if len(args) < 3 {
		text := ""
		if len(args) == 2 && args[1] == "customtext" {
			text = localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "admin.welcomeCustom"})
		} else {
			text = localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "admin.welcomeHelp"})
		}
		if len(text) > 0 {
			return s.sendAutoDeleteMessage(ctx, time.Second*30, text)
		}
	}
	value := args[2]
	status := map[string]bool{"on": true, "off": false}[value]
	switch args[1] {
	case "status":
		gs.Welcome.IsOpen = status
	case "desc":
		gs.Welcome.Desc = status
	case "pinned":
		gs.Welcome.Pinned = status
	case "killme":
		i, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		gs.Welcome.KillMe = i
	case "customtext":
		gs.Welcome.Template = value
	}
	err := s.gs.SetWelcome(context.Background(), ctx.Chat().Username, gs.Welcome)
	if err != nil {
		s.sendAutoDeleteMessage(ctx, AfterDelTime(), localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "admin.setFailed"}))
		return nil
	}
	v := localize.MustLocalize(&i18n.LocalizeConfig{MessageID: "admin.setSuccess", TemplateData: map[string]string{"Value": fmt.Sprintf("%s: %s\n", args[1], value)}})
	s.sendAutoDeleteMessage(ctx, AfterDelTime(), v)
	return nil
}
