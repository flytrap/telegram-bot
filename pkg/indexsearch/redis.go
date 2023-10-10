package indexsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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

func NewRedisSearch(client *rueidis.CoreClient, index string, language string, prefix string, config IndexInfo) IndexSearch {
	return &IndexSearchOnRedis{client: *client, Name: index, Prefix: prefix, Language: language, Score: 1, IndexInfo: config}
}

type IndexSearchOnRedis struct {
	client    rueidis.CoreClient
	Name      string // 索引名称
	Language  string
	Prefix    string
	Score     int64
	IndexInfo IndexInfo
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
	cmd := s.client.B().FtCreate().Index(s.Name).OnJson().Prefix(1).Prefix(s.Prefix).Language(s.Language).Nohl()
	tagBuild := cmd.Schema().FieldName("$name").As("name").Text().Weight(s.IndexInfo.Name).FieldName("$category").As("category").Tag()
	for _, tag := range s.IndexInfo.Tags {
		tagBuild = tagBuild.FieldName(fmt.Sprintf("$%s", tag)).As(tag).Tag()
	}
	build := tagBuild.FieldName("$type").As("type").Numeric().FieldName("$number").As("number").Numeric().Sortable().FieldName("$time").As("time").Numeric().Sortable().FieldName("$weight").As("weight").Numeric().Sortable().FieldName("$desc").As("desc").Text().Weight(s.IndexInfo.Desc).Build()

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
func (s *IndexSearchOnRedis) GetItem(ctx context.Context, key string) (map[string]string, error) {
	cmd := s.client.B().JsonGet().Key(s.PrefixKey(key)).Path("$").Build()
	return s.client.Do(ctx, cmd).AsStrMap()
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
func (s *IndexSearchOnRedis) Search(ctx context.Context, query SearchReq) (int64, []map[string]interface{}, error) {
	q := ""
	if len(query.Tags) > 0 {
		li := []string{}
		for k, v := range query.Tags {
			li = append(li, fmt.Sprintf("@%s:{%s}", k, v))
		}
		q = strings.Join(li, "|")
	} else {
		text := filterQuery(query.Q) // 过滤特殊字符
		q = fmt.Sprintf("(@category:{%s})|%s", text, text)
	}
	if len(query.Category) > 0 {
		q = fmt.Sprintf("@category:{%s} %s", query.Category, q)
	}
	n, resp, err := s.Query(ctx, q, query.Order, (query.Page-1)*query.Size, query.Size)
	if err != nil {
		return 0, nil, err
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

	return n, results, nil
}

func (s *IndexSearchOnRedis) Query(ctx context.Context, query string, order string, offset int64, num int64) (int64, []rueidis.FtSearchDoc, error) {
	cmd := s.client.B().FtSearch().Index(s.Name).Query(query)
	var build rueidis.Completed
	if len(order) > 0 {
		build = cmd.Sortby(order).Desc().Limit().OffsetNum(offset, num).Build()
	} else {
		build = cmd.Limit().OffsetNum(offset, num).Build()
	}
	n, resp, err := s.client.Do(ctx, build).AsFtSearch()
	logrus.Info("query: ", query, ";result: ", n)
	if err != nil {
		return 0, nil, err
	}
	return n, resp, nil
}

func filterQuery(q string) string {
	for _, key := range []string{"(", ")", "[", "]", "$", "@", "."} {
		q = strings.ReplaceAll(q, key, "|")
	}
	return strings.ReplaceAll(q, "||", "|")
}
