package user_handlers

import (
	"reflect"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/w1png/go-htmx-ecommerce-template/errors"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	user_templates "github.com/w1png/go-htmx-ecommerce-template/templates/user"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
)

func GatherProductsRoutes(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	user_page_group.GET("/products/:slug", ProductHandler)
	user_api_group.GET("/products/:slug", ProductApiHandler)
}

func ProductHandler(c echo.Context) error {
	product, err := storage.StorageInstance.GetProductBySlug(c.Param("slug"))
	if err != nil {
		return err
	}

	cart_product, err := storage.StorageInstance.GetCartProductByProductIdAndCartID(product.ID, utils.GetCartFromContext(c.Request().Context()).ID)
	if err != nil && reflect.TypeOf(err) != reflect.TypeOf(&errors.ObjectNotFoundError{}) {
		log.Error(err)
		return err
	}

	if cart_product == nil {
		cart_product = &models.CartProduct{
			Quantity: 0,
		}
	}

	return utils.Render(c, user_templates.Product(product, cart_product))
}

func ProductApiHandler(c echo.Context) error {
	product, err := storage.StorageInstance.GetProductBySlug(c.Param("slug"))
	if err != nil {
		return err
	}

	cart_product, err := storage.StorageInstance.GetCartProductByProductIdAndCartID(product.ID, utils.GetCartFromContext(c.Request().Context()).ID)
	if err != nil && reflect.TypeOf(err) != reflect.TypeOf(&errors.ObjectNotFoundError{}) {
		log.Error(err)
		return err
	}

	if cart_product == nil {
		cart_product = &models.CartProduct{
			Quantity: 0,
		}
	}

	return utils.Render(c, user_templates.ProductApi(product, cart_product))
}
