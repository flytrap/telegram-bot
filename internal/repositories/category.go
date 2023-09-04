package repositories

import (
	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Get(id uint) (*models.Category, error)
	GetAll(offset int, limit int) (data *[]models.Category, err error)
	Query(name string, size int) (data *[]models.Category, err error)
	GetMany(ids []uint) (data *[]models.Category, err error)

	Create(*models.Category) error
	Delete(id uint) (err error)
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &CategoryRepositoryImp{Db: db}
}

type CategoryRepositoryImp struct {
	Db *gorm.DB
}

func (s *CategoryRepositoryImp) Get(id uint) (data *models.Category, err error) {
	if err = s.Db.First(&data, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (s *CategoryRepositoryImp) Query(name string, size int) (data *[]models.Category, err error) {
	if err := s.Db.Where("en_name like ?", name+"%").Limit(size).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (s *CategoryRepositoryImp) GetMany(ids []uint) (data *[]models.Category, err error) {
	if err := s.Db.Where("id in ?", ids).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (s *CategoryRepositoryImp) GetAll(offset int, limit int) (data *[]models.Category, err error) {
	if err := s.Db.Select("id", "name").Offset(offset).Limit(limit).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (s *CategoryRepositoryImp) Create(data *models.Category) (err error) {
	result := s.Db.Create(data)
	return errors.WithStack(result.Error)
}

func (s *CategoryRepositoryImp) Delete(id uint) (err error) {
	result := s.Db.Where("id=?", id).Delete(models.Category{})
	return errors.WithStack(result.Error)
}
