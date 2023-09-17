package services

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/flytrap/telegram-bot/internal/config"
	"github.com/flytrap/telegram-bot/pkg/indexsearch"
	"github.com/redis/rueidis"
	"github.com/sirupsen/logrus"
)

type IndexMangerService interface {
	IndexName(key string) string
	LoadData(ctx context.Context, indexName string, language string) error             // 加载索引数据
	InitIndex(ctx context.Context) error                                               // 初始化索引数据
	DeleteAllIndex(ctx context.Context)                                                // 清除索引信息
	AddItems(ctx context.Context, indexName string, data map[string]interface{}) error // 添加词条
	RemoveItem(ctx context.Context, indexName string, key string) error                // 删除词条
	Query(ctx context.Context, indexName string, text string, category string, page int64, size int64) ([]map[string]interface{}, error)
}

func NewIndexMangerService(client rueidis.CoreClient, dataService DataService) IndexMangerService {
	i := IndexMangerServiceImp{Client: client, dataService: dataService, indexes: map[string]indexsearch.IndexSearch{}}
	return &i
}

type IndexMangerServiceImp struct {
	Client      rueidis.CoreClient
	dataService DataService
	indexes     map[string]indexsearch.IndexSearch
}

func (s *IndexMangerServiceImp) IndexName(key string) string {
	return fmt.Sprintf("index:%s", key)
}

func (s *IndexMangerServiceImp) addIndex(ctx context.Context, indexName string, language string, prefix string) error {
	index := indexsearch.NewRedisSearch(&s.Client, indexName, language, fmt.Sprintf("%s:%s", prefix, language))
	s.indexes[indexName] = index
	return index.Init(ctx)
}

func (s *IndexMangerServiceImp) AddItems(ctx context.Context, indexName string, data map[string]interface{}) error {
	index, ok := s.indexes[indexName]
	if !ok {
		return errors.New("index not found")
	}
	for key, item := range data {
		err := index.AddItem(ctx, key, item)
		if err != nil {
			logrus.Warning(err)
		}
	}
	return nil
}

func (s *IndexMangerServiceImp) RemoveItem(ctx context.Context, indexName string, key string) error {
	index, ok := s.indexes[indexName]
	if !ok {
		return errors.New("index not found")
	}
	return index.DeleteItem(ctx, key)
}

func (s *IndexMangerServiceImp) Query(ctx context.Context, indexName string, text string, category string, page int64, size int64) ([]map[string]interface{}, error) {
	index, ok := s.indexes[indexName]
	if !ok {
		return nil, errors.New("index not found")
	}
	return index.Search(ctx, text, category, page, size)
}

func (s *IndexMangerServiceImp) RemoveIndex(ctx context.Context, indexName string) error {
	index, ok := s.indexes[indexName]
	if !ok {
		return errors.New("index not found")
	}
	return index.Delete(ctx)
}

// 加载数据
func (s *IndexMangerServiceImp) LoadData(ctx context.Context, indexName string, language string) error {
	var (
		wg sync.WaitGroup
		c  = make(chan map[string]interface{})
	)

	wg.Add(1)
	// 写redis
	go func() {
		defer wg.Done()
		for v := range c {
			logrus.Info("write items: ", len(v))
			s.AddItems(ctx, indexName, v)
		}
	}()

	wg.Add(1)
	// 读数据库
	go func() {
		defer wg.Done()
		n := int64(1)
		for {
			_, data, err := s.dataService.List("", "", language, n, 1000, "")
			if err != nil {
				logrus.Warning(err)
				break
			}
			if len(data) == 0 {
				logrus.Info("query ok")
				break
			}
			items := map[string]interface{}{}
			for _, item := range data {
				items[item["code"].(string)] = item
			}
			logrus.Info("read items: ", len(items))
			c <- items
			n += 1
		}
		close(c)
	}()
	wg.Wait()
	logrus.Info("db init ok")
	return nil
}

// 初始化索引
func (s *IndexMangerServiceImp) InitIndex(ctx context.Context) error {
	logrus.Debug("init index", config.C.Bot.Languages)
	for _, lang := range config.C.Bot.Languages {
		name := s.IndexName(lang)
		err := s.addIndex(ctx, name, lang, config.C.Redis.KeyPrefix)
		if err != nil {
			logrus.Warning(err)
			return err
		}
	}
	return nil
}

func (s *IndexMangerServiceImp) DeleteIndex(ctx context.Context, lang string) error {
	name := s.IndexName(lang)
	index := indexsearch.NewRedisSearch(&s.Client, name, lang, fmt.Sprintf("%s:%s", config.C.Redis.KeyPrefix, lang))
	return index.Delete(ctx)
}

func (s *IndexMangerServiceImp) DeleteAllIndex(ctx context.Context) {
	for _, lang := range config.C.Bot.Languages {
		s.DeleteIndex(ctx, lang)
		logrus.Info("delete index :" + lang)
	}
}
