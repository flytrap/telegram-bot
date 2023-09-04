package services

import (
	"errors"

	"github.com/flytrap/telegram-bot/internal/repositories"
)

type CategoryService interface {
	Load() error // 加载所有分类信息
	GetName(id uint) (string, error)
	GetId(name string) (uint, error)
}

func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	c := CategoryServiceImp{repo: repo, idMap: map[uint]string{}, nameMap: map[string]uint{}}
	c.Load()
	return &c
}

type CategoryServiceImp struct {
	repo    repositories.CategoryRepository
	idMap   map[uint]string
	nameMap map[string]uint
}

func (s *CategoryServiceImp) Load() error {
	results, err := s.repo.GetAll(0, 10000)
	if err != nil {
		return err
	}
	for _, item := range *results {
		s.idMap[item.ID] = item.Name
		s.nameMap[item.Name] = item.ID
	}
	return nil
}

func (s *CategoryServiceImp) GetName(id uint) (string, error) {
	name, ok := s.idMap[id]
	if ok {
		return name, nil
	}
	return "", errors.New("not found")
}

func (s *CategoryServiceImp) GetId(name string) (uint, error) {
	id, ok := s.nameMap[name]
	if ok {
		return id, nil
	}
	return 0, errors.New("not found")
}
