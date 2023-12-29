package user_handlers

import "github.com/labstack/echo"

func HealthHandler(c echo.Context) error {
	return c.String(200, "OK")
}
