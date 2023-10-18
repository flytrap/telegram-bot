package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/flytrap/telegram-bot/pkg/human"
	"gorm.io/gorm"
)

type JSON json.RawMessage

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	if len(bytes) == 0 {
		return nil
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSON(result)
	return err
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

type DataInfo struct {
	gorm.Model

	Category uint   `json:"category" mapstructure:"category" gorm:"comment:分类"`
	Language string `json:"language" mapstructure:"language" gorm:"size:32;comment:语言"`
	Name     string `json:"name" mapstructure:"name" gorm:"size:256;comment:名字"`
	Code     string `json:"code" mapstructure:"code" gorm:"unique;size:64;comment:唯一标识"`
	Tid      int64  `json:"tid" mapstructure:"tid" gorm:"comment:TgId"`
	Type     int8   `json:"type" mapstructure:"type" gorm:"default:1;comment:类型,区分group|channel"`
	Number   uint32 `json:"number" mapstructure:"number" gorm:"comment:人数"`
	Desc     string `json:"desc" mapstructure:"desc" gorm:"type:text;comment:描述信息"`
	Weight   int32  `json:"weight" mapstructure:"weight" gorm:"default:0;comment:权重"`
	Private  string `json:"private" mapstructure:"private" gorm:"comment:私密信息"`
	Location string `json:"location" mapstructure:"location" gorm:"size:64;comment:地理位置"`
	Extend   string `json:"extend" mapstructure:"extend" gorm:"type:text;comment:扩展信息"`
	Time     int64  `json:"time" mapstructure:"time" gorm:"comment:更新时间"`
	Images   JSON   `json:"images,omitempty" mapstructure:"images,omitempty" gorm:"type:json;comment:图片"`

	Tags []Tag `json:"tags" gorm:"many2many:data_tag"`
}

func (s *DataInfo) HumanSize() string {
	return human.HumanSize(int64(s.Number))
}
