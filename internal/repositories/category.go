package repositories

import (
	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Get(name string) (*models.Category, error)
	Query(name string, size int) (data *[]models.Category, err error)
	GetMany(ids []uint) (data *[]models.Category, err error)
	List(q string, offset int64, limit int64, ordering string) (n int64, data []*models.Tag, err error)

	Create(*models.Category) error
	Update(id uint, name string, weight int32) error
	Delete(ids []uint) (err error)
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &CategoryRepositoryImp{Db: db}
}

type CategoryRepositoryImp struct {
	Db *gorm.DB
}

func (s *CategoryRepositoryImp) Get(name string) (data *models.Category, err error) {
	if err = s.Db.First(&data, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (s *CategoryRepositoryImp) Query(name string, size int) (data *[]models.Category, err error) {
	if err := s.Db.Where("name like ?", name+"%").Limit(size).Find(&data).Error; err != nil {
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
func (s *CategoryRepositoryImp) List(q string, offset int64, limit int64, ordering string) (n int64, data []*models.Tag, err error) {
	query := s.Db
	if len(q) > 0 {
		query = query.Where("name like ?", q+"%")
	}
	if len(ordering) == 0 {
		ordering = "id desc"
	}

	if err = query.Offset(int(offset)).Limit(int(limit)).Order(ordering).Find(&data).Error; err != nil {
		return 0, nil, err
	}
	query.Count(&n)
	return n, data, nil
}

func (s *CategoryRepositoryImp) Create(data *models.Category) (err error) {
	result := s.Db.Create(data)
	return errors.WithStack(result.Error)
}

func (s *CategoryRepositoryImp) Update(id uint, name string, weight int32) error {
	info := map[string]interface{}{"weight": weight}
	if len(name) > 0 {
		info["name"] = name
	}
	result := s.Db.Model(&models.Category{}).Where("id = ?", id).Updates(info)
	return errors.WithStack(result.Error)
}

func (s *CategoryRepositoryImp) Delete(ids []uint) (err error) {
	result := s.Db.Where("id in ?", ids).Delete(models.Category{})
	return errors.WithStack(result.Error)
}
