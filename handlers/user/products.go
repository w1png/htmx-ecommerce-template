package user_handlers

import (
	"github.com/labstack/echo"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	user_templates "github.com/w1png/go-htmx-ecommerce-template/templates/user"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
)

func ProductHandler(c echo.Context) error {
	product, err := storage.StorageInstance.GetProductBySlug(c.Param("slug"))
	if err != nil {
		return err
	}

	return utils.Render(c, user_templates.Product(product))
}

func ProductApiHandler(c echo.Context) error {
	product, err := storage.StorageInstance.GetProductBySlug(c.Param("slug"))
	if err != nil {
		return err
	}

	return utils.Render(c, user_templates.ProductApi(product))
}
