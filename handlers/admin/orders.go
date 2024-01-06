package admin_handlers

import (
	"fmt"
	"math"
	"net/http"
	"reflect"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/w1png/go-htmx-ecommerce-template/errors"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	admin_orders_templates "github.com/w1png/go-htmx-ecommerce-template/templates/admin/orders"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
)

func GatherOrdersRoutes(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	admin_page_group.GET("/orders", OrdersIndexHandler)
	admin_api_group.GET("/orders", OrdersIndexApiHandler)

	admin_api_group.GET("/orders", OrdersIndexApiHandler)
	admin_api_group.GET("/orders/:id/modal", GetOrderModalHandler)
	admin_api_group.PUT("/orders/:id", UpdateOrderStatusHandler)
	admin_api_group.GET("/orders/:id/status", GetOrderStatusHandler)
	admin_api_group.GET("/orders/page/:page", GetOrdersPageHandler)
}

func getStatusFromQueryParam(c echo.Context) (models.OrderStatus, error) {
	statusRaw := c.QueryParam("status")
	if statusRaw == "" {
		return models.OrderStatusAny, nil
	}
	var status models.OrderStatus
	statusInt, err := strconv.Atoi(statusRaw)
	if err != nil {
		return status, c.String(http.StatusBadRequest, "Неверный запрос")
	}
	return models.OrderStatus(statusInt), nil
}

func OrdersIndexHandler(c echo.Context) error {
	status, err := getStatusFromQueryParam(c)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	orders, err := storage.StorageInstance.GetOrders(status, 0, models.ORDERS_PER_PAGE)
	if err != nil {
		return err
	}

	return utils.Render(c, admin_orders_templates.Index(orders, status))
}

func OrdersIndexApiHandler(c echo.Context) error {
	status, err := getStatusFromQueryParam(c)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	orders, err := storage.StorageInstance.GetOrders(status, 0, models.ORDERS_PER_PAGE)
	if err != nil {
		return err
	}

	return utils.Render(c, admin_orders_templates.IndexApi(orders, status))
}

func GetOrdersPageHandler(c echo.Context) error {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	status, err := getStatusFromQueryParam(c)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	offset, limit := utils.GetOffsetAndLimit(page, models.ORDERS_PER_PAGE)
	orders, err := storage.StorageInstance.GetOrders(status, offset, limit)
	if err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	count, err := storage.StorageInstance.GetOrdersCount(status)
	if err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	next_page := page + 1
	total_pages := int(math.Ceil(float64(count) / float64(limit)))
	if total_pages <= page {
		next_page = -1
	}

	return utils.Render(c, admin_orders_templates.Orders(orders, next_page, status))
}

func GetOrderModalHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	order, err := storage.StorageInstance.GetOrderById(uint(id))
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&errors.ObjectNotFoundError{}) {
			return c.String(http.StatusNotFound, "Заказ не найден")
		}
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	return utils.Render(c, admin_orders_templates.OrderModal(order))
}

func UpdateOrderStatusHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	statusRaw, err := strconv.Atoi(c.QueryParam("status"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	status := models.OrderStatus(statusRaw)
	if status == models.OrderStatusAny {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	order, err := storage.StorageInstance.GetOrderById(uint(id))
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&errors.ObjectNotFoundError{}) {
			return c.String(http.StatusNotFound, "Заказ не найден")
		}
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	order.Status = status

	if err := storage.StorageInstance.UpdateOrder(order); err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	c.Response().Header().Set("HX-Trigger", fmt.Sprintf("update_status_%d", order.ID))
	return c.NoContent(http.StatusOK)
}

func GetOrderStatusHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	order, err := storage.StorageInstance.GetOrderById(uint(id))
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&errors.ObjectNotFoundError{}) {
			return c.String(http.StatusNotFound, "Заказ не найден")
		}
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	return utils.Render(c, admin_orders_templates.OrderStatusDropdown(order))
}
