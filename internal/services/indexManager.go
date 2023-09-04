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
	LoadData(ctx context.Context, language string) error                          // 更新索引数据
	InitIndex(ctx context.Context) error                                          // 初始化索引数据
	DeleteAllIndex(ctx context.Context)                                           // 清除索引信息
	AddItems(ctx context.Context, name string, data map[string]interface{}) error // 添加词条
	RemoveItem(ctx context.Context, name string, key string) error                // 删除词条
	Query(ctx context.Context, name string, text string, category string, page int64, size int64) ([]map[string]string, error)
}

func NewIndexMangerService(client rueidis.CoreClient, gs GroupService) IndexMangerService {
	i := IndexMangerServiceImp{Client: client, gs: gs, indexes: map[string]indexsearch.IndexSearch{}}
	return &i
}

type IndexMangerServiceImp struct {
	Client  rueidis.CoreClient
	gs      GroupService
	indexes map[string]indexsearch.IndexSearch
}

func (s *IndexMangerServiceImp) indexName(key string) string {
	return fmt.Sprintf("index:%s", key)
}

func (s *IndexMangerServiceImp) addIndex(ctx context.Context, name string, language string, prefix string) error {
	index := indexsearch.NewRedisSearch(&s.Client, name, language, fmt.Sprintf("%s:%s", prefix, language))
	s.indexes[index.GetName()] = index
	return index.Init(ctx)
}

func (s *IndexMangerServiceImp) AddItems(ctx context.Context, name string, data map[string]interface{}) error {
	index, ok := s.indexes[name]
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

func (s *IndexMangerServiceImp) RemoveItem(ctx context.Context, name string, key string) error {
	index, ok := s.indexes[name]
	if !ok {
		return errors.New("index not found")
	}
	return index.DeleteItem(ctx, key)
}

func (s *IndexMangerServiceImp) Query(ctx context.Context, name string, text string, category string, page int64, size int64) ([]map[string]string, error) {
	index, ok := s.indexes[name]
	if !ok {
		return nil, errors.New("index not found")
	}
	return index.Search(ctx, text, category, page, size)
}

func (s *IndexMangerServiceImp) RemoveIndex(ctx context.Context, name string) error {
	index, ok := s.indexes[name]
	if !ok {
		return errors.New("index not found")
	}
	return index.Delete(ctx)
}

// 加载数据
func (s *IndexMangerServiceImp) LoadData(ctx context.Context, language string) error {
	indexName := s.indexName(language)

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
			data, err := s.gs.GetMany("", language, n, 1000)
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
		name := s.indexName(lang)
		err := s.addIndex(ctx, name, lang, config.C.Redis.KeyPrefix)
		if err != nil {
			logrus.Warning(err)
			return err
		}
	}
	return nil
}

func (s *IndexMangerServiceImp) DeleteIndex(ctx context.Context, lang string) error {
	name := s.indexName(lang)
	index := indexsearch.NewRedisSearch(&s.Client, name, lang, fmt.Sprintf("%s:%s", config.C.Redis.KeyPrefix, lang))
	return index.Delete(ctx)
}

func (s *IndexMangerServiceImp) DeleteAllIndex(ctx context.Context) {
	for _, lang := range config.C.Bot.Languages {
		name := s.indexName(lang)
		index := indexsearch.NewRedisSearch(&s.Client, name, lang, fmt.Sprintf("%s:%s", config.C.Redis.KeyPrefix, lang))
		index.Delete(ctx)
		logrus.Info("delete index :" + name)
	}
}
