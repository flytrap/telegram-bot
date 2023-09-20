package indexsearch

import (
	"context"
)

type IndexSearch interface {
	Init(ctx context.Context) error                                  // 初始化, 创建索引
	AddItem(ctx context.Context, key string, data interface{}) error // 设置词条
	DeleteItem(ctx context.Context, key string) error                // 删除词条
	Delete(ctx context.Context) error                                // 删除词条
	Search(ctx context.Context, text string, category string, page int64, size int64) (int64, []map[string]interface{}, error)
	GetName() string
}

func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}
