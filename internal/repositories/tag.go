package repositories

import (
	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TagRepository interface {
	Get(id uint) (*models.Tag, error)
	QueryGroup(name string, pageSize int, page int) (data *[]models.Tag, err error)
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

func (s *TagRepositoryImp) Get(id uint) (data *models.Tag, err error) {
	if err = s.Db.First(&data, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (s *TagRepositoryImp) QueryGroup(name string, pageSize int, page int) (data *[]models.Tag, err error) {
	if err := s.Db.Where("name = ?", name+"%").Limit(pageSize).Offset((page - 1) * pageSize).Find(&data).Error; err != nil {
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
