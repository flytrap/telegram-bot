package services

import (
	"errors"

	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/flytrap/telegram-bot/internal/repositories"
	"github.com/sirupsen/logrus"
)

func NewDataService(repo repositories.DataInfoRepository, tagService DataTagService, categoryService CategoryService) DataService {
	return &DataInfoServiceImp{repo: repo, tagService: tagService, categoryService: categoryService}
}

type DataService interface {
	List(q string, category string, language string, page int64, size int64, ordering string, data interface{}) (int64, error)
	SearchTag(tag string, page int64, size int64, data interface{}) (total int64, err error)
	GetNeedUpdateCode(days int, page int64, size int64) ([]string, error)
	Update(code string, tid int64, name string, desc string, num uint32, weight int, lang string, category uint) error
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

func (s *DataInfoServiceImp) Update(code string, tid int64, name string, desc string, num uint32, weight int, lang string, category uint) error {
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
	if category > 0 {
		params["category"] = category
	}
	if len(desc) > 0 {
		params["desc"] = desc
	}
	if len(lang) > 0 {
		params["language"] = lang
	}
	return s.repo.Update(code, params)
}

func (s *DataInfoServiceImp) Delete(codes []string) (err error) {
	return s.repo.Delete(codes)
}

func (s *DataInfoServiceImp) UpdateOrCreate(code string, tid int64, name string, desc string, num uint32, tags []string, category string, lang string) error {
	cid := uint(0)
	if len(category) > 0 {
		c, err := s.categoryService.GetOrCreate(category, 0)
		if err == nil {
			cid = c.ID
		}
	}
	if s.Exists(code) {
		return s.Update(code, tid, name, desc, num, 0, lang, cid)
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
	data := models.DataInfo{Code: code, Tid: tid, Name: name, Desc: desc, Number: num, Tags: ts, Category: cid, Language: lang}
	return s.repo.Create(&data)
}
