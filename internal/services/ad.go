package services

import (
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/flytrap/telegram-bot/internal/repositories"
	"github.com/mitchellh/mapstructure"
)

type AdService interface {
	KeywordAd(keyword string) (map[string]interface{}, error)
	List(q string, page int64, size int64) (data []map[string]interface{}, err error)
	Update(id uint, info map[string]interface{}) error
}

func NewAdService(repo repositories.AdRepository) AdService {
	res := AdServiceImp{repo: repo, keywordMap: map[string]map[string]interface{}{}, globalList: []map[string]interface{}{}}
	return &res
}

type AdServiceImp struct {
	repo            repositories.AdRepository
	categoryService CategoryService
	isLoad          bool
	keywordMap      map[string]map[string]interface{}
	globalList      []map[string]interface{}
}

func (s *AdServiceImp) List(q string, page int64, size int64) (data []map[string]interface{}, err error) {
	result, err := s.repo.GetMany(q, time.Time{}, time.Time{}, (page-1)*size, size)
	if err != nil {
		return nil, err
	}
	json.Marshal(result)
	for _, item := range result {
		c, _ := s.categoryService.GetName(item.Category)
		info := item.ToMap()
		info["category"] = c
		data = append(data, info)
	}
	return
}

func (s *AdServiceImp) Update(id uint, info map[string]interface{}) error {
	t := &models.User{}
	err := mapstructure.Decode(info, t)
	if err != nil {
		return err
	}

	return s.repo.Update(id, info)
}

func (s *AdServiceImp) Load() error {
	if s.isLoad {
		return nil
	}
	results, err := s.repo.GetMany("", time.Now(), time.Time{}, 0, 10000)
	if err != nil {
		return err
	}
	for _, item := range results {
		info := item.ToMap()
		if len(item.Keyword) > 0 {
			s.keywordMap[item.Keyword] = info
		}
		if item.Global != 0 {
			s.globalList = append(s.globalList, info)
		}
	}
	s.isLoad = true
	return nil
}

func (s *AdServiceImp) KeywordAd(keyword string) (map[string]interface{}, error) {
	item, ok := s.keywordMap[keyword]
	if ok {
		return item, nil
	}
	if len(s.globalList) == 0 {
		return nil, errors.New("keyword not found")
	}
	i := rand.Intn(len(s.globalList))
	return s.globalList[i], nil
}
