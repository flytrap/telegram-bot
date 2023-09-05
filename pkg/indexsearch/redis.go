package indexsearch

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/rueidis"
	"github.com/sirupsen/logrus"
)

var LanguageMap = map[string]string{
	"english":    "English",
	"chinese":    "简体中文",
	"arabic":     "عربي",      // 阿拉伯语
	"danish":     "dansk",     // 丹麦语
	"french":     "français",  // 法语
	"german":     "Deutsch",   // 德语
	"hungarian":  "Magyar",    // 匈牙利, 暂时用不到
	"italian":    "italiano",  // 意大利
	"norwegian":  "norsk",     // 挪威
	"portuguese": "português", // 葡萄牙
	"romanian":   "Română",    // 罗马尼亚
	"russian":    "русский",   //俄语
	"serbian":    "Српски",    // 塞尔维亚
	"spanish":    "español",   // 西班牙
	"swedish":    "svenska",   // 瑞典
	"tamil":      "தமிழ்",     // 泰米尔语
	"turkish":    "Türkçe",    // 土耳其语
	"yiddish":    "ייִדיש",    // 意第绪语
}

func NewRedisSearch(client *rueidis.CoreClient, index string, language string, prefix string) IndexSearch {
	return &IndexSearchOnRedis{client: *client, Name: index, Prefix: prefix, Language: language, Score: 1, ScoreField: "weight",
		IndexInfo: IndexInfo{Name: 1, Category: 1, Code: 1, Type: 1, Body: 0.5}}
}

type IndexField struct {
	Name   string
	Type   string   // geo, text, num, tag, vector
	Weight float64  // 权重
	Algo   string   // vector 使用
	Args   []string // vector使用
}

type IndexInfo struct {
	Name     float64 // 标题(string)
	Category float64 // 分类(string, tag)
	Code     float64 // 数据代号(string, tag)
	Type     float64 // 数据类型(int, num)
	Body     float64 // 详细内容(string)
}

type IndexSearchOnRedis struct {
	client     rueidis.CoreClient
	Name       string // 索引名称
	Language   string
	Prefix     string
	Score      int64
	ScoreField string
	IndexInfo  IndexInfo
}

// 初始化索引信息
func (s *IndexSearchOnRedis) Init(ctx context.Context) error {
	cmd1 := s.client.B().FtList().Build()
	vs, err := s.client.Do(ctx, cmd1).AsStrSlice()
	if err != nil {
		return err
	}
	logrus.Info(vs)
	if IsContain(vs, s.Name) {
		return nil
	}
	// cmd := s.client.B().FtCreate().Index(s.Name).OnJson().Prefix(1).Prefix(s.Prefix).Language(s.Language).Score(float64(s.Score)).ScoreField(s.ScoreField).Nohl()
	cmd := s.client.B().FtCreate().Index(s.Name).OnJson().Prefix(1).Prefix(s.Prefix).Language(s.Language).Nohl()
	build := cmd.Schema().FieldName("$name").As("name").Text().Weight(s.IndexInfo.Name).FieldName("$category").As("category").Tag().FieldName("$code").As("code").Tag().FieldName("$type").As("type").Numeric().FieldName("$body").As("body").Text().Weight(s.IndexInfo.Body).Build()

	err = s.client.Do(ctx, build).Error()
	return err
}

func (s *IndexSearchOnRedis) GetName() string {
	return s.Name
}

func (s *IndexSearchOnRedis) PrefixKey(key string) string {
	return fmt.Sprintf("%s:%s", s.Prefix, key)
}

// 添加词条
func (s *IndexSearchOnRedis) AddItem(ctx context.Context, key string, data interface{}) error {
	cmd := s.client.B().JsonSet().Key(s.PrefixKey(key)).Path("$").Value(rueidis.JSON(data)).Build()
	return s.client.Do(ctx, cmd).Error()
}

// 删除词条
func (s *IndexSearchOnRedis) DeleteItem(ctx context.Context, key string) error {
	cmd := s.client.B().JsonDel().Key(s.PrefixKey(key)).Build()
	return s.client.Do(ctx, cmd).Error()
}

// 删除索引
func (s *IndexSearchOnRedis) Delete(ctx context.Context) error {
	cmd := s.client.B().FtDropindex().Index(s.Name).Build()
	return s.client.Do(ctx, cmd).Error()
}

// 搜索
func (s *IndexSearchOnRedis) Search(ctx context.Context, text string, category string, page int64, size int64) ([]map[string]interface{}, error) {
	q := text
	if len(category) > 0 {
		q = fmt.Sprintf("@category:%s %s", category, q)
	}
	resp, err := s.Query(ctx, q, page*size, size)
	if err != nil {
		return nil, err
	}
	results := []map[string]interface{}{}
	for _, item := range resp {
		res := map[string]interface{}{}
		err = json.Unmarshal([]byte(item.Doc["$"]), &res)
		if err != nil {
			continue
		}
		results = append(results, res)
	}

	return results, nil
}

func (s *IndexSearchOnRedis) Query(ctx context.Context, query string, offset int64, num int64) ([]rueidis.FtSearchDoc, error) {
	cmd := s.client.B().FtSearch().Index(s.Name).Query(query).Limit().OffsetNum(offset, num).Build()
	n, resp, err := s.client.Do(ctx, cmd).AsFtSearch()
	logrus.Info(n)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
