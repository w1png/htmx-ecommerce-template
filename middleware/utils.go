package middleware

import (
	"context"

	"github.com/labstack/echo"
)

func UseUrl(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), "url", c.Request().URL.String())))
		return next(c)
	}
}
