package middleware

import (
	"context"

	"github.com/labstack/echo"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
)

func UseUrl(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), "url", c.Request().URL.String())))
		return next(c)
	}
}

func UseCategories(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		categories, err := storage.StorageInstance.GetMainCategories()
		if err != nil {
			return err
		}

		c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), "categories", categories)))
		return next(c)
	}
}
