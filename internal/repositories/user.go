package repositories

import (
	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UserRepository interface {
	Get(userId int64) (*models.User, error)
	List(q string, offset int64, limit int64, ordering string) (int64, []*models.User, error)

	Create(*models.User) error
	Update(userId int64, info map[string]interface{}) (err error)
	Delete(userIds []int64) (err error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImp{Db: db}
}

type UserRepositoryImp struct {
	Db *gorm.DB
}

func (s *UserRepositoryImp) Get(userId int64) (data *models.User, err error) {
	if err = s.Db.First(&data, "user_id = ?", userId).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (s *UserRepositoryImp) List(q string, offset int64, limit int64, ordering string) (n int64, data []*models.User, err error) {
	query := s.Db
	if len(q) > 0 {
		query = query.Where("username like ?", q+"%")
	}
	if len(ordering) == 0 {
		ordering = "id desc"
	}

	if err := query.Offset(int(offset)).Limit(int(limit)).Order(ordering).Find(&data).Error; err != nil {
		return 0, nil, err
	}
	query.Count(&n)
	return n, data, nil
}

func (s *UserRepositoryImp) Create(data *models.User) (err error) {
	result := s.Db.Create(data)
	return errors.WithStack(result.Error)
}

func (s *UserRepositoryImp) Delete(userIds []int64) (err error) {
	result := s.Db.Where("user_id in ?", userIds).Delete(&models.User{})
	return errors.WithStack(result.Error)
}

func (s *UserRepositoryImp) Update(userId int64, info map[string]interface{}) (err error) {
	result := s.Db.Model(&models.User{UserID: userId}).Where("user_id = ?", userId).Updates(info)
	return errors.WithStack(result.Error)
}
