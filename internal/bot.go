package app

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/flytrap/telegram-bot/internal/config"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

var (
	menu = &tele.ReplyMarkup{ResizeKeyboard: true}
)

func InitBot() (*tele.Bot, error) {
	if len(config.C.Bot.Token) < 10 {
		logrus.Warning("tg token not config")
		return nil, nil
	}
	initMenu() // 初始化菜单项
	c := http.Client{}
	if len(config.C.Proxy.Protocal) > 0 { //设置代理
		proxyURL, err := url.Parse(fmt.Sprintf("%s://%s:%d", config.C.Proxy.Protocal, config.C.Proxy.Ip, config.C.Proxy.Port))
		if err != nil {
			return nil, err
		}
		c = http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	}

	pref := tele.Settings{
		Token:  config.C.Bot.Token, //  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
		Client: &c,
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	b.Handle("/start", sendMenu)

	return b, nil
}

func initMenu() {
	items := []tele.Row{}
	for _, item := range config.C.Bot.Menus {
		subItem := []tele.Btn{}
		for _, su := range item {
			subItem = append(subItem, menu.Text(su))
		}
		items = append(items, menu.Row(subItem...))
	}
	menu.Reply(items...)
}

func sendMenu(c tele.Context) error {
	return c.Send(config.C.Bot.Start, menu, tele.ModeMarkdown)
}
