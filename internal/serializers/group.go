package serializers

import (
	"github.com/flytrap/telegram-bot/pkg/human"
)

type GroupSerializer struct {
	Code   string `json:"Code"`
	Name   string `json:"name"`
	Number int64  `json:"number"`
	Type   int8   `json:"type"`
}

func (s *GroupSerializer) ItemInfo(i int) string {
	return human.TgGroupItemInfo(i, s.Code, string(rune(s.Type)), s.Name, s.Number)
}

func (s *GroupSerializer) Url() string {
	return human.TgGroupUrl(s.Code)
}
