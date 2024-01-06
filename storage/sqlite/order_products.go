package sqlite

import "github.com/w1png/go-htmx-ecommerce-template/models"

func (s *SqliteStorage) CreateOrderProduct(orderProduct *models.OrderProduct) error {
	return s.DB.Create(orderProduct).Error
}
