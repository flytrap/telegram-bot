package repositories

import (
	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type DataTagRepository interface {
	Get(name string) (*models.DataTag, error)
	GetMany(ids []uint) (data *[]models.DataTag, err error)

	Create(*models.DataTag) error
	Delete(id uint) (err error)
}

func NewDataTagRepository(db *gorm.DB) DataTagRepository {
	return &DataTagRepositoryImp{Db: db}
}

type DataTagRepositoryImp struct {
	Db *gorm.DB
}

func (s *DataTagRepositoryImp) Get(name string) (data *models.DataTag, err error) {
	if err = s.Db.First(&data, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (s *DataTagRepositoryImp) GetMany(ids []uint) (data *[]models.DataTag, err error) {
	if err := s.Db.Where("id in ?", ids).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (s *DataTagRepositoryImp) Create(data *models.DataTag) (err error) {
	result := s.Db.Create(data)
	return errors.WithStack(result.Error)
}

func (s *DataTagRepositoryImp) Delete(id uint) (err error) {
	result := s.Db.Where("id=?", id).Delete(models.DataTag{})
	return errors.WithStack(result.Error)
}
