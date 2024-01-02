package user_handlers

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/w1png/go-htmx-ecommerce-template/errors"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	"github.com/w1png/go-htmx-ecommerce-template/templates/components"
	user_templates "github.com/w1png/go-htmx-ecommerce-template/templates/user"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
)

func GetCartHandler(c echo.Context) error {
	return utils.Render(c, user_templates.CartProducts(utils.GetCartFromContext(c.Request().Context()).Products))
}

func ChangeCartProductQuantityHandler(c echo.Context) error {
	should_decrease := c.QueryParam("decrease") == "true"

	product_id, err := strconv.ParseUint(c.Param("product_id"), 10, 64)
	if err != nil {
		return err
	}

	cart := utils.GetCartFromContext(c.Request().Context())

	cart_product, err := storage.StorageInstance.GetCartProductByProductIdAndCartID(uint(product_id), cart.ID)
	if err != nil && reflect.TypeOf(err) != reflect.TypeOf(&errors.ObjectNotFoundError{}) {
		log.Error(err)
		return err
	}

	if cart_product == nil {
		product, err := storage.StorageInstance.GetProductById(uint(product_id))
		if err != nil {
			if reflect.TypeOf(err) == reflect.TypeOf(&errors.ObjectNotFoundError{}) {
				return c.String(http.StatusNotFound, "Товар не найден")
			}
			log.Error(err)
			return err
		}

		cart_product = models.NewCartProduct(
			product.ID,
			cart.ID,
			product.Slug,
			product.Name,
			product.Price,
			product.DiscountPrice,
			0,
		)
	}

	fmt.Println(1)
	fmt.Printf("Cart product: %+v\n", cart_product)

	if should_decrease && cart_product.Quantity > 0 {
		cart_product.Quantity--
	} else {
		cart_product.Quantity++
	}

	if cart_product.ID == 0 {
		if err := storage.StorageInstance.CreateCartProduct(cart_product); err != nil {
			return err
		}
	} else {
		if err := storage.StorageInstance.UpdateCartProduct(cart_product); err != nil {
			return err
		}
	}

	fmt.Println(2)

	return utils.Render(c, components.AddToCartButton(cart_product.Product.ID, cart_product.Quantity))
}
