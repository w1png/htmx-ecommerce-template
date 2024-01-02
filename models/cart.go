package models

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/labstack/echo"
	"github.com/w1png/go-htmx-ecommerce-template/config"
	"gorm.io/gorm"
)

func GenerateUUID(c echo.Context) string {
	hash_input := fmt.Sprintf("%d%s%s", time.Now().UnixNano(), c.RealIP(), config.ConfigInstance.JWTSecret)
	return fmt.Sprintf("%x", sha256.Sum256([]byte(hash_input)))
}

type Cart struct {
	gorm.Model

	ID       uint
	UUID     string
	Products []*CartProduct
}

func (c *Cart) AfterFind(tx *gorm.DB) error {
	if err := tx.Model(&CartProduct{}).Where("cart_id = ?", c.ID).Find(&c.Products).Error; err != nil {
		return err
	}
	return nil
}

func NewCart(uuid string) *Cart {
	return &Cart{UUID: uuid}
}

func (c *Cart) GetTotalPrice() int {
	total := 0
	for _, product := range c.Products {
		p := product.Price
		if product.DiscountPrice != -1 {
			p = product.DiscountPrice
		}
		total += p * product.Quantity
	}
	return total
}
