package serializers

import (
	"encoding/json"
	"strings"

	"github.com/yanyiwu/gojieba"
)

var x = gojieba.NewJieba()

type DataCache struct {
	Category string  `json:"category"` // 分类
	Language string  `json:"language"` // 语言
	Name     string  `json:"name"`     // 名称
	Code     string  `json:"code" `    // 代号，唯一标识
	Type     int8    `json:"type"`     // 数据类型
	Number   uint32  `json:"number"`   // 数量
	Desc     string  `json:"desc"`     // 详细信息
	Weight   float32 `json:"weight"`   // 数据权重
	Location string  `json:"location"` // 地理位置
	Extend   string  `json:"extend"`   // 拓展数据
	Time     int64   `json:"time"`     // 更新时间
	Private  string  `json:"private"`  // 私密数据
	Images   []byte  `json:"images"`   // 图片
}

// 小内存模式
type DataCacheLocationMin struct {
	Category string   `json:"category"` // 分类
	Name     string   `json:"name"`     // 名称
	Code     string   `json:"code" `    // 代号，唯一标识
	Number   uint32   `json:"number"`   // 数量
	Weight   float32  `json:"weight"`   // 数据权重
	Time     int64    `json:"time"`     // 更新时间
	Tags     []string `json:"tags"`     // 标签
}

type DataCacheLocation struct {
	DataCacheLocationMin
	Language string `json:"language"` // 语言
	Type     int8   `json:"type"`     // 数据类型
	Desc     string `json:"desc"`     // 详细信息
}

// 大内存模式
type DataCacheLocationMax struct {
	DataCacheLocation

	Location string   `json:"location"` // 地理位置
	Private  string   `json:"private"`  // 私密数据
	Extend   string   `json:"extend"`   // 拓展数据
	Images   []string `json:"images"`   // 图片
}

func (s *DataCacheLocationMin) ParseLocation(location string) {
	li := strings.Split(location, "-")
	if len(s.Tags) == 0 {
		s.Tags = []string{s.Category}
	}
	if len(li) > 0 {
		s.Tags = append(s.Tags, x.CutAll(strings.TrimSpace(li[0]))...)
	}
	if len(li) > 1 {
		s.Tags = append(s.Tags, x.CutAll(strings.TrimSpace(li[1]))...)
	}
	if len(li) > 2 {
		s.Tags = append(s.Tags, x.CutAll(strings.TrimSpace(li[2]))...)
	}
}

func (s *DataCacheLocationMax) ParseImages(images []byte) error {
	if len(s.Images) > 0 {
		return json.Unmarshal(images, &s.Images)
	}
	return nil
}
