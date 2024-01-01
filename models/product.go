package models

import (
	"strings"

	"github.com/w1png/htmx-template/gorm_types"
	"gorm.io/gorm"
)

type StockType int

var STOCK_TYPES_ARRAY = []StockType{
	StockTypeInStock,
	StockTypeOutOfStock,
	StockTypeOrder,
}

const PRODUCTS_PER_PAGE = 20

const (
	StockTypeInStock StockType = iota
	StockTypeOutOfStock
	StockTypeOrder
)

func (s StockType) ToString() string {
	switch s {
	case StockTypeInStock:
		return "В наличии"
	case StockTypeOutOfStock:
		return "Нет в наличии"
	case StockTypeOrder:
		return "Под заказ"
	default:
		return "Неизвестно"
	}
}

type Product struct {
	gorm.Model

	ID            uint
	Slug          string `gorm:"unique"`
	Name          string
	Description   string
	Price         int
	DiscountPrice int
	StockType     StockType
	Tags          string
	CategoryId    uint
	Category      Category               `gorm:"foreignKey:CategoryId"`
	Images        gorm_types.StringArray `gorm:"type:text[]"`
	IsEnabled     bool
	IsFeatured    bool
}

func NewProduct(
	slug string,
	name string,
	description string,
	price int,
	stock_type StockType,
	tags string,
	category_id uint,
	images []string,
) *Product {
	dashed_slug := strings.ReplaceAll(slug, " ", "-")

	return &Product{
		Slug:          dashed_slug,
		Name:          name,
		Description:   description,
		Price:         price,
		DiscountPrice: -1,
		StockType:     stock_type,
		Tags:          tags,
		CategoryId:    category_id,
		Images:        images,
		IsEnabled:     false,
		IsFeatured:    false,
	}
}
