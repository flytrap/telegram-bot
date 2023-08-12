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
