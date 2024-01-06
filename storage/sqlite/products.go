package sqlite

import (
	"fmt"

	"github.com/w1png/go-htmx-ecommerce-template/errors"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"gorm.io/gorm"
)

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

func (s *SqliteStorage) GetEnabledProducts(offset, limit int) ([]*models.Product, error) {
	var products []*models.Product
	if err := s.DB.Where("is_enabled = ?", true).Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (s *SqliteStorage) GetFeaturedProducts() ([]*models.Product, error) {
	var products []*models.Product
	if err := s.DB.Where("is_featured = ?", true).Where("is_enabled = ?", true).Find(&products).Error; err != nil {
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
