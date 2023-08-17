package services

import (
	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/flytrap/telegram-bot/internal/repositories"
	"github.com/flytrap/telegram-bot/internal/serializers"
	"github.com/sirupsen/logrus"
)

func NewGroupService(repo repositories.GroupRepository, tagService TagService) GroupService {
	return &GroupServiceImp{repo: repo, tagService: tagService}
}

type GroupService interface {
	SearchTag(tag string, size int, page int) (data []*serializers.GroupSerilizer, err error)
	GetNeedUpdateCode(days int, pageSize int, page int) ([]string, error)
	Update(code string, tid int64, name string, desc string, num uint32) error
	Delete(code string) (err error)

	UpdateOrCreate(code string, tid int64, name string, desc string, num uint32, tags []string, category string) error
}

type GroupServiceImp struct {
	repo       repositories.GroupRepository
	tagService TagService
}

func (s *GroupServiceImp) Exists(code string) bool {
	data, _ := s.repo.Get(code)
	return data != nil
}

func (s *GroupServiceImp) SearchTag(tag string, pageSize int, page int) (data []*serializers.GroupSerilizer, err error) {
	res, err := s.repo.QueryTag(tag, pageSize, page)
	data = []*serializers.GroupSerilizer{}
	for _, item := range res {
		data = append(data, &serializers.GroupSerilizer{Code: item.Code, Name: item.Name, Number: item.HumanSize()})
	}
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *GroupServiceImp) GetNeedUpdateCode(days int, pageSize int, page int) ([]string, error) {
	res, err := s.repo.GetNeedUpdate(days, pageSize, page)
	if err != nil {
		return nil, err
	}
	results := []string{}
	for _, item := range res {
		results = append(results, item.Code)
	}
	return results, nil
}

func (s *GroupServiceImp) Update(code string, tid int64, name string, desc string, num uint32) error {
	var params = map[string]interface{}{"tid": tid, "name": name, "number": num}
	if len(desc) > 0 {
		params["desc"] = desc
	}
	return s.repo.Update(code, params)
}

func (s *GroupServiceImp) Delete(code string) (err error) {
	return s.repo.Delete(code)
}

func (s *GroupServiceImp) UpdateOrCreate(code string, tid int64, name string, desc string, num uint32, tags []string, category string) error {
	if s.Exists(code) {
		return s.Update(code, tid, name, desc, num)
	}
	ts := []*models.Tag{}
	for _, t := range tags {
		tag, err := s.tagService.GetOrCreate(t)
		if err != nil {
			logrus.Warn(err)
			continue
		}
		ts = append(ts, tag)
	}
	data := models.Group{Code: code, Tid: tid, Name: name, Desc: desc, Number: num, Tags: ts}
	err := s.repo.Create(&data)
	return err
}
