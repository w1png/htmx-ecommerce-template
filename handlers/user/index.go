package user_handlers

import (
	"github.com/labstack/echo"
	user_templates "github.com/w1png/htmx-template/templates/user"
	"github.com/w1png/htmx-template/utils"
)

func IndexApiHandler(c echo.Context) error {
	return utils.Render(c, user_templates.IndexApi())
}

func IndexHandler(c echo.Context) error {
	return utils.Render(c, user_templates.Index())
}
