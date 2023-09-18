package services

import (
	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/flytrap/telegram-bot/internal/repositories"
)

type DataTagService interface {
	List(q string, page int64, size int64, ordering string, data interface{}) (n int64, err error)
	Update(id uint, name string, weight int32) error
	Create(name string, weight int32) error
	Delete(ids []uint) (err error)

	GetOrCreate(name string, weight int32) (*models.Tag, error)
}

func NewDataTagService(repo repositories.DataTagRepository) DataTagService {
	return &DataTagServiceImp{repo: repo}
}

type DataTagServiceImp struct {
	repo repositories.DataTagRepository
}

func (s *DataTagServiceImp) List(q string, page int64, size int64, ordering string, data interface{}) (n int64, err error) {
	n, err = s.repo.List(q, (page-1)*size, size, ordering, data)
	return
}

func (s *DataTagServiceImp) Create(name string, weight int32) error {
	data := models.Tag{Name: name, Weight: weight}
	return s.repo.Create(&data)
}

func (s *DataTagServiceImp) GetOrCreate(name string, weight int32) (*models.Tag, error) {
	t, err := s.repo.Get(name)
	if t == nil {
		t = &models.Tag{Name: name, Weight: weight}
		err = s.repo.Create(t)
		if err != nil {
			return t, err
		}
	}
	return t, err
}

func (s *DataTagServiceImp) Update(id uint, name string, weight int32) error {
	err := s.repo.Update(id, name, weight)
	if err != nil {
		return err
	}
	return nil
}

func (s *DataTagServiceImp) Delete(ids []uint) (err error) {
	return s.repo.Delete(ids)
}
