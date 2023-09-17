package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model

	Name   string `json:"name" gorm:"unique;size:256;comment:名字"`
	Weight int32  `json:"weight" gorm:"default:0;comment:权重"`
}
