package admin_handlers

import (
	"github.com/labstack/echo"
	admin_templates "github.com/w1png/htmx-template/templates/admin"
	"github.com/w1png/htmx-template/utils"
)

func AdminIndexHandler(c echo.Context) error {
	return utils.Render(c, admin_templates.Index(c.Request().Context()))
}

func AdminApiIndexHandler(c echo.Context) error {
	return utils.Render(c, admin_templates.IndexApiNavbar())
}
