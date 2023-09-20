package serializers

import (
	"strings"

	"github.com/flytrap/telegram-bot/pkg/human"
)

type DataSerializer struct {
	Code   string `json:"Code"`
	Name   string `json:"name"`
	Number int64  `json:"number"`
	Type   int8   `json:"type"`
}

func (s *DataSerializer) ItemInfo(i int) string {
	name := strings.ReplaceAll(s.Name, "[", "|")
	name = strings.ReplaceAll(name, "]", "|")
	name = strings.ReplaceAll(name, "@", "__")
	return human.TgGroupItemInfo(i, s.Code, int(s.Type), name, s.Number)
}

func (s *DataSerializer) Url() string {
	return human.TgGroupUrl(s.Code)
}
