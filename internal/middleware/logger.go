package middleware

import (
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type HandlerUserFunc func(map[string]interface{}) error
type HandlerUserQueryFunc func(int64, string) error

// Logger returns a middleware that logs incoming updates.
// If no custom logger provided, log.Default() will be used.
func Logger(userFunc HandlerUserFunc, userQueryFunc HandlerUserQueryFunc) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			user := c.Sender()
			go userFunc(map[string]interface{}{"username": user.Username, "UserId": user.ID, "firstName": user.FirstName, "LastName": user.LastName,
				"LanguageCode": user.LanguageCode, "IsBot": user.IsBot, "IsPremium": user.IsPremium})
			q := c.Text()
			if c.Update().Message == nil && c.Update().Callback != nil {
				q = c.Update().Callback.Data
			}
			go userQueryFunc(user.ID, q)
			logUser(user, q)
			return next(c)
		}
	}
}

func logUser(u *tele.User, text string) error {
	logrus.Println(u.ID, u.Username, text)
	return nil
}
