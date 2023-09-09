package middleware

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type HandlerUserFunc func(map[string]interface{}) error

// Logger returns a middleware that logs incoming updates.
// If no custom logger provided, log.Default() will be used.
func Logger(userFunc HandlerUserFunc, logger ...*logrus.Logger) tele.MiddlewareFunc {
	var l *logrus.Logger
	if len(logger) > 0 {
		l = logger[0]
	} else {
		l = logrus.StandardLogger()
	}

	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			data, _ := json.MarshalIndent(c.Update(), "", "  ")
			user := c.Sender()
			go userFunc(map[string]interface{}{"Username": user.Username, "userID": user.ID, "FirstName": user.FirstName, "LastName": user.LastName,
				"LanguageCode": user.LanguageCode, "IsBot": user.IsBot, "IsPremium": user.IsPremium})
			logUser(user, c.Text())
			l.Println(string(data))
			return next(c)
		}
	}
}

func logUser(u *tele.User, text string) error {
	return nil
}
