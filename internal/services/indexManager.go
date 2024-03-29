package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/flytrap/telegram-bot/internal/config"
	"github.com/flytrap/telegram-bot/internal/serializers"
	"github.com/flytrap/telegram-bot/pkg/indexsearch"
	"github.com/jinzhu/copier"
	"github.com/redis/rueidis"
	"github.com/sirupsen/logrus"
)

type IndexMangerService interface {
	IndexName(key string) string
	LoadData(ctx context.Context, indexName string, language string) error                        // 加载索引数据
	InitIndex(ctx context.Context) error                                                          // 初始化索引数据
	DeleteAllIndex(ctx context.Context)                                                           // 清除索引信息
	AddItems(ctx context.Context, indexName string, data map[string]*serializers.DataCache) error // 添加词条
	GetItem(ctx context.Context, indexName string, code string) (map[string]interface{}, error)   // 添加词条
	RemoveItem(ctx context.Context, indexName string, key string) error                           // 删除词条
	Query(ctx context.Context, indexName string, info indexsearch.SearchReq) (int64, []map[string]interface{}, error)
}

func NewIndexMangerService(client rueidis.CoreClient, dataService DataService, categoryService CategoryService) IndexMangerService {
	i := indexMangerServiceImp{Client: client, dataService: dataService, categoryService: categoryService, indexes: map[string]indexsearch.IndexSearch{}}
	return &i
}

type indexMangerServiceImp struct {
	Client          rueidis.CoreClient
	dataService     DataService
	categoryService CategoryService
	indexes         map[string]indexsearch.IndexSearch
}

func (s *indexMangerServiceImp) IndexName(key string) string {
	return fmt.Sprintf("%s:%s", config.C.Index.Name, key)
}

func (s *indexMangerServiceImp) addIndex(ctx context.Context, indexName string, language string, prefix string) error {
	dw := config.C.Index.DescWeight
	if dw != 0 && config.C.Index.Mode == "min" {
		dw = 0 // 小内存模式，没有desc字段
	}
	config := indexsearch.IndexInfo{Name: 1, Category: 1, Code: 1, Type: 1, Desc: config.C.Index.DescWeight, NumberFields: config.C.Index.NumFilter}
	index := indexsearch.NewRedisSearch(&s.Client, indexName, language, fmt.Sprintf("%s:%s", prefix, language), config)
	s.indexes[indexName] = index
	return index.Init(ctx)
}

func (s *indexMangerServiceImp) AddItems(ctx context.Context, indexName string, data map[string]*serializers.DataCache) error {
	index, ok := s.indexes[indexName]
	if !ok {
		return errors.New("index not found")
	}
	var info interface{}
	for key, item := range data {
		switch config.C.Index.Mode {
		case "min":
			info = &serializers.DataCacheLocationMin{}
		case "max":
			info = &serializers.DataCacheLocationMax{}
		default:
			info = &serializers.DataCacheLocation{}
		}
		err := copier.Copy(info, item)
		if err != nil {
			logrus.Warning(err)
			continue
		}
		switch config.C.Index.Mode {
		case "min":
			info.(*serializers.DataCacheLocationMin).ParseLocation(item.Location)
		case "max":
			info.(*serializers.DataCacheLocationMax).ParseLocation(item.Location)
			info.(*serializers.DataCacheLocationMax).ParseImages(item.Images)
		default:
			info.(*serializers.DataCacheLocation).ParseLocation(item.Location)
		}

		err = index.AddItem(ctx, key, &info)
		if err != nil {
			logrus.Warning(err)
		}
	}
	return nil
}

func (s *indexMangerServiceImp) GetItem(ctx context.Context, indexName string, code string) (map[string]interface{}, error) {
	index, ok := s.indexes[indexName]
	if !ok {
		return nil, errors.New("index not found")
	}
	return index.GetItem(ctx, code)
}

func (s *indexMangerServiceImp) RemoveItem(ctx context.Context, indexName string, key string) error {
	index, ok := s.indexes[indexName]
	if !ok {
		return errors.New("index not found")
	}
	return index.DeleteItem(ctx, key)
}

func (s *indexMangerServiceImp) Query(ctx context.Context, indexName string, info indexsearch.SearchReq) (int64, []map[string]interface{}, error) {
	index, ok := s.indexes[indexName]
	if !ok {
		return 0, nil, errors.New("index not found")
	}
	return index.Search(ctx, info)
}

func (s *indexMangerServiceImp) RemoveIndex(ctx context.Context, indexName string) error {
	index, ok := s.indexes[indexName]
	if !ok {
		return errors.New("index not found")
	}
	return index.Delete(ctx)
}

// 加载数据
func (s *indexMangerServiceImp) LoadData(ctx context.Context, indexName string, language string) error {
	var (
		wg sync.WaitGroup
		c  = make(chan map[string]*serializers.DataCache)
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
		size := int64(5000)
		if config.C.Index.Mode == "min" {
			size = 1000
		} else if config.C.Index.Mode == "max" {
			size = 10000
		}
		for {
			data := []*serializers.DataCache{}
			_, err := s.dataService.List("", "", language, n, size, "", &data)
			if err != nil {
				logrus.Warning(err)
				break
			}
			if len(data) == 0 {
				logrus.Info("query ok")
				break
			}
			items := map[string]*serializers.DataCache{}
			for _, item := range data {
				id, err := strconv.Atoi(item.Category)
				if err == nil && id > 0 {
					ca, _ := s.categoryService.GetName(uint(id))
					item.Category = ca
				}
				item.Weight += float32(item.Number) / 10000
				items[item.Code] = item
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
func (s *indexMangerServiceImp) InitIndex(ctx context.Context) error {
	logrus.Debug("init index", config.C.Index.Language)
	name := s.IndexName(config.C.Index.Language)
	return s.addIndex(ctx, name, config.C.Index.Language, config.C.Redis.KeyPrefix)
}

func (s *indexMangerServiceImp) DeleteIndex(ctx context.Context, lang string) error {
	name := s.IndexName(lang)
	err := s.addIndex(ctx, name, lang, config.C.Redis.KeyPrefix)
	if err != nil {
		return err
	}
	index := s.indexes[name]
	delete(s.indexes, name)
	return index.Delete(ctx)
}

func (s *indexMangerServiceImp) DeleteAllIndex(ctx context.Context) {
	s.DeleteIndex(ctx, config.C.Index.Language)
	logrus.Info("delete index: " + config.C.Index.Language)
}
