package repositories

import (
	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type GroupRepository interface {
	Get(id uint) (*models.Group, error)
	Query(name string, size int) (data *[]models.Group, err error)
	GetMany(ids []uint) (data *[]models.Group, err error)
	QueryTag(tag string, pageSize int, page int) ([]*models.Group, error)

	Create(*models.Group) error
	Delete(id uint) (err error)
}

func NewGroupRepository(db *gorm.DB) GroupRepository {
	return &GroupRepositoryImp{Db: db}
}

type GroupRepositoryImp struct {
	Db *gorm.DB
}

func (s *GroupRepositoryImp) Get(id uint) (data *models.Group, err error) {
	if err = s.Db.First(&data, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (s *GroupRepositoryImp) Query(name string, size int) (data *[]models.Group, err error) {
	if err := s.Db.Where("en_name like ?", name+"%").Limit(size).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (s *GroupRepositoryImp) QueryTag(tag string, pageSize int, page int) (data []*models.Group, err error) {
	q := s.Db.Model(&models.Group{}).Select("tg_group.name", "tg_group.tid", "tg_group.code", "tg_group.type", "tg_group.number").Joins("inner JOIN tg_group_tag on tg_group_tag.group_id=tg_group.id").Joins("LEFT JOIN tg_tag on tg_tag.id=tg_group_tag.tag_id").Where("tg_tag.name = ?", tag).Preload("Tags").Limit(pageSize).Offset((page - 1) * pageSize).Find(&data)
	return data, q.Error
}

func (s *GroupRepositoryImp) GetMany(ids []uint) (data *[]models.Group, err error) {
	if err := s.Db.Where("id in ?", ids).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (s *GroupRepositoryImp) Create(data *models.Group) (err error) {
	result := s.Db.Create(data)
	return errors.WithStack(result.Error)
}

func (s *GroupRepositoryImp) Delete(id uint) (err error) {
	result := s.Db.Where("id=?", id).Delete(models.Group{})
	return errors.WithStack(result.Error)
}
