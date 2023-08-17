package services

import (
	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/flytrap/telegram-bot/internal/repositories"
)

type TagService interface {
	GetOrCreate(name string) (*models.Tag, error)
}

func NewTagService(repo repositories.TagRepository) TagService {
	return &TagServiceImp{repo: repo}
}

type TagServiceImp struct {
	repo repositories.TagRepository
}

func (s *TagServiceImp) GetOrCreate(name string) (*models.Tag, error) {
	t, err := s.repo.Get(name)
	if t == nil {
		t = &models.Tag{Name: name}
		err = s.repo.Create(t)
	}
	return t, err
}
