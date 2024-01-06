package sqlite

import (
	"fmt"

	"github.com/w1png/go-htmx-ecommerce-template/errors"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"gorm.io/gorm"
)

func (s *SqliteStorage) CreateCartProduct(cartProduct *models.CartProduct) error {
	return s.DB.Create(cartProduct).Error
}

func (s *SqliteStorage) GetCartProductById(id uint) (*models.CartProduct, error) {
	var cartProduct models.CartProduct
	if err := s.DB.Where("id = ?", id).First(&cartProduct).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewObjectNotFoundError(fmt.Sprintf("CartProduct with id: %d", id))
		}
		return nil, err
	}

	return &cartProduct, nil
}

func (s *SqliteStorage) GetCartProductByProductIdAndCartID(id uint, cart_id uint) (*models.CartProduct, error) {
	var cartProduct models.CartProduct
	if err := s.DB.Where("product_id = ?", id).Where("cart_id = ?", cart_id).First(&cartProduct).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewObjectNotFoundError(fmt.Sprintf("CartProduct with id: %d and cart id: %d", id, cart_id))
		}
		return nil, err
	}

	return &cartProduct, nil
}

func (s *SqliteStorage) UpdateCartProduct(cartProduct *models.CartProduct) error {
	if _, err := s.GetCartProductById(cartProduct.ID); err != nil {
		return err
	}
	return s.DB.Save(cartProduct).Error
}

func (s *SqliteStorage) DeleteCartProductById(id uint) error {
	if _, err := s.GetCartProductById(id); err != nil {
		return err
	}
	return s.DB.Unscoped().Delete(&models.CartProduct{}, id).Error
}

func (s *SqliteStorage) CreateOrder(order *models.Order) error {
	return s.DB.Create(order).Error
}
