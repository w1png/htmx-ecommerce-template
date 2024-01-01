package storage

import (
	"fmt"

	"github.com/w1png/htmx-template/config"
	"github.com/w1png/htmx-template/errors"
	"github.com/w1png/htmx-template/models"
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
	)
}

func (s *SqliteStorage) CreateUser(user *models.User) error {
	return s.DB.Create(user).Error
}

func (s *SqliteStorage) GetUserById(id uint) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewObjectNotFoundError(fmt.Sprintf("User with id: %d", id))
		}
		return nil, err
	}

	return &user, nil
}

func (s *SqliteStorage) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewObjectNotFoundError(fmt.Sprintf("User with username: %s", username))
		}
		return nil, err
	}

	return &user, nil
}

func (s *SqliteStorage) GetAllUsersByUsernameFuzzy(username string) ([]*models.User, error) {
	var users []*models.User
	if err := s.DB.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username)).Find(&users).Error; err != nil {
		return users, err
	}

	return users, nil
}

func (s *SqliteStorage) GetUsersByUsernameFuzzy(username string, offset, limit int) ([]*models.User, error) {
	var users []*models.User
	if err := s.DB.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username)).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return users, err
	}

	return users, nil
}

func (s *SqliteStorage) GetAllUsers() ([]*models.User, error) {
	var users []*models.User
	if err := s.DB.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (s *SqliteStorage) GetUsers(offset, limit int) ([]*models.User, error) {
	var users []*models.User
	if err := s.DB.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (s *SqliteStorage) GetUsersCount() (int, error) {
	var count int64
	if err := s.DB.Model(&models.User{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (s *SqliteStorage) UpdateUser(user *models.User) error {
	if _, err := s.GetUserById(user.ID); err != nil {
		return err
	}

	return s.DB.Save(user).Error
}

func (s *SqliteStorage) DeleteUserById(id uint) error {
	if id == 1 {
		return errors.NewMainAdminDeletionError("Cannot delete main admin")
	}

	if _, err := s.GetUserById(id); err != nil {
		return err
	}

	return s.DB.Delete(&models.User{}, id).Error
}

func (s *SqliteStorage) CreateCategory(category *models.Category) error {
	if err := s.DB.Create(category).Error; err != nil && err.Error() == "UNIQUE constraint failed: categories.slug" {
		return errors.NewSlugNotUniqueError(category.Slug)
	} else if err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) GetCategoryById(id uint) (*models.Category, error) {
	var category models.Category
	if err := s.DB.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewObjectNotFoundError(fmt.Sprintf("Category with id: %d", id))
		}
		return nil, err
	}

	return &category, nil
}

func (s *SqliteStorage) GetCategoryBySlug(slug string) (*models.Category, error) {
	var category models.Category
	if err := s.DB.Where("slug = ?", slug).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewObjectNotFoundError(fmt.Sprintf("Category with slug: %s", slug))
		}
		return nil, err
	}

	return &category, nil
}

func (s *SqliteStorage) GetAllCategories() ([]*models.Category, error) {
	var categories []*models.Category
	if err := s.DB.Find(&categories).Error; err != nil {
		return categories, err
	}

	return categories, nil
}

func (s *SqliteStorage) GetCategories(offset, limit int) ([]*models.Category, error) {
	var categories []*models.Category
	if err := s.DB.Offset(offset).Limit(limit).Find(&categories).Error; err != nil {
		return categories, err
	}

	return categories, nil
}

func (s *SqliteStorage) GetCategoriesCount() (int, error) {
	var count int64
	if err := s.DB.Model(&models.Category{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (s *SqliteStorage) GetAllCategoriesByNameFuzzy(name string) ([]*models.Category, error) {
	var categories []*models.Category
	if err := s.DB.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name)).Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *SqliteStorage) GetCategoriesByNameFuzzy(name string, offset, limit int) ([]*models.Category, error) {
	var categories []*models.Category
	if err := s.DB.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name)).Offset(offset).Limit(limit).Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *SqliteStorage) UpdateCategory(category *models.Category) error {
	if _, err := s.GetCategoryById(category.ID); err != nil {
		return err
	}

	return s.DB.Save(category).Error
}

func (s *SqliteStorage) DeleteCategoryById(id uint) error {
	if _, err := s.GetCategoryById(id); err != nil {
		return err
	}

	return s.DB.Unscoped().Delete(&models.Category{}, id).Error
}

func (s *SqliteStorage) CreateProduct(product *models.Product) error {
	if err := s.DB.Create(product).Error; err != nil && err.Error() == "UNIQUE constraint failed: products.slug" {
		return errors.NewSlugNotUniqueError(fmt.Sprintf("Product with slug: %s", product.Slug))
	} else if err != nil {
		return err
	}

	return nil
}

func (s *SqliteStorage) GetProductById(id uint) (*models.Product, error) {
	var product models.Product
	if err := s.DB.Where("id = ?", id).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewObjectNotFoundError(fmt.Sprintf("Product with id: %d", id))
		}
		return nil, err
	}

	return &product, nil
}

func (s *SqliteStorage) GetProductBySlug(slug string) (*models.Product, error) {
	var product models.Product
	if err := s.DB.Where("slug = ?", slug).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewObjectNotFoundError(fmt.Sprintf("Product with slug: %s", slug))
		}
		return nil, err
	}

	return &product, nil
}

func (s *SqliteStorage) GetAllProducts() ([]*models.Product, error) {
	var products []*models.Product
	if err := s.DB.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (s *SqliteStorage) GetProducts(offset, limit int) ([]*models.Product, error) {
	var products []*models.Product
	if err := s.DB.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (s *SqliteStorage) GetProductsCount() (int, error) {
	var count int64
	if err := s.DB.Model(&models.Product{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (s *SqliteStorage) GetAllProductsByNameFuzzy(name string) ([]*models.Product, error) {
	var products []*models.Product
	if err := s.DB.Where("name LIKE ?", "%"+name+"%").Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (s *SqliteStorage) GetProductsByNameFuzzy(name string, offset, limit int) ([]*models.Product, error) {
	var products []*models.Product
	if err := s.DB.Where("name LIKE ?", "%"+name+"%").Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (s *SqliteStorage) GetProductsByTags(tag string, offset, limit int) ([]*models.Product, error) {
	var products []*models.Product
	if err := s.DB.Where("tags LIKE ?", "%"+tag+"%").Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (s *SqliteStorage) UpdateProduct(product *models.Product) error {
	if err := s.DB.Save(product).Error; err != nil && err.Error() == "UNIQUE constraint failed: products.slug" {
		return errors.NewSlugNotUniqueError(fmt.Sprintf("Product with slug: %s", product.Slug))
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return errors.NewObjectNotFoundError(fmt.Sprintf("Product with id: %d", product.ID))
	} else if err != nil {
		return err
	}

	return s.DB.Save(product).Error
}

func (s *SqliteStorage) DeleteProductById(id uint) error {
	if _, err := s.GetProductById(id); err != nil {
		return err
	}

	return s.DB.Unscoped().Delete(&models.Product{}, id).Error
}
