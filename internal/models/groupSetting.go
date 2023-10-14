package models

import "gorm.io/gorm"

type GroupSetting struct {
	gorm.Model

	Name string `json:"name" gorm:"size:256;comment:名字"`
	Code string `json:"code" gorm:"unique;size:64;comment:唯一标识"`

	// 配置信息
	NotRobot        bool   `json:"not_robot" gorm:"default:false;comment:机器人校验开启"`
	RobotTimeout    int    `json:"robot_timeout" gorm:"default:30;comment:机器人校验时间"`
	Welcome         bool   `json:"welcome" gorm:"default:true;comment:是否开启新人欢迎词"`
	WelcomeDesc     bool   `json:"welcome_desc" gorm:"default:true;comment:是否显示描述信息"`
	WelcomePinned   bool   `json:"welcome_pinned" gorm:"default:true;comment:是否显示置顶消息"`
	WelcomeKillMe   int    `json:"welcome_kill_me" gorm:"default:30;comment:开启欢迎词自毁时间(s,0不自毁)"`
	WelcomeTemplate string `json:"welcome_template" gorm:"size:256;comment:自定义欢迎词模版"`
}
