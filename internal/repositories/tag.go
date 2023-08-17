package repositories

import (
	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TagRepository interface {
	Get(name string) (*models.Tag, error)
	GetMany(ids []uint) (data *[]models.Tag, err error)

	Create(*models.Tag) error
	Delete(id uint) (err error)
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &TagRepositoryImp{Db: db}
}

type TagRepositoryImp struct {
	Db *gorm.DB
}

func (s *TagRepositoryImp) Get(name string) (data *models.Tag, err error) {
	if err = s.Db.First(&data, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (s *TagRepositoryImp) GetMany(ids []uint) (data *[]models.Tag, err error) {
	if err := s.Db.Where("id in ?", ids).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (s *TagRepositoryImp) Create(data *models.Tag) (err error) {
	result := s.Db.Create(data)
	return errors.WithStack(result.Error)
}

func (s *TagRepositoryImp) Delete(id uint) (err error) {
	result := s.Db.Where("id=?", id).Delete(models.Tag{})
	return errors.WithStack(result.Error)
}
