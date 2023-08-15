package services

import (
	"github.com/flytrap/telegram-bot/internal/repositories"
	"github.com/flytrap/telegram-bot/internal/serializers"
)

func NewGroupService(repo repositories.GroupRepository) GroupService {
	return &GroupServiceImp{repo: repo}
}

type GroupService interface {
	SearchTag(tag string, size int, page int) (data []*serializers.GroupSerilizer, err error)
	GetNeedUpdateCode(days int, pageSize int, page int) ([]string, error)
	Update(code string, tid int, name string, desc string, num int) error
	Delete(code string) (err error)
}

type GroupServiceImp struct {
	repo repositories.GroupRepository
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

func (s *GroupServiceImp) Update(code string, tid int, name string, desc string, num int) error {
	var params = map[string]interface{}{"tid": tid, "name": name, "number": num}
	if len(desc) > 0 {
		params["desc"] = desc
	}
	return s.repo.Update(code, params)
}

func (s *GroupServiceImp) Delete(code string) (err error) {
	return s.repo.Delete(code)
}
