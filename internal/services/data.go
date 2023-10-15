package services

import (
	"encoding/json"
	"errors"

	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/flytrap/telegram-bot/internal/repositories"
	"github.com/flytrap/telegram-bot/pkg/human"
	"github.com/sirupsen/logrus"
)

type DataService interface {
	Get(code string) (map[string]interface{}, error)
	List(q string, category string, language string, page int64, size int64, ordering string, data interface{}) (int64, error)
	SearchTag(tag string, page int64, size int64, data interface{}) (total int64, err error)
	GetNeedUpdateCode(days int, page int64, size int64) ([]string, error)
	Update(code string, data map[string]interface{}) error
	Delete(codes []string) (err error)

	UpdateOrCreate(code string, data map[string]interface{}) error
}

func NewDataService(repo repositories.DataInfoRepository, tagService DataTagService, categoryService CategoryService) DataService {
	return &DataInfoServiceImp{repo: repo, tagService: tagService, categoryService: categoryService}
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

func (s *DataInfoServiceImp) Get(code string) (map[string]interface{}, error) {
	res, err := s.repo.Get(code)
	if err != nil {
		return nil, err
	}
	return human.Decode(res)
}

func (s *DataInfoServiceImp) List(q string, category string, language string, page int64, size int64, ordering string, data interface{}) (int64, error) {
	cs := uint(0)
	var err error
	if len(category) > 0 {
		cs, err = s.categoryService.GetId(category)
		if err != nil {
			return 0, errors.New("category not found: " + category)
		}
	}
	n, err := s.repo.List(q, cs, language, (page-1)*size, size, ordering, data)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (s *DataInfoServiceImp) SearchTag(tag string, page int64, size int64, data interface{}) (total int64, err error) {
	total, err = s.repo.QueryTag(tag, (page-1)*size, size, data)
	return total, err
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

func (s *DataInfoServiceImp) Update(code string, params map[string]interface{}) error {
	if v, ok := params["tid"]; !ok || v.(int64) == 0 {
		delete(params, "tid")
	}
	if v, ok := params["number"]; !ok || v.(int64) == 0 {
		delete(params, "number")
	}
	if v, ok := params["weight"]; !ok || v.(int64) <= 0 {
		delete(params, "weight")
	}
	if v, ok := params["category"]; !ok || v.(uint) <= 0 {
		delete(params, "category")
	}
	if v, ok := params["desc"]; !ok || len(v.(string)) == 0 {
		delete(params, "desc")
	}
	if v, ok := params["language"]; !ok || len(v.(string)) == 0 {
		delete(params, "language")
	}
	return s.repo.Update(code, params)
}

func (s *DataInfoServiceImp) Delete(codes []string) (err error) {
	return s.repo.Delete(codes)
}

func (s *DataInfoServiceImp) UpdateOrCreate(code string, params map[string]interface{}) error {
	cid := uint(0)
	if v, ok := params["category"]; ok || len(v.(string)) > 0 {
		c, err := s.categoryService.GetOrCreate(v.(string), 0)
		if err == nil {
			cid = c.ID
		}
	}
	params["category"] = cid
	images := []byte{}
	if v, ok := params["images"]; ok {
		delete(params, "images")
		res, err := json.Marshal(v)
		if err == nil {
			images = res
		}
	}
	params["images"] = images
	if s.Exists(code) {
		delete(params, "code")
		return s.Update(code, params)
	}
	ts := []models.Tag{}
	if _, ok := params["tags"]; ok {
		for _, t := range params["tags"].([]string) {
			tag, err := s.tagService.GetOrCreate(t, 0)
			if err != nil {
				logrus.Warn(err)
				continue
			}
			ts = append(ts, *tag)
		}
	}
	data := models.DataInfo{}
	err := human.Encode(params, &data)
	if err != nil {
		logrus.Warning(err)
		return err
	}
	data.Tags = ts
	// data.Images = images
	return s.repo.Create(&data)
}
