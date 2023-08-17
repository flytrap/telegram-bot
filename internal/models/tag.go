package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model

	Name   string `json:"name" gorm:"unique;size:256;comment:名字"`
	EnName string `json:"en_name" gorm:"index;size:256;comment:英文名字"`
	Weight int32  `json:"weight" gorm:"default:0;comment:权重"`

	Groups []*Group `json:"groups" gorm:"many2many:group_tag"`
}
