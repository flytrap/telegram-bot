package models

import (
	"time"

	"gorm.io/gorm"
)

// 广告
type Ad struct {
	gorm.Model

	Keyword  string    `json:"keyword" gorm:"comment:关键词"`
	IsGlobal bool      `json:"is_global" gorm:"default:false;comment:是否是全局词条"`
	Expire   time.Time `json:"expire" gorm:"comment:过期时间"`
	Category uint      `json:"category" gorm:"comment:分类"`
	Name     string    `json:"name" gorm:"index;size:256;comment:名字"`
	Code     string    `json:"code" gorm:"unique;size:64;comment:唯一标识"`
	Type     int8      `json:"type" gorm:"default:1;comment:类型,区分group|channel"`
	Number   uint32    `json:"number" gorm:"comment:人数"`
	Desc     string    `json:"desc" gorm:"type:text;comment:描述信息"`
}
