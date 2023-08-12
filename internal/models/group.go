package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Group struct {
	gorm.Model

	Category uint   `json:"category" gorm:"comment:分类"`
	Name     string `json:"name" gorm:"size:256;comment:名字"`
	Tid      uint64 `json:"tid" gorm:"comment:TgId"`
	Type     int8   `json:"type" gorm:"default:1;comment:类型，区分组|channel"`
	Code     string `json:"code" gorm:"size:32;comment:唯一标识"`
	Number   uint32 `json:"number" gorm:"comment:人数"`
	Desc     string `json:"desc" gorm:"type:text;comment:描述信息"`

	Tags []*Tag `json:"tags" gorm:"many2many:group_tag"`
}

func (s *Group) HumanSize() string {
	if s.Number >= 1000000000 {
		return fmt.Sprintf("%.2fb", float32(s.Number/1000000000))
	} else if s.Number >= 1000000 {
		return fmt.Sprintf("%.2fm", float32(s.Number/1000000))
	} else if s.Number >= 1000 {
		return fmt.Sprintf("%.2fK", float32(s.Number/1000))
	}
	return fmt.Sprintf("%d", s.Number)
}
