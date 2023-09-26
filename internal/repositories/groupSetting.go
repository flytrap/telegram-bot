package repositories

import (
	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type GroupSettingRepository interface {
	Get(code string) (*models.GroupSetting, error)
	List(q string, offset int64, limit int64, ordering string, data interface{}) (int64, error)

	Create(*models.GroupSetting) error
	Update(code string, info map[string]interface{}) (err error)
	Delete(codes []string) (err error)
}

func NewGroupSettingRepository(db *gorm.DB) GroupSettingRepository {
	return &groupSettingRepositoryImp{Db: db}
}

type groupSettingRepositoryImp struct {
	Db *gorm.DB
}

func (s *groupSettingRepositoryImp) Get(code string) (data *models.GroupSetting, err error) {
	if err = s.Db.First(&data, "code = ?", code).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (s *groupSettingRepositoryImp) List(q string, offset int64, limit int64, ordering string, data interface{}) (n int64, err error) {
	query := s.Db
	if len(q) > 0 {
		query = query.Where("name like ?", q+"%")
	}
	if len(ordering) == 0 {
		ordering = "id desc"
	}

	if err := query.Offset(int(offset)).Limit(int(limit)).Order(ordering).Find(data).Error; err != nil {
		return 0, err
	}
	query.Count(&n)
	return n, nil
}

func (s *groupSettingRepositoryImp) Create(data *models.GroupSetting) (err error) {
	result := s.Db.Create(data)
	return errors.WithStack(result.Error)
}

func (s *groupSettingRepositoryImp) Delete(codes []string) (err error) {
	result := s.Db.Where("code in ?", codes).Delete(&models.GroupSetting{})
	return errors.WithStack(result.Error)
}

func (s *groupSettingRepositoryImp) Update(code string, info map[string]interface{}) (err error) {
	result := s.Db.Model(&models.GroupSetting{Code: code}).Where("code = ?", code).Updates(info)
	return errors.WithStack(result.Error)
}
