package repositories

import (
	"time"

	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type DataInfoRepository interface {
	Get(code string) (*models.DataInfo, error)
	List(q string, category uint, language string, offset int64, limit int64, ordering string, data interface{}) (int64, error)
	GetNeedUpdate(days int, offset int64, limit int64) (data []*models.DataInfo, err error)
	QueryTag(tag string, offset int64, limit int64, data interface{}) (int64, error)

	Create(*models.DataInfo) error
	Update(code string, info map[string]interface{}) (err error)
	Delete(codes []string) (err error)
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

func (s *DataInfoRepositoryImp) List(q string, category uint, language string, offset int64, limit int64, ordering string, data interface{}) (n int64, err error) {
	query := s.Db.Model(models.DataInfo{})
	if len(q) > 0 {
		query = query.Where("name like ?", q+"%")
	}
	if category > 0 {
		query = query.Where("category = ?", category)
	}
	if len(language) > 0 {
		query = query.Where("language = ?", language)
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

func (s *DataInfoRepositoryImp) QueryTag(tag string, offset int64, limit int64, data interface{}) (n int64, err error) {
	qs := s.Db.Model(&models.DataInfo{}).Select("tg_data_info.name", "tg_data_info.tid", "tg_data_info.code", "tg_data_info.type", "tg_data_info.number").Joins("inner JOIN tg_data_tag on tg_data_tag.data_info_id=tg_data_info.id").Joins("LEFT JOIN tg_tag on tg_tag.id=tg_data_tag.tag_id").Where("tg_tag.name = ?", tag).Preload("Tags").Limit(int(limit)).Offset(int(offset)).Order("number desc")
	if err = qs.Find(data).Error; err != nil {
		return 0, nil
	}
	qs.Count(&n)
	return n, nil
}

func (s *DataInfoRepositoryImp) GetNeedUpdate(days int, offset int64, limit int64) (data []*models.DataInfo, err error) {
	if err := s.Db.Select("code").Where("tid = 0 or updated_at < ?", time.Now().AddDate(0, 0, -days)).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (s *DataInfoRepositoryImp) Create(data *models.DataInfo) (err error) {
	result := s.Db.Create(data)
	return errors.WithStack(result.Error)
}

func (s *DataInfoRepositoryImp) Delete(codes []string) (err error) {
	result := s.Db.Where("code in ?", codes).Delete(&models.DataInfo{})
	return errors.WithStack(result.Error)
}

func (s *DataInfoRepositoryImp) Update(code string, info map[string]interface{}) (err error) {
	result := s.Db.Model(&models.DataInfo{Code: code}).Where("code = ?", code).Updates(info)
	return errors.WithStack(result.Error)
}
