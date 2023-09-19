package human

import "fmt"

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
	return fmt.Sprintf("%d、 %s [%s - %s](%s)", index, _tp, name, n, TgGroupUrl(code))
}

func TgGroupUrl(code string) string {
	return fmt.Sprintf("https://t.me/%s", code)
}
