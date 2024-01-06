package sqlite

import (
	"fmt"

	"github.com/w1png/go-htmx-ecommerce-template/errors"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"gorm.io/gorm"
)

func (s *SqliteStorage) GetOrderById(id uint) (*models.Order, error) {
	var order *models.Order
	if err := s.DB.Where("id = ?", id).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewObjectNotFoundError(fmt.Sprintf("Order with id: %d", id))
		}
		return nil, err
	}

	return order, nil
}

func (s *SqliteStorage) GetOrders(status models.OrderStatus, offset, limit int) ([]*models.Order, error) {
	var orders []*models.Order
	tx := s.DB
	if status != models.OrderStatusAny {
		tx = s.DB.Where("status = ?", status)
	}

	if err := tx.Offset(offset).Limit(limit).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *SqliteStorage) GetOrdersCount(status models.OrderStatus) (int, error) {
	var count int64
	tx := s.DB
	if status != models.OrderStatusAny {
		tx = s.DB.Where("status = ?", status)
	}

	if err := tx.Model(&models.Order{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (s *SqliteStorage) UpdateOrder(order *models.Order) error {
	if _, err := s.GetOrderById(order.ID); err != nil {
		return err
	}

	return s.DB.Save(order).Error
}
