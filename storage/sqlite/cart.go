package sqlite

import "github.com/w1png/go-htmx-ecommerce-template/models"

func (s *SqliteStorage) CreateCart(cart *models.Cart) error {
	return s.DB.Create(cart).Error
}

func (s *SqliteStorage) GetCartByUUID(uuid string) (*models.Cart, error) {
	var cart models.Cart
	if err := s.DB.Where("uuid = ?", uuid).First(&cart).Error; err != nil {
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
