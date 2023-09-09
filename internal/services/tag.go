package services

import (
	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/flytrap/telegram-bot/internal/repositories"
)

type DataTagService interface {
	GetOrCreate(name string) (*models.DataTag, error)
}

func NewDataTagService(repo repositories.DataTagRepository) DataTagService {
	return &DataTagServiceImp{repo: repo}
}

type DataTagServiceImp struct {
	repo repositories.DataTagRepository
}

func (s *DataTagServiceImp) GetOrCreate(name string) (*models.DataTag, error) {
	t, err := s.repo.Get(name)
	if t == nil {
		t = &models.DataTag{Name: name}
		err = s.repo.Create(t)
	}
	return t, err
}
