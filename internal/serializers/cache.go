package serializers

import "strings"

type DataCache struct {
	Category string   `json:"category"` // 分类
	Language string   `json:"language"` // 语言
	Name     string   `json:"name"`     // 名称
	Code     string   `json:"code" `    // 代号，唯一标识
	Type     int8     `json:"type"`     // 数据类型
	Number   uint32   `json:"number"`   // 数量
	Desc     string   `json:"desc"`     // 详细信息
	Weight   float32  `json:"weight"`   // 数据权重
	Location string   `json:"location"` // 地理位置
	Area     string   `json:"area"`     // 区
	Extend   string   `json:"extend"`   // 拓展数据
	Images   []string `json:"images"`   // 图片
}

type DataCacheLocation struct {
	Category string   `json:"category"` // 分类
	Language string   `json:"language"` // 语言
	Name     string   `json:"name"`     // 名称
	Code     string   `json:"code" `    // 代号，唯一标识
	Type     int8     `json:"type"`     // 数据类型
	Number   uint32   `json:"number"`   // 数量
	Desc     string   `json:"desc"`     // 详细信息
	Weight   float32  `json:"weight"`   // 数据权重
	Time     float32  `json:"time"`     // 更新时间
	Private  string   `json:"private"`  // 私密数据
	Province string   `json:"province"` // 省
	City     string   `json:"city"`     // 市
	Area     string   `json:"area"`     // 区
	Extend   string   `json:"extend"`   // 拓展数据
	Images   []string `json:"images"`   // 图片
}

func (s *DataCacheLocation) ParseLocation(location string) {
	li := strings.Split(location, "-")
	if len(li) > 0 {
		s.Province = li[0]
	}
	if len(li) > 1 {
		s.City = li[1]
	}
	if len(li) > 2 {
		s.Area = li[2]
	}
}
