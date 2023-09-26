package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/flytrap/telegram-bot/internal/config"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

func AfterDelTime() time.Duration {
	return time.Second * time.Duration(config.C.Bot.Manager.DeleteDelay)
}

func isAdmin(ctx tele.Context) bool {
	userId := ctx.Sender().ID
	if userId == 0 {
		return true
	}
	c := ctx.Chat()
	u := ctx.Sender()
	m, err := ctx.Bot().ChatMemberOf(c, u)
	if err != nil {
		logrus.Warning(err)
		return false
	}
	r := m.Role
	if r == tele.Creator || r == tele.Administrator {
		return true
	}
	return false
}

func getName(user *tele.User) string {
	text := ""
	if user != nil {
		if len(user.FirstName) > 0 {
			text += user.FirstName
		}
		if len(user.LastName) > 0 {
			text += user.LastName
		}
		if len(text) == 0 {
			text = user.Username
		}
	}
	if len(text) > 0 {
		return text
	}
	return "神秘人"
}

func (s *HandlerManagerImp) removetVerifyStatus(ChatId int64, userId int64) error {
	chat := tele.Chat{ID: ChatId}
	m, err := s.Bot.ChatMemberOf(&chat, &tele.User{ID: userId})
	if err != nil {
		logrus.Warning(err)
		return err
	}
	m.CanSendMessages = true
	m.CanSendMedia = true
	m.CanSendOther = true
	m.CanAddPreviews = true
	return s.Bot.Restrict(&chat, m)
}

// 发送消息并自动删除
func (s *HandlerManagerImp) sendAutoDeleteMessage(ctx tele.Context, d time.Duration, what interface{}, opts ...interface{}) error {
	m, _ := ctx.Bot().Send(ctx.Recipient(), what, opts...)
	key := fmt.Sprintf("del:message:%d", m.ID)
	data := fmt.Sprintf("%d:%d:%d", time.Now().Add(d).Unix(), m.Chat.ID, m.ID)
	s.store.Set(context.Background(), key, data, d+time.Hour*12)
	HandlerAfterFunc(d, func() {
		s.deleteMessage(key, m.Chat.ID, strconv.Itoa(m.ID))
	})
	return nil
}

// 发送消息并自动删除用户
func (s *HandlerManagerImp) sendAutoDeleteUser(ctx tele.Context, d time.Duration, what interface{}, opts ...interface{}) error {
	m, _ := ctx.Bot().Send(ctx.Recipient(), what, opts...)
	key := fmt.Sprintf("del:user:%d", m.ID)
	data := fmt.Sprintf("%d:%d:%d:%d", time.Now().Add(d).Unix(), m.Chat.ID, m.ID, m.Sender.ID)
	s.store.Set(context.Background(), key, data, d+time.Hour*12)
	HandlerAfterFunc(d, func() {
		if s.store.IsExist(context.Background(), key) {
			s.removeGroupMember(m.Chat.ID, m.Sender.ID)
		}
		s.deleteMessage(key, m.Chat.ID, strconv.Itoa(m.ID))
	})
	return nil
}

// 延时处理任务
func HandlerAfterFunc(d time.Duration, f func()) {
	time.AfterFunc(d, f)
}

func (s *HandlerManagerImp) deleteMessage(key string, chatId int64, msgId string) {
	msg := tele.StoredMessage{ChatID: chatId, MessageID: msgId}
	err := s.Bot.Delete(&msg)
	if err != nil {
		logrus.Warning(err)
	}
	if len(key) > 0 {
		s.store.Delete(context.Background(), key)
	}
}

// 踢人
func (s *HandlerManagerImp) removeGroupMember(ChatId int64, userId int64) error {
	chat := tele.Chat{ID: ChatId}
	m, err := s.Bot.ChatMemberOf(&chat, &tele.User{ID: userId})
	if err != nil {
		logrus.Warning(err)
		return err
	}
	return s.Bot.Ban(&chat, m)
}

// 检查遗留删除任务
func (s *HandlerManagerImp) CheckDeleteMessage(ctx context.Context) {
	results := s.store.Keys(ctx, "del:*")
	for _, item := range results {
		key := s.store.RawKey(item)
		delUser := strings.HasPrefix(key, "del:user:")
		v := s.store.Get(ctx, key)
		if v == nil {
			continue
		}
		li := strings.Split(v.(string), ":")
		if len(li) >= 3 {
			s.store.Delete(ctx, key)
			continue
		}
		i, err := strconv.ParseInt(li[0], 10, 64)
		if err != nil {
			logrus.Warning(err)
			continue
		}
		cid, err := strconv.ParseInt(li[1], 10, 64)
		if err != nil {
			logrus.Warning(err)
			continue
		}
		if i-10 < time.Now().Unix() {
			if delUser && len(li) == 4 {
				uid, err := strconv.ParseInt(li[3], 10, 64)
				if err == nil {
					s.removeGroupMember(cid, uid)
				}
			}
			s.deleteMessage(key, cid, li[2])
			continue
		}
		go HandlerAfterFunc(time.Second*time.Duration((i-time.Now().Unix())), func() {
			if delUser && len(li) == 4 {
				uid, err := strconv.ParseInt(li[3], 10, 64)
				if err == nil {
					s.removeGroupMember(cid, uid)
				}
			}
			s.deleteMessage(key, cid, li[2])
		})
	}
}
