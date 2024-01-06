package admin_handlers

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
	admin_products_templates "github.com/w1png/go-htmx-ecommerce-template/templates/admin/products"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
)

func GatherProductsRoutes(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	admin_page_group.GET("/products", ProductsIndexHandler)
	admin_api_group.GET("/products", ProductsIndexApiHandler)

	admin_api_group.GET("/products", ProductsIndexApiHandler)
	admin_api_group.POST("/products", PostProductHandler)
	admin_api_group.DELETE("/products/:id", DeleteProductHandler)
	admin_api_group.GET("/products/add", GetAddProductFormHandler)
	admin_api_group.GET("/products/:id", GetProductHandler)
	admin_api_group.GET("/products/page/:page", GetProductsPage)
	admin_api_group.GET("/products/:id/edit", GetEditProductFormHandler)
	admin_api_group.PUT("/products/:id", PutProductHandler)
}

func ProductsIndexHandler(c echo.Context) error {
	products, err := storage.StorageInstance.GetProducts(utils.GetOffsetAndLimit(1, models.PRODUCTS_PER_PAGE))
	if err != nil {
		return err
	}

	categories, err := storage.StorageInstance.GetAllCategories()
	if err != nil {
		return err
	}

	next_page, err := utils.GetNextPage(1, storage.StorageInstance.GetProductsCount, models.PRODUCTS_PER_PAGE)
	if err != nil {
		return err
	}

	return utils.Render(c, admin_products_templates.Index(products, categories, next_page))
}

func ProductsIndexApiHandler(c echo.Context) error {
	products, err := storage.StorageInstance.GetProducts(utils.GetOffsetAndLimit(1, models.PRODUCTS_PER_PAGE))
	if err != nil {
		return err
	}

	categories, err := storage.StorageInstance.GetAllCategories()
	if err != nil {
		return err
	}

	next_page, err := utils.GetNextPage(1, storage.StorageInstance.GetProductsCount, models.PRODUCTS_PER_PAGE)
	if err != nil {
		return err
	}

	return utils.Render(c, admin_products_templates.IndexApi(products, categories, next_page))
}

func PostProductHandler(c echo.Context) error {
	c.Response().Header().Set("HX-Reswap", "innerHTML")

	if err := c.Request().ParseMultipartForm(30 * 1024 * 1024); err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	name := c.FormValue("name")
	slug := c.FormValue("slug")
	tags := c.FormValue("tags")
	price, err := strconv.Atoi(c.FormValue("price"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}
	description := c.FormValue("description")
	category_id, err := strconv.ParseUint(c.FormValue("category"), 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}
	if category_id != 0 {
		if _, err := storage.StorageInstance.GetCategoryById(uint(category_id)); err != nil {
			if reflect.TypeOf(err) == reflect.TypeOf(&errors.ObjectNotFoundError{}) {
				return c.String(http.StatusBadRequest, "Категория не найдена")
			}
			log.Error(err)
			return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
	}

	stock_type_int, err := strconv.Atoi(c.FormValue("stock_type"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	stock_type := models.StockType(stock_type_int)

	multipart_form, err := c.MultipartForm()
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	var image_filenames []string
	images := multipart_form.File["images"]
	if len(images) == 0 {
		return c.String(http.StatusBadRequest, "Фотографии обязательны для заполнения")
	}
	for _, image := range images {
		filename, err := utils.SaveImage(image)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Ошибка сохранения файла \"%s\"", image.Filename))
		}
		image_filenames = append(image_filenames, filename)
	}

	product := models.NewProduct(
		slug,
		name,
		description,
		price,
		stock_type,
		tags,
		uint(category_id),
		image_filenames,
	)

	if err := storage.StorageInstance.CreateProduct(product); err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&errors.SlugNotUniqueError{}) {
			return c.String(http.StatusBadRequest, "Товар с такой ссылкой уже существует")
		}
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	c.Response().Header().Del("HX-Reswap")
	c.Response().Header().Set("HX-Trigger", "save_product")

	return utils.Render(c, admin_products_templates.Product(product))
}

func DeleteProductHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	if err := storage.StorageInstance.DeleteProductById(uint(id)); err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&errors.ObjectNotFoundError{}) {
			return c.String(http.StatusBadRequest, "Товар не найден")
		}

		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	return c.NoContent(http.StatusOK)
}

func GetEditProductFormHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}
	product, err := storage.StorageInstance.GetProductById(uint(id))
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&errors.ObjectNotFoundError{}) {
			return c.String(http.StatusNotFound, "Товар не найден")
		}
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	categories, err := storage.StorageInstance.GetAllCategories()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	return utils.Render(c, admin_products_templates.EditProductForm(product, categories, models.STOCK_TYPES_ARRAY))
}

func GetAddProductFormHandler(c echo.Context) error {
	categories, err := storage.StorageInstance.GetAllCategories()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	return utils.Render(c, admin_products_templates.AddProductForm(categories, models.STOCK_TYPES_ARRAY))
}

func PutProductHandler(c echo.Context) error {
	c.Response().Header().Set("HX-Reswap", "innerHTML")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	if err := c.Request().ParseMultipartForm(30 * 1024 * 1024); err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	name := c.FormValue("name")
	slug := c.FormValue("slug")
	tags := c.FormValue("tags")
	is_enabled := c.FormValue("is_enabled") == "true"
	is_featured := c.FormValue("is_featured") == "true"
	price, err := strconv.Atoi(c.FormValue("price"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}
	discount_price, err := strconv.Atoi(c.FormValue("discount_price"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}
	description := c.FormValue("description")
	category_id, err := strconv.ParseUint(c.FormValue("category"), 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}
	if category_id != 0 {
		if _, err := storage.StorageInstance.GetCategoryById(uint(category_id)); err != nil {
			if reflect.TypeOf(err) == reflect.TypeOf(&errors.ObjectNotFoundError{}) {
				return c.String(http.StatusBadRequest, "Категория не найдена")
			}
			log.Error(err)
			return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
	}

	stock_type_int, err := strconv.Atoi(c.FormValue("stock_type"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	stock_type := models.StockType(stock_type_int)

	multipart_form, err := c.MultipartForm()
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	var image_filenames []string
	images := multipart_form.File["images"]
	for _, image := range images {
		filename, err := utils.SaveImage(image)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Ошибка сохранения файла \"%s\"", image.Filename))
		}
		image_filenames = append(image_filenames, filename)
	}

	product, err := storage.StorageInstance.GetProductById(uint(id))
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&errors.ObjectNotFoundError{}) {
			return c.String(http.StatusNotFound, "Товар не найден")
		}
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	product.Slug = slug
	product.IsEnabled = is_enabled
	product.IsFeatured = is_featured
	product.Name = name
	product.CategoryId = uint(category_id)
	product.Description = description
	product.Price = price
	product.DiscountPrice = discount_price
	product.StockType = stock_type
	product.Tags = tags
	if len(image_filenames) != 0 {
		product.Images = image_filenames
	}

	if err := storage.StorageInstance.UpdateProduct(product); err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&errors.SlugNotUniqueError{}) {
			return c.String(http.StatusBadRequest, "Товар с такой ссылкой уже существует")
		}
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	c.Response().Header().Del("HX-Reswap")
	c.Response().Header().Set("HX-Trigger", fmt.Sprintf("product_saved_%d", product.ID))

	categories, err := storage.StorageInstance.GetAllCategories()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	return utils.Render(c, admin_products_templates.AddProductForm(categories, models.STOCK_TYPES_ARRAY))
}

func GetProductHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	product, err := storage.StorageInstance.GetProductById(uint(id))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	return utils.Render(c, admin_products_templates.Product(product))
}

func GetProductsPage(c echo.Context) error {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	products, err := storage.StorageInstance.GetProducts(utils.GetOffsetAndLimit(page, models.PRODUCTS_PER_PAGE))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	next_page, err := utils.GetNextPage(page, storage.StorageInstance.GetProductsCount, models.PRODUCTS_PER_PAGE)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	return utils.Render(c, admin_products_templates.Products(products, next_page))
}
