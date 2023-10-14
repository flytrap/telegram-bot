package human

import (
	"fmt"
	"strings"
)

func HumanSize(num int64) string {
	if num == 0 {
		return "View"
	}
	if num >= 1000000000 {
		return fmt.Sprintf("%.2fb", float32(num/1000000000))
	} else if num >= 1000000 {
		return fmt.Sprintf("%.2fm", float32(num/1000000))
	} else if num >= 1000 {
		return fmt.Sprintf("%.2fK", float32(num/1000))
	}
	return fmt.Sprintf("%d", num)
}

func TgGroupItemInfo(index int, code string, tp int, name string, num int64) string {
	_tp := "👥"
	if tp == 2 {
		_tp = "📢"
	}
	n := HumanSize(num)
	name = strings.ReplaceAll(name, "[", "【")
	name = strings.ReplaceAll(name, "]", "】")
	return fmt.Sprintf("%d、 %s [%s - %s](%s)", index, _tp, name, n, TgGroupUrl(code))
}

func TgGroupUrl(code string) string {
	return fmt.Sprintf("https://t.me/%s", code)
}

func DetailItemInfo(name string, desc string, extend string, location string, code string) string {
	if len(code) > 0 {
		name = fmt.Sprintf("[%s](@%s)", name, code)
	}
	return fmt.Sprintf("%s\n%s\n %s\n\n %s", name, location, desc, extend)
}
