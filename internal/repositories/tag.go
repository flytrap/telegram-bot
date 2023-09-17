package repositories

import (
	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type DataTagRepository interface {
	Get(name string) (*models.Tag, error)
	GetMany(ids []uint) (data *[]models.Tag, err error)
	List(q string, offset int64, limit int64, ordering string) (int64, []*models.Tag, error)

	Create(*models.Tag) error
	Update(id uint, name string, weight int32) error
	Delete(ids []uint) (err error)
}

func NewDataTagRepository(db *gorm.DB) DataTagRepository {
	return &DataTagRepositoryImp{Db: db}
}

type DataTagRepositoryImp struct {
	Db *gorm.DB
}

func (s *DataTagRepositoryImp) Get(name string) (data *models.Tag, err error) {
	if err = s.Db.First(&data, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (s *DataTagRepositoryImp) GetMany(ids []uint) (data *[]models.Tag, err error) {
	if err := s.Db.Where("id in ?", ids).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (s *DataTagRepositoryImp) List(q string, offset int64, limit int64, ordering string) (n int64, data []*models.Tag, err error) {
	query := s.Db.Model(models.Tag{})
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

func (s *DataTagRepositoryImp) Create(data *models.Tag) (err error) {
	result := s.Db.Create(data)
	return errors.WithStack(result.Error)
}

func (s *DataTagRepositoryImp) Update(id uint, name string, weight int32) error {
	info := map[string]interface{}{"weight": weight}
	if len(name) > 0 {
		info["name"] = name
	}
	result := s.Db.Model(&models.Tag{}).Where("id = ?", id).Updates(info)
	return errors.WithStack(result.Error)
}

func (s *DataTagRepositoryImp) Delete(ids []uint) (err error) {
	result := s.Db.Where("id in ?", ids).Delete(models.Tag{})
	return errors.WithStack(result.Error)
}
