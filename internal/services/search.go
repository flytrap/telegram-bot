package services

import (
	"context"
	"fmt"

	"github.com/flytrap/telegram-bot/internal/config"
	"github.com/flytrap/telegram-bot/internal/serializers"
	"github.com/flytrap/telegram-bot/pkg/human"
)

type SearchService interface {
	QueryItems(context.Context, string, int64, int64) ([]string, bool, error) // 搜素信息
}

func NewSearchService(dataService DataService, im IndexMangerService, adService AdService) SearchService {
	return &searchServiceImp{dataService: dataService, IndexManager: im, adService: adService}
}

type searchServiceImp struct {
	dataService  DataService
	IndexManager IndexMangerService
	adService    AdService
}

func (s *searchServiceImp) QueryItems(ctx context.Context, text string, page int64, size int64) (items []string, hasNext bool, err error) {
	if page >= config.C.Bot.MaxPage {
		page = config.C.Bot.MaxPage // 阻止过多翻页
	}
	var n int64
	if config.C.Bot.UseCache {
		n, items, err = s.QueryCacheItems(ctx, "chinese", text, "", page, config.C.Bot.PageSize)
	} else {
		n, items, err = s.QueryDbItems(text, page, config.C.Bot.PageSize)
	}
	if err != nil {
		return nil, false, err
	}
	hasNext = (page-1)*size+int64(len(items)) < n
	ad := s.loadAd(text)
	if len(ad) > 0 {
		items = append([]string{ad, ""}, items...) // 增加广告
	}
	return items, hasNext, nil
}

func (s *searchServiceImp) loadAd(keyword string) string {
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

func (s *searchServiceImp) QueryDbItems(text string, page int64, size int64) (int64, []string, error) {
	data := []*serializers.DataSerializer{}
	n, err := s.dataService.SearchTag(text, page, 15, &data)
	if err != nil {
		return 0, nil, err
	}
	items := []string{}
	for i, item := range data {
		items = append(items, item.ItemInfo(i+1))
	}
	return n, items, nil
}

func (s *searchServiceImp) QueryCacheItems(ctx context.Context, name string, text string, category string, page int64, size int64) (int64, []string, error) {
	index := s.IndexManager.IndexName(name)
	n, data, err := s.IndexManager.Query(ctx, index, text, category, page, size)
	if err != nil {
		return 0, nil, err
	}
	items := []string{}
	for i, item := range data {
		items = append(items, human.TgGroupItemInfo(int(page-1)*int(size)+i+1, item["code"].(string), int(item["type"].(float64)), item["name"].(string), int64(item["number"].(float64))))
	}
	return n, items, nil
}
