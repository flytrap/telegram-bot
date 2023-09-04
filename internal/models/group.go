package models

import (
	"github.com/flytrap/telegram-bot/pkg/human"
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model

	Category uint   `json:"category" gorm:"comment:分类"`
	Language string `json:"language" gorm:"size:32;comment:语言"`
	Name     string `json:"name" gorm:"size:256;comment:名字"`
	Tid      int64  `json:"tid" gorm:"comment:TgId"`
	Type     int8   `json:"type" gorm:"default:1;comment:类型，区分组|channel"`
	Code     string `json:"code" gorm:"unique;size:64;comment:唯一标识"`
	Number   uint32 `json:"number" gorm:"comment:人数"`
	Desc     string `json:"desc" gorm:"type:text;comment:描述信息"`
	Weight   int32  `json:"weight" gorm:"default:0;comment:权重"`

	Tags []*Tag `json:"tags" gorm:"many2many:group_tag"`
}

func (s *Group) HumanSize() string {
	return human.HumanSize(int64(s.Number))
}
