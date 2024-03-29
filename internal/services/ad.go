package services

import (
	"errors"
	"math/rand"
	"time"

	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/flytrap/telegram-bot/internal/repositories"
	"github.com/mitchellh/mapstructure"
)

type AdService interface {
	KeywordAd(keyword string) (*models.Ad, error)
	List(q string, page int64, size int64, ordering string, data interface{}) (n int64, err error)
	Update(id uint, info map[string]interface{}) error
	Create(info map[string]interface{}) error
	Delete(ids []uint) (err error)
}

func NewAdService(repo repositories.AdRepository, categoryService CategoryService) AdService {
	res := AdServiceImp{repo: repo, categoryService: categoryService, keywordMap: map[string]*models.Ad{}, globalList: []*models.Ad{}}
	return &res
}

type AdServiceImp struct {
	repo            repositories.AdRepository
	categoryService CategoryService
	isLoad          bool
	keywordMap      map[string]*models.Ad
	globalList      []*models.Ad
}

func (s *AdServiceImp) List(q string, page int64, size int64, ordering string, data interface{}) (n int64, err error) {
	n, err = s.repo.List(q, time.Time{}, time.Time{}, (page-1)*size, size, ordering, data)
	if err != nil {
		return
	}
	return
}

func (s *AdServiceImp) Update(id uint, info map[string]interface{}) error {
	return s.repo.Update(id, info)
}

func (s *AdServiceImp) Create(info map[string]interface{}) error {
	data := models.Ad{}
	err := mapstructure.Decode(info, &data)
	if err != nil {
		return err
	}
	return s.repo.Create(&data)
}

func (s *AdServiceImp) Delete(ids []uint) (err error) {
	return s.repo.Delete(ids)
}

func (s *AdServiceImp) Load() error {
	if s.isLoad {
		return nil
	}
	results := []models.Ad{}
	_, err := s.repo.List("", time.Now(), time.Time{}, 0, 10000, "", results)
	if err != nil {
		return err
	}
	for _, item := range results {
		if len(item.Keyword) > 0 {
			s.keywordMap[item.Keyword] = &item
		}
		if item.Global != 0 {
			s.globalList = append(s.globalList, &item)
		}
	}
	s.isLoad = true
	return nil
}

func (s *AdServiceImp) KeywordAd(keyword string) (*models.Ad, error) {
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
