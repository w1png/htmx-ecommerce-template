package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model

	ID        uint
	Name      string
	Slug      string `gorm:"unique"`
	ImagePath string
	Tags      string

	Parent   *Category
	ParentId uint

	IsEnabled bool
}

const CATEGORIES_PER_PAGE = 20

func (c *Category) BeforeDelete(tx *gorm.DB) error {
	return tx.Model(&Product{}).Where("category_id = ?", c.ID).Update("category_id", 0).Error
}

func (c Category) AfterFind(tx *gorm.DB) error {
	if c.ParentId == 0 {
		return nil
	}
	return tx.Model(&Category{}).Where("id = ?", c.ParentId).First(&c.Parent).Error
}

func NewCategory(name, slug, image_path, tags string, parent_id uint) *Category {
	return &Category{
		Name:      name,
		Slug:      slug,
		ImagePath: image_path,
		Tags:      tags,
		ParentId:  parent_id,
		IsEnabled: true,
	}
}
