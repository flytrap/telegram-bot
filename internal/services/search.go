package services

import (
	"context"
	"fmt"

	"github.com/flytrap/telegram-bot/internal/config"
	"github.com/flytrap/telegram-bot/internal/serializers"
	"github.com/flytrap/telegram-bot/pkg/human"
	"github.com/flytrap/telegram-bot/pkg/indexsearch"
)

type SearchService interface {
	QueryItems(context.Context, string, string, string, int64, int64) ([]map[string]interface{}, bool, error) // 搜素信息
	GetPrivate(context.Context, string) (string, error)                                                       // 获取隐私信息
	GetDetail(ctx context.Context, code string) (map[string]interface{}, error)                               // 获取详细词条数据
	LoadAd(keyword string) string                                                                             // 获取广告
}

func NewSearchService(dataService DataService, im IndexMangerService, adService AdService) SearchService {
	return &searchServiceImp{dataService: dataService, IndexManager: im, adService: adService}
}

type searchServiceImp struct {
	dataService  DataService
	IndexManager IndexMangerService
	adService    AdService
}

func (s *searchServiceImp) QueryItems(ctx context.Context, category string, tag string, q string, page int64, size int64) (items []map[string]interface{}, hasNext bool, err error) {
	if page >= config.C.Index.MaxPage {
		page = config.C.Index.MaxPage // 阻止过多翻页
	}
	var n int64
	if config.C.Bot.UseCache {
		n, items, err = s.QueryCacheItems(ctx, "chinese", category, tag, q, page, size)
	} else {
		n, items, err = s.QueryDbItems(q, page, size)
	}
	if err != nil {
		return nil, false, err
	}
	hasNext = (page-1)*size+int64(len(items)) < n
	return items, hasNext, nil
}

func (s *searchServiceImp) LoadAd(keyword string) string {
	item, err := s.adService.KeywordAd(keyword)
	if err != nil {
		return ""
	}
	showAd := ""
	if len(item.AdTag) > 0 {
		showAd = fmt.Sprintf("[%s] ", item.AdTag)
	}
	return fmt.Sprintf("%s[%s](%s)", showAd, item.Name, human.TgGroupUrl(item.Code))
}

func (s *searchServiceImp) QueryDbItems(text string, page int64, size int64) (int64, []map[string]interface{}, error) {
	data := []*serializers.DataSerializer{}
	n, err := s.dataService.SearchTag(text, page, 15, &data)
	if err != nil {
		return 0, nil, err
	}
	items := []map[string]interface{}{}
	for _, item := range data {
		info, err := human.Decode(item)
		if err != nil {
			continue
		}
		items = append(items, info)
	}
	return n, items, nil
}

func (s *searchServiceImp) QueryCacheItems(ctx context.Context, name string, category string, tag string, q string, page int64, size int64) (int64, []map[string]interface{}, error) {
	index := s.IndexManager.IndexName(name)
	query := indexsearch.SearchReq{Q: q, Category: category, Page: page, Size: size, Tag: tag} // 查询条件

	return s.IndexManager.Query(ctx, index, query)
}

func (s *searchServiceImp) GetPrivate(ctx context.Context, code string) (string, error) {
	index := s.IndexManager.IndexName("chinese")
	res, err := s.IndexManager.GetItem(ctx, index, code)
	if _, ok := res["private"]; err != nil && ok {
		return "", err
	}
	return fmt.Sprintf("%s\n\n%s", res["name"].(string), res["private"].(string)), nil
}

func (s *searchServiceImp) GetDetail(ctx context.Context, code string) (map[string]interface{}, error) {
	index := s.IndexManager.IndexName("chinese")
	return s.IndexManager.GetItem(ctx, index, code)
}
