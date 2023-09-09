package services

import (
	"encoding/json"
	"time"

	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/flytrap/telegram-bot/internal/repositories"
	"github.com/mitchellh/mapstructure"
)

type AdService interface {
	List(q string, page int64, size int64) (data []map[string]interface{}, err error)
	Update(id uint, info map[string]interface{}) error
}

func NewAdService(repo repositories.AdRepository) AdService {
	return &AdServiceImp{repo: repo}
}

type AdServiceImp struct {
	repo            repositories.AdRepository
	categoryService CategoryService
	codeMap         map[string]map[string]interface{}
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
	results, err := s.repo.GetMany("", time.Now(), time.Time{}, 0, 10000)
	if err != nil {
		return err
	}
	for _, item := range results {
		info := item.ToMap()
		s.codeMap[item.Code] = info
	}
	return nil
}
