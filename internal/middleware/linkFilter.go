package middleware

import (
	"context"
	"regexp"
	"strings"

	"github.com/flytrap/telegram-bot/internal/config"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

var linkRe = regexp.MustCompile(`@[[a-zA-Z0-9_]]{5-64}`)

// Logger returns a middleware that logs incoming updates.
// If no custom logger provided, log.Default() will be used.
func (s *middleWareManagerImp) LinkFilter() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			user := c.Sender()
			q := c.Text()
			if checkLink(q) {
				logrus.Println("delete link", user.ID, user.Username, q)
				c.Delete()
				i := s.userService.AddWarning(context.Background(), user.ID)
				if i >= 3 {
					c.Bot().BanSenderChat(c.Chat(), c.Recipient())
				}
			}
			if config.C.Index.Recommend.Channel == c.Message().OriginalChat.Username {
				return nil
			}
			return next(c)
		}
	}
}

func checkLink(q string) bool {
	if len(q) == 0 {
		return true
	}
	if strings.Contains(q, "http://") {
		return true
	}
	if strings.Contains(q, "https://") {
		return true
	}
	if strings.Contains(q, "t.me") {
		return true
	}
	if len(linkRe.FindAllStringSubmatch(q, -1)) > 0 {
		return true
	}
	return false
}
