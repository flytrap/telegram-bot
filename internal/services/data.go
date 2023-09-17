package services

import (
	"errors"

	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/flytrap/telegram-bot/internal/repositories"
	"github.com/flytrap/telegram-bot/internal/serializers"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

func NewDataService(repo repositories.DataInfoRepository, tagService DataTagService, categoryService CategoryService) DataService {
	return &DataInfoServiceImp{repo: repo, tagService: tagService, categoryService: categoryService}
}

type DataService interface {
	List(q string, category string, language string, page int64, size int64, ordering string) (int64, []map[string]interface{}, error)
	SearchTag(tag string, page int64, size int64) (data []*serializers.DataSerializer, err error)
	GetNeedUpdateCode(days int, page int64, size int64) ([]string, error)
	Update(code string, tid int64, name string, desc string, num uint32, weight int) error
	Delete(codes []string) (err error)

	UpdateOrCreate(code string, tid int64, name string, desc string, num uint32, tags []string, category string, lang string) error
}

type DataInfoServiceImp struct {
	repo            repositories.DataInfoRepository
	tagService      DataTagService
	categoryService CategoryService
}

func (s *DataInfoServiceImp) Exists(code string) bool {
	data, _ := s.repo.Get(code)
	return data != nil
}

func (s *DataInfoServiceImp) List(q string, category string, language string, page int64, size int64, ordering string) (int64, []map[string]interface{}, error) {
	cs := uint(0)
	var err error
	if len(category) > 0 {
		cs, err = s.categoryService.GetId(category)
		if err != nil {
			return 0, nil, errors.New("category not found: " + category)
		}
	}
	n, data, err := s.repo.List(q, cs, language, (page-1)*size, size, ordering)
	if err != nil {
		return 0, nil, err
	}
	results := []map[string]interface{}{}
	for _, item := range data {
		c, _ := s.categoryService.GetName(item.Category)
		o := map[string]interface{}{}
		mapstructure.Decode(item, &o)
		o["Category"] = c

		results = append(results, o)
	}
	return n, results, nil
}

func (s *DataInfoServiceImp) SearchTag(tag string, page int64, size int64) (data []*serializers.DataSerializer, err error) {
	res, err := s.repo.QueryTag(tag, (page-1)*size, size)
	data = []*serializers.DataSerializer{}
	for _, item := range res {
		data = append(data, &serializers.DataSerializer{Code: item.Code, Name: item.Name, Number: int64(item.Number)})
	}
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *DataInfoServiceImp) GetNeedUpdateCode(days int, page int64, size int64) ([]string, error) {
	res, err := s.repo.GetNeedUpdate(days, (page-1)*size, size)
	if err != nil {
		return nil, err
	}
	results := []string{}
	for _, item := range res {
		results = append(results, item.Code)
	}
	return results, nil
}

func (s *DataInfoServiceImp) Update(code string, tid int64, name string, desc string, num uint32, weight int) error {
	var params = map[string]interface{}{"name": name}
	if tid != 0 {
		params["tid"] = tid
	}
	if num != 0 {
		params["number"] = num
	}
	if weight > 0 {
		params["weight"] = weight
	}
	if len(desc) > 0 {
		params["desc"] = desc
	}
	return s.repo.Update(code, params)
}

func (s *DataInfoServiceImp) Delete(codes []string) (err error) {
	return s.repo.Delete(codes)
}

func (s *DataInfoServiceImp) UpdateOrCreate(code string, tid int64, name string, desc string, num uint32, tags []string, category string, lang string) error {
	if s.Exists(code) {
		return s.Update(code, tid, name, desc, num, -1)
	}
	c, err := s.categoryService.GetOrCreate(category, 0)
	if err != nil {
		return err
	}
	ts := []models.Tag{}
	for _, t := range tags {
		tag, err := s.tagService.GetOrCreate(t, 0)
		if err != nil {
			logrus.Warn(err)
			continue
		}
		ts = append(ts, *tag)
	}
	data := models.DataInfo{Code: code, Tid: tid, Name: name, Desc: desc, Number: num, Tags: ts, Category: c.ID, Language: lang}
	err = s.repo.Create(&data)
	return err
}
