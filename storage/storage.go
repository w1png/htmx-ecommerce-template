package storage

import (
	"github.com/w1png/htmx-template/config"
	"github.com/w1png/htmx-template/errors"
	"github.com/w1png/htmx-template/models"
)

type Storage interface {
	CreateUser(user *models.User) error
	GetUserById(id uint) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	GetUsers(offset, limit int) ([]*models.User, error)
	GetAllUsersByUsernameFuzzy(username string) ([]*models.User, error)
	GetUsersByUsernameFuzzy(username string, offset, limit int) ([]*models.User, error)
	GetUsersCount() (int, error)
	UpdateUser(user *models.User) error
	DeleteUserById(id uint) error

	CreateCategory(category *models.Category) error
	GetCategoryById(id uint) (*models.Category, error)
	GetCategoryBySlug(slug string) (*models.Category, error)
	GetAllCategories() ([]*models.Category, error)
	GetCategories(offset, limit int) ([]*models.Category, error)
	GetCategoriesCount() (int, error)
	GetAllCategoriesByNameFuzzy(name string) ([]*models.Category, error)
	GetCategoriesByNameFuzzy(name string, offset, limit int) ([]*models.Category, error)
	UpdateCategory(category *models.Category) error
	DeleteCategoryById(id uint) error
}

var StorageInstance Storage

func InitStorage() error {
	var err error
	switch config.ConfigInstance.StorageType {
	case "sqlite":
		if StorageInstance, err = NewSQLiteStorage(); err != nil {
			return err
		}
	default:
		return errors.NewUnknownDatabaseTypeError(config.ConfigInstance.StorageType)
	}

	return nil
}
