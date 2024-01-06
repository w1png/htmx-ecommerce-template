package sqlite

import (
	"fmt"

	"github.com/w1png/go-htmx-ecommerce-template/errors"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"gorm.io/gorm"
)

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

func (s *SqliteStorage) GetMainCategories() ([]*models.Category, error) {
	var categories []*models.Category
	if err := s.DB.Where("parent_id = 0").Find(&categories).Error; err != nil {
		return categories, err
	}

	return categories, nil
}

func (s *SqliteStorage) GetCategoryChildren(id uint) ([]*models.Category, error) {
	var categories []*models.Category
	if err := s.DB.Where("parent_id = ?", id).Where("is_enabled = true").Find(&categories).Error; err != nil {
		return categories, err
	}

	return categories, nil
}

func (s *SqliteStorage) GetCategoryProducts(id uint) ([]*models.Product, error) {
	var products []*models.Product
	if err := s.DB.Where("category_id = ?", id).Where("is_enabled = true").Find(&products).Error; err != nil {
		return products, err
	}

	return products, nil
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
