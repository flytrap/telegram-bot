package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model

	Name   string `json:"name" gorm:"index;size:256;comment:名字"`
	EnName string `json:"en_name" gorm:"index;size:256;comment:英文名字"`
	Weight int    `json:"weight" gorm:"default:0;comment:权重"`
}
