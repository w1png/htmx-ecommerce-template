package sqlite

import (
	"fmt"

	"github.com/w1png/go-htmx-ecommerce-template/errors"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"gorm.io/gorm"
)

func (s *SqliteStorage) CreateCart(cart *models.Cart) error {
	return s.DB.Create(cart).Error
}

func (s *SqliteStorage) GetCartByUUID(uuid string) (*models.Cart, error) {
	var cart models.Cart
	if err := s.DB.Where("uuid = ?", uuid).First(&cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewObjectNotFoundError(fmt.Sprintf("Cart with uuid: %s", uuid))
		}
		return nil, err
	}

	return &cart, nil
}

func (s *SqliteStorage) GetCartById(id uint) (*models.Cart, error) {
	var cart models.Cart
	if err := s.DB.Where("id = ?", id).First(&cart).Error; err != nil {
		return nil, err
	}

	return &cart, nil
}
