package services

import (
	"errors"
	"strings"

	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/flytrap/telegram-bot/internal/repositories"
)

type CategoryService interface {
	List(q string, page int64, size int64, ordering string, data interface{}) (n int64, err error)
	Update(id uint, name string, weight int32) error
	Create(name string, weight int32) error
	Delete(ids []uint) (err error)

	Load() error // 加载所有分类信息
	GetName(id uint) (string, error)
	GetId(name string) (uint, error)
	GetOrCreate(name string, weight int32) (*models.Category, error)
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

func (s *CategoryServiceImp) List(q string, page int64, size int64, ordering string, data interface{}) (n int64, err error) {
	n, err = s.repo.List(q, (page-1)*size, size, ordering, data)
	return
}

func (s *CategoryServiceImp) Create(name string, weight int32) error {
	data := models.Category{Name: name, Weight: weight}
	return s.repo.Create(&data)
}

func (s *CategoryServiceImp) GetOrCreate(name string, weight int32) (*models.Category, error) {
	t, err := s.repo.Get(name)
	if t == nil {
		t = &models.Category{Name: name, Weight: weight}
		err = s.repo.Create(t)
		if err != nil {
			return t, err
		}
	}
	return t, err
}

func (s *CategoryServiceImp) Update(id uint, name string, weight int32) error {
	err := s.repo.Update(id, name, weight)
	if err != nil {
		return err
	}
	return nil
}

func (s *CategoryServiceImp) Delete(ids []uint) (err error) {
	return s.repo.Delete(ids)
}

func (s *CategoryServiceImp) Load() error {
	results := []*models.Category{}
	_, err := s.repo.List("", 0, 10000, "", &results)
	if err != nil {
		return err
	}
	for _, item := range results {
		s.idMap[item.ID] = item.Name
		s.nameMap[item.Name] = item.ID
	}
	return nil
}

func (s *CategoryServiceImp) GetName(id uint) (string, error) {
	name, ok := s.idMap[id]
	if ok {
		ns := strings.Split(name, ":")
		if len(ns) > 1 {
			name = ns[len(ns)-1]
		}
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
