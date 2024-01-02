package middleware

import (
	"context"
	"net/http"
	"reflect"

	"github.com/labstack/echo"
	"github.com/w1png/go-htmx-ecommerce-template/errors"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
)

func UseCart(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cart_uuid, err := c.Cookie("cart_uuid")
		if err != nil && err != http.ErrNoCookie {
			return err
		}

		var cart *models.Cart
		if cart_uuid != nil {
			cart, err = storage.StorageInstance.GetCartByUUID(cart_uuid.Value)
			if err != nil && reflect.TypeOf(err) != reflect.TypeOf(&errors.ObjectNotFoundError{}) {
				return err
			}
		}

		if cart == nil {
			cart = models.NewCart(models.GenerateUUID(c))
			if err := storage.StorageInstance.CreateCart(cart); err != nil {
				return err
			}

			c.SetCookie(&http.Cookie{
				Name:  "cart_uuid",
				Path:  "/",
				Value: cart.UUID,
			})
		}

		c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), "cart", cart)))
		return next(c)
	}
}
