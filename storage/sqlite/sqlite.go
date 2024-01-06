package sqlite

import (
	"github.com/w1png/go-htmx-ecommerce-template/config"
	"github.com/w1png/go-htmx-ecommerce-template/errors"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SqliteStorage struct {
	DB *gorm.DB
}

func NewSQLiteStorage() (*SqliteStorage, error) {
	storage := &SqliteStorage{}

	var err error
	if storage.DB, err = gorm.Open(sqlite.Open(config.ConfigInstance.SqlitePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}); err != nil {
		return nil, errors.NewDatabaseConnectionError(err.Error())
	}

	if err := storage.autoMigrate(); err != nil {
		return nil, errors.NewDatabaseMigrationError(err.Error())
	}

	return storage, nil
}

func (s *SqliteStorage) autoMigrate() error {
	return s.DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.CartProduct{},
		&models.Cart{},
		&models.OrderProduct{},
		&models.Order{},
	)
}
