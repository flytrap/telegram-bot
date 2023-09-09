package middleware

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Logger returns a middleware that logs incoming updates.
// If no custom logger provided, log.Default() will be used.
func Logger(logger ...*logrus.Logger) tele.MiddlewareFunc {
	var l *logrus.Logger
	if len(logger) > 0 {
		l = logger[0]
	} else {
		l = logrus.StandardLogger()
	}
	logrus.Println()

	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			data, _ := json.MarshalIndent(c.Update(), "", "  ")
			logUser(c.Sender(), c.Text())
			l.Println(string(data))
			return next(c)
		}
	}
}

func logUser(u *tele.User, text string) error {
	return nil
}
