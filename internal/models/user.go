package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	// 系统信息
	UserId       int64  `json:"user_id" gorm:"unique;comment:用户id"`
	FirstName    string `json:"first_name" gorm:"size:256"`
	LastName     string `json:"last_name" gorm:"size:256"`
	Username     string `json:"username" gorm:"size:64;comment:用户名"`
	LanguageCode string `json:"language_code" gorm:"size:32;comment:语言代号"`
	IsBot        bool   `json:"is_bot" gorm:"comment:是否是机器人"`
	IsPremium    bool   `json:"is_premium" gorm:"comment:是否是会员"`
	// 配置信息
	Lang string `json:"lang" gorm:"size:32;comment:语言"`
}
