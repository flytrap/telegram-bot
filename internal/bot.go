package app

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/flytrap/telegram-bot/internal/config"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

func InitBot() (*tele.Bot, error) {
	token := config.C.Bot.Token
	if len(token) < 5 {
		token = os.Getenv("BotToken")
	}
	if len(token) < 10 {
		logrus.Warning("tg token not config")
		return nil, nil
	}
	c := http.Client{}
	if len(config.C.Proxy.Protocal) > 0 { //设置代理
		proxyURL, err := url.Parse(fmt.Sprintf("%s://%s:%d", config.C.Proxy.Protocal, config.C.Proxy.Ip, config.C.Proxy.Port))
		if err != nil {
			return nil, err
		}
		c = http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	}

	pref := tele.Settings{
		Token:  token, //  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
		Client: &c,
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// b.Handle("/start", sendMenu)
	// b.Handle("/chinese", sendMenu)

	return b, nil
}
