package user_handlers

import (
	"github.com/labstack/echo"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	user_templates "github.com/w1png/go-htmx-ecommerce-template/templates/user"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
)

func GatherIndexHandlers(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	user_page_group.GET("/", IndexHandler)
	user_api_group.GET("/index", IndexApiHandler)
}

func IndexApiHandler(c echo.Context) error {
	featured_products, err := storage.StorageInstance.GetFeaturedProducts()
	if err != nil {
		return err
	}

	return utils.Render(c, user_templates.IndexApi(featured_products))
}

func IndexHandler(c echo.Context) error {
	featured_products, err := storage.StorageInstance.GetFeaturedProducts()
	if err != nil {
		return err
	}

	return utils.Render(c, user_templates.Index(featured_products))
}
