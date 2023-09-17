package models

import (
	"time"

	"gorm.io/gorm"
)

// 广告
type Ad struct {
	gorm.Model

	Keyword  string    `json:"keyword" gorm:"size:64;comment:关键词"`
	Global   int8      `json:"global" gorm:"default:0;comment:全局词条,0非全局"`
	IsShowAd bool      `json:"is_show_ad" gorm:"default:false;comment:是否显示广告字样"`
	Expire   time.Time `json:"expire" gorm:"comment:过期时间"`
	Category uint      `json:"category" gorm:"comment:分类"`
	Language string    `json:"language" gorm:"size:32;comment:语言"`
	Name     string    `json:"name" gorm:"index;size:256;comment:名字"`
	Code     string    `json:"code" gorm:"unique;size:64;comment:唯一标识"`
	Type     int8      `json:"type" gorm:"default:1;comment:类型,区分group|channel"`
	Number   uint32    `json:"number" gorm:"comment:人数"`
	Desc     string    `json:"desc" gorm:"type:text;comment:描述信息"`
}

func (item *Ad) ToMap() map[string]interface{} {
	return map[string]interface{}{"id": item.ID, "type": item.Type, "code": item.Code, "global": item.Global, "category_id": item.Category,
		"language": item.Language, "desc": item.Desc, "num": item.Number, "name": item.Name, "keyword": item.Keyword, "is_show_ad": item.IsShowAd}
}
