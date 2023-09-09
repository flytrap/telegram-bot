package repositories

import (
	"time"

	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type DataInfoRepository interface {
	Get(code string) (*models.DataInfo, error)
	GetMany(category uint, language string, offset int64, limit int64) ([]*models.DataInfo, error)
	GetNeedUpdate(days int, offset int64, limit int64) (data []*models.DataInfo, err error)
	QueryTag(tag string, offset int64, limit int64) ([]*models.DataInfo, error)

	Create(*models.DataInfo) error
	Update(code string, info map[string]interface{}) (err error)
	Delete(code string) (err error)
}

func NewDataInfoRepository(db *gorm.DB) DataInfoRepository {
	return &DataInfoRepositoryImp{Db: db}
}

type DataInfoRepositoryImp struct {
	Db *gorm.DB
}

func (s *DataInfoRepositoryImp) Get(code string) (data *models.DataInfo, err error) {
	if err = s.Db.First(&data, "code = ?", code).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (s *DataInfoRepositoryImp) GetMany(category uint, language string, offset int64, limit int64) (data []*models.DataInfo, err error) {
	query := s.Db
	if category > 0 {
		query = query.Where("category = ?", category)
	}
	if len(language) > 0 {
		query = query.Where("language = ?", language)
	}
	if err := query.Offset(int(offset)).Limit(int(limit)).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (s *DataInfoRepositoryImp) QueryTag(tag string, offset int64, limit int64) (data []*models.DataInfo, err error) {
	q := s.Db.Model(&models.DataInfo{}).Select("tg_group.name", "tg_group.tid", "tg_group.code", "tg_group.type", "tg_group.number").Joins("inner JOIN tg_group_tag on tg_group_tag.group_id=tg_group.id").Joins("LEFT JOIN tg_tag on tg_tag.id=tg_group_tag.tag_id").Where("tg_tag.name = ?", tag).Preload("Tags").Limit(int(limit)).Offset(int(offset)).Order("number desc").Find(&data)
	return data, q.Error
}

func (s *DataInfoRepositoryImp) GetNeedUpdate(days int, offset int64, limit int64) (data []*models.DataInfo, err error) {
	if err := s.Db.Select("code").Where("tid is null or updated_at < ?", time.Now().AddDate(0, 0, -days)).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (s *DataInfoRepositoryImp) Create(data *models.DataInfo) (err error) {
	result := s.Db.Create(data)
	return errors.WithStack(result.Error)
}

func (s *DataInfoRepositoryImp) Delete(code string) (err error) {
	result := s.Db.Where("code=?", code).Delete(&models.DataInfo{})
	return errors.WithStack(result.Error)
}

func (s *DataInfoRepositoryImp) Update(code string, info map[string]interface{}) (err error) {
	result := s.Db.Model(&models.DataInfo{Code: code}).Where("code = ?", code).Updates(info)
	return errors.WithStack(result.Error)
}