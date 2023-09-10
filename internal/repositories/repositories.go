package repositories

import (
	"strings"

	"github.com/flytrap/telegram-bot/internal/config"
	"github.com/flytrap/telegram-bot/internal/models"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var RepositorySet = wire.NewSet(NewDataInfoRepository, NewDataTagRepository, NewCategoryRepository, NewUserRepository, NewAdRepository)

func AutoMigrate(db *gorm.DB) error {
	if dbType := config.C.Gorm.DBType; strings.ToLower(dbType) == "mysql" {
		db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	}

	return db.AutoMigrate(
		new(models.DataInfo),
		new(models.Category),
		new(models.Tag),
		new(models.User),
		new(models.Ad),
	)
}
