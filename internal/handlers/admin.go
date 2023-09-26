package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/flytrap/telegram-bot/internal/serializers"
	tele "gopkg.in/telebot.v3"
)

func (s *HandlerManagerImp) AdminHandler(ctx tele.Context) error {
	ctx.DeleteAfter(AfterDelTime())
	if ctx.Chat().Type != tele.ChatGroup && ctx.Chat().Type != tele.ChatSuperGroup {
		return s.sendAutoDeleteMessage(ctx, AfterDelTime(), "请将本机器人加入您的群中再使用此命令")
	}
	if !isAdmin(ctx) {
		return s.sendAutoDeleteMessage(ctx, AfterDelTime(), "只有本群管理才可以使用此命令")
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
	if len(req) > 0 {
		switch req {
		case "notrobot":
			return s.notrobot(ctx, gs)
		case "welcome":
			return s.welcome(ctx, gs)
		default:
			return s.sendAutoDeleteMessage(ctx, AfterDelTime(), "欢迎使用Admin命令\n\n输入 [admin welcome](/admin welcome) 设置欢迎词\n输入 [admin notRobot](/admin notrobot) 设置进群验证", tele.ModeMarkdown)
		}
	}
	bts := &tele.ReplyMarkup{}
	welcomeBt := bts.Data("欢迎词设置", "welcome", "welcome")
	notrbotBt := bts.Data("进群验证", "notrobot", "welcome")
	s.Bot.Handle(&welcomeBt, s.AdminHandler)
	s.Bot.Handle(&notrbotBt, s.AdminHandler)
	bts.Inline(bts.Row(welcomeBt, notrbotBt))
	return s.sendAutoDeleteMessage(ctx, AfterDelTime(), "欢迎使用Admin命令\n\n输入 [admin welcome](/admin welcome) 设置欢迎词\n输入 [admin notRobot](/admin notrobot) 设置进群验证", tele.ModeMarkdown, bts)
}

func (s *HandlerManagerImp) notrobot(ctx tele.Context, gs *serializers.GroupSetting) error {
	args := ctx.Args()
	if ctx.Callback() != nil {
		args = strings.Split(ctx.Callback().Data, "|")
	}
	text := ""
	timeout := -1
	isOpen := gs.NotRobot.IsOpen
	if len(args) > 2 && args[1] == "status" {
		if args[2] == "on" {
			gs.NotRobot.IsOpen = true
			isOpen = true
			text = "进群验证按钮启用成功"
		} else if args[2] == "off" {
			gs.NotRobot.IsOpen = false
			isOpen = false
			text = "进群验证按钮禁用成功"
		}
		s.gs.SetNotRobot(context.Background(), ctx.Chat().Username, gs.NotRobot)
	} else if len(args) > 2 && args[1] == "timeout" {
		timeout, err := strconv.Atoi(args[2])
		if err != nil {
			return s.sendAutoDeleteMessage(ctx, AfterDelTime(), "倒计时只能是数字，单位为妙钟")
		}
		gs.NotRobot.Timeout = timeout
		text = "进群验证倒计时已经设置为: " + args[2] + "\n注意，倒计时踢人需要在谷歌GAS设置触发器调用 removePendingUsers() 函数, 如果无法自动踢人请联系机器人服务器架设者"
		if timeout == 0 {
			text += "倒计时设置为0时，将不会自动踢除未验证的用户"
		}
		s.gs.SetNotRobot(context.Background(), ctx.Chat().Username, gs.NotRobot)
	} else {
		text = "使用以下命令设置进群验证\n\n"
		text += "/admin notrobot status - 开关进群验证，值可以是 on 或 off\n"
		text += "/admin notrobot timeout - 设置验证倒计时，可以是任何正数(秒)\n"
		text += "超过这个时间还没验证的用户会被踢出群. 设置成0则不踢\n"
		text += "注意！ 需要在谷歌脚本使用触发器定期踢出未验证的用户\n"
	}
	if len(text) > 0 {
		bts := s.setButtonsForAdminNotRobot(isOpen, timeout)
		return s.sendAutoDeleteMessage(ctx, AfterDelTime(), text, bts)
	}
	return nil
}

func (s *HandlerManagerImp) setButtonsForAdminNotRobot(isOpen bool, timeout int) *tele.ReplyMarkup {
	selector := tele.ReplyMarkup{}
	status := "【关】"
	statusTotal := "on"
	if isOpen {
		status = "【开】"
		statusTotal = "off"
	}
	selector.Row(selector.Data("进群验证按钮"+status, "joinGroup", "notrobot:status:"+statusTotal))
	if !isOpen {
		return &selector
	}
	ts := strconv.FormatInt(int64(timeout), 10)
	bts := []tele.Btn{selector.Data("踢出倒计时【"+ts+"秒】"+"✅", "notrobot", "notrobot", "timeout", ts)}
	if timeout != 120 {
		bts = append(bts, selector.Data("踢出倒计时【120秒】"+"✅", "notrobot", "notrobot", "timeout", "120"))
	}
	if timeout != 3600 {
		bts = append(bts, selector.Data("踢出倒计时【1小时】"+"✅", "notrobot", "notrobot", "timeout", "3600"))
	}
	if timeout != 3600*24 {
		bts = append(bts, selector.Data("踢出倒计时【1天】"+"✅", "notrobot", "notrobot", "timeout", "86400"))
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
	if len(args) < 3 {
		text := ""
		if len(args) == 2 && args[1] == "customtext" {
			text = "要设置自定义欢迎词，请回复本消息把欢迎词发给我，你可以在欢迎词中使用以下变量：\n"
			text += "{{.desc}} - 代表群描述\n"
			text += "{{.pinned}} - 代表置顶消息\n"
			text += "{{.name}} - 代表新成员的名字\n"
		} else {
			text = "在这里可以打开欢迎词，欢迎词可选包含群描述和置顶消息，如果30秒自毁开启，欢迎词将在30秒后被删除\n"
			text += "请问还有什么我可以帮你的吗？\n"
			text += "/admin welcome status on/of 开关欢迎词\n"
			text += "/admin welcome desc on/of 是否显示群描述信息\n"
			text += "/admin welcome pinned on/of 是否显示群置顶信息\n"
			text += "/admin welcome killme 0 欢迎消息销毁停留时间，0不删除\n"
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
		s.sendAutoDeleteMessage(ctx, AfterDelTime(), "设置失败 \n")
		return nil
	}
	s.sendAutoDeleteMessage(ctx, AfterDelTime(), fmt.Sprintf("设置成功 - %s: %s\n", args[1], value))
	return nil
}
