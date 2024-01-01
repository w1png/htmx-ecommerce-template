package utils

import (
	"context"

	"github.com/labstack/echo"
	"github.com/w1png/go-htmx-ecommerce-template/models"
)

type ResponseData struct {
	User *models.User
	Data interface{}
}

func MarshalResponse(c echo.Context, data interface{}) *ResponseData {
	var user *models.User
	userAny := c.Request().Context().Value("user")
	if userAny == nil {
		user = nil
	} else {
		user = userAny.(*models.User)
	}

	return &ResponseData{
		User: user,
		Data: data,
	}
}

func GetUserFromContext(ctx context.Context) *models.User {
	var user *models.User
	userAny := ctx.Value("user")
	if userAny == nil {
		user = nil
	} else {
		user = userAny.(*models.User)
	}

	return user
}

func GetUrlFromContext(ctx context.Context) string {
	var url string
	urlAny := ctx.Value("url")
	if urlAny == nil {
		url = ""
	} else {
		url = urlAny.(string)
	}

	return url
}
