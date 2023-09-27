package middleware

import (
	"context"

	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Logger returns a middleware that logs incoming updates.
// If no custom logger provided, log.Default() will be used.
func (s *middleWareManagerImp) Logger() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			user := c.Sender()
			q := c.Text()
			if c.Update().Message == nil && c.Update().Callback != nil {
				q = c.Update().Callback.Data
			}
			s.logUser(user, q)
			return next(c)
		}
	}
}

func (s *middleWareManagerImp) logUser(user *tele.User, text string) error {
	ctx := context.Background()
	if !s.userService.Check(ctx, user.ID) {
		s.userService.GetOrCreate(ctx, map[string]interface{}{"username": user.Username, "UserId": user.ID, "firstName": user.FirstName, "LastName": user.LastName,
			"LanguageCode": user.LanguageCode, "IsBot": user.IsBot, "IsPremium": user.IsPremium})
	}
	logrus.Println(user.ID, user.Username, text)
	return s.store.Xadd(ctx, "log:query", map[string]interface{}{"user_id": user.ID, "content": text, "type": "tg-query", "username": user.Username})
}
