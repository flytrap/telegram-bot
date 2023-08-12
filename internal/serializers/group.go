package serializers

import "fmt"

type GroupSerilizer struct {
	Code   string `json:"Code"`
	Name   string `json:"name"`
	Number string `json:"number"`
	Type   int8   `json:"type"`
}

func (s *GroupSerilizer) ItemInfo(i int) string {
	tp := "👥"
	if s.Type == 2 {
		tp = "📢"
	}
	return fmt.Sprintf("%d、 %s [%s - %s](%s)", i, tp, s.Name, s.Number, s.Url())
}

func (s *GroupSerilizer) Url() string {
	return fmt.Sprintf("https://t.me/%s", s.Code)
}
