package repositories

import (
	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type AdRepository interface {
	Get(id uint) (*models.Ad, error)
	GetMany(q string, offset int64, limit int64) ([]*models.Ad, error)

	Create(*models.Ad) error
	Update(id uint, info map[string]interface{}) (err error)
	Delete(ids []uint) (err error)
}

func NewAdRepository(db *gorm.DB) AdRepository {
	return &AdRepositoryImp{Db: db}
}

type AdRepositoryImp struct {
	Db *gorm.DB
}

func (s *AdRepositoryImp) Get(id uint) (data *models.Ad, err error) {
	if err = s.Db.First(&data, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (s *AdRepositoryImp) GetMany(q string, offset int64, limit int64) (data []*models.Ad, err error) {
	query := s.Db
	if len(q) > 0 {
		query = query.Where("name like ?", q+"%")
	}

	if err := query.Offset(int(offset)).Limit(int(limit)).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (s *AdRepositoryImp) Create(data *models.Ad) (err error) {
	result := s.Db.Create(data)
	return errors.WithStack(result.Error)
}

func (s *AdRepositoryImp) Delete(ids []uint) (err error) {
	result := s.Db.Where("id in ?", ids).Delete(&models.Ad{})
	return errors.WithStack(result.Error)
}

func (s *AdRepositoryImp) Update(id uint, info map[string]interface{}) (err error) {
	result := s.Db.Model(&models.Ad{}).Where("id = ?", id).Updates(info)
	return errors.WithStack(result.Error)
}
