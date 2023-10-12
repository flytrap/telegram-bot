package indexsearch

import (
	"context"
)

type IndexSearch interface {
	Init(ctx context.Context) error                                          // 初始化, 创建索引
	GetItem(ctx context.Context, key string) (map[string]interface{}, error) // 获取词条信息
	AddItem(ctx context.Context, key string, data interface{}) error         // 设置词条
	DeleteItem(ctx context.Context, key string) error                        // 删除词条
	Delete(ctx context.Context) error                                        // 删除词条
	Search(ctx context.Context, query SearchReq) (int64, []map[string]interface{}, error)
	GetName() string
}

type IndexInfo struct {
	Name     float64 // 标题(string)
	Category float64 // 分类(string, tag)
	Code     float64 // 数据代号(string, tag)
	Type     float64 // 数据类型(int, num)
	Desc     float64 // 详细内容(string)
}

type SearchReq struct {
	Category string // 分类
	Q        string // 搜索字符串
	Page     int64  // 分页
	Size     int64  // 数量
	Order    string // 排序
	Tag      string // 标签搜索
}

func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}
