package serializers

import (
	"github.com/flytrap/telegram-bot/pkg/human"
)

type DataSerializer struct {
	Code   string `json:"Code"`
	Name   string `json:"name"`
	Number int64  `json:"number"`
	Type   int8   `json:"type"`
}

func (s *DataSerializer) ItemInfo(i int) string {
	return human.TgGroupItemInfo(i, s.Code, int(s.Type), s.Name, s.Number)
}

func (s *DataSerializer) Url() string {
	return human.TgGroupUrl(s.Code)
}
