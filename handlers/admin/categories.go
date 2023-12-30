package admin_handlers

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/w1png/htmx-template/errors"
	"github.com/w1png/htmx-template/models"
	"github.com/w1png/htmx-template/storage"
	categories_admin_templates "github.com/w1png/htmx-template/templates/admin/categories"
	"github.com/w1png/htmx-template/utils"
)

func CategoriesIndexHandler(c echo.Context) error {
	categories, err := storage.StorageInstance.GetCategories(utils.GetOffsetAndLimit(1, models.CATEGORIES_PER_PAGE))
	if err != nil {
		return err
	}

	next_page, err := utils.GetNextPage(1, storage.StorageInstance.GetCategoriesCount, models.CATEGORIES_PER_PAGE)
	if err != nil {
		return err
	}

	return utils.Render(c, categories_admin_templates.Index(categories, next_page))
}

func GetCategoryHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	category, err := storage.StorageInstance.GetCategoryById(uint(id))
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&errors.ObjectNotFoundError{}) {
			return c.String(http.StatusNotFound, "Пользователь не найден")
		}
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	return utils.Render(c, categories_admin_templates.Category(category))
}

func GetAddCategoryHandler(c echo.Context) error {
	categories, err := storage.StorageInstance.GetAllCategories()
	if err != nil {
		return err
	}

	return utils.Render(c, categories_admin_templates.AddCategoryForm(categories))
}

func EditCategoryHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	category, err := storage.StorageInstance.GetCategoryById(uint(id))
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&errors.ObjectNotFoundError{}) {
			return c.String(http.StatusNotFound, "Пользователь не найден")
		}
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	categories, err := storage.StorageInstance.GetAllCategories()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	return utils.Render(c, categories_admin_templates.CategoryEdit(category, categories))
}

func PostCategoryHandler(c echo.Context) error {
	c.Response().Header().Set("HX-Reswap", "innerHTML")

	if err := c.Request().ParseMultipartForm(30 * 1024 * 1024); err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	name := c.FormValue("name")
	slug := c.FormValue("slug")
	tags := c.FormValue("tags")
	parent, err := strconv.ParseUint(c.FormValue("parent"), 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	if name == "" || slug == "" || tags == "" {
		return c.String(http.StatusBadRequest, "Поля обязательны для заполнения")
	}

	image, err := c.FormFile("image")
	if err != nil {
		return c.String(http.StatusBadRequest, "Фотография обязательна")
	}
	image_filename, err := utils.SaveImage(image)
	if err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Ошибка сохранения фотографии")
	}

	if parent != 0 {
		if _, err := storage.StorageInstance.GetCategoryById(uint(parent)); err != nil {
			if reflect.TypeOf(err) == reflect.TypeOf(&errors.ObjectNotFoundError{}) {
				return c.String(http.StatusBadRequest, "Родительская категория не найдена")
			}

			return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
	}

	category := models.NewCategory(
		name,
		slug,
		image_filename,
		tags,
		uint(parent),
	)

	if err := storage.StorageInstance.CreateCategory(category); err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&errors.SlugNotUniqueError{}) {
			return c.String(http.StatusBadRequest, "Категория с такой ссылкой уже существует")
		}
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	c.Response().Header().Del("HX-Reswap")

	return utils.Render(c, categories_admin_templates.Category(category))
}

func DeleteCategoryHandler(c echo.Context) error {

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	if err := storage.StorageInstance.DeleteCategoryById(uint(id)); err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&errors.ObjectNotFoundError{}) {
			return c.String(http.StatusNotFound, "Категория не найдена")
		}
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	return c.NoContent(http.StatusOK)
}

func PutCategoryHandler(c echo.Context) error {
	c.Response().Header().Set("HX-Reswap", "innerHTML")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	var image_filename string
	if err := c.Request().ParseMultipartForm(30 * 1024 * 1024); err == nil {
		image, err := c.FormFile("image")
		if err == nil {
			image_filename, _ = utils.SaveImage(image)
		}
	} else {
		if err := c.Request().ParseForm(); err != nil {
			return c.String(http.StatusBadRequest, "Неверный запрос")
		}
	}

	name := c.FormValue("name")
	slug := c.FormValue("slug")
	tags := c.FormValue("tags")
	parent, err := strconv.ParseUint(c.FormValue("parent"), 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	if name == "" || slug == "" || tags == "" {
		return c.String(http.StatusBadRequest, "Поля обязательны для заполнения")
	}

	category, err := storage.StorageInstance.GetCategoryById(uint(id))
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&errors.ObjectNotFoundError{}) {
			return c.String(http.StatusNotFound, "Категория не найдена")
		}
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	if parent != 0 {
		if _, err := storage.StorageInstance.GetCategoryById(uint(parent)); err != nil {
			if reflect.TypeOf(err) == reflect.TypeOf(&errors.ObjectNotFoundError{}) {
				return c.String(http.StatusBadRequest, "Родительская категория не найдена")
			}

			return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
	}

	category.Name = name
	category.Slug = slug
	category.Tags = tags
	category.ParentId = uint(parent)

	if image_filename != "" {
		category.ImagePath = image_filename
	}

	if err := storage.StorageInstance.UpdateCategory(category); err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&errors.SlugNotUniqueError{}) {
			return c.String(http.StatusBadRequest, "Категория с такой ссылкой уже существует")
		}
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	categories, err := storage.StorageInstance.GetAllCategories()
	if err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	c.Response().Header().Del("HX-Reswap")
	c.Response().Header().Set("HX-Trigger", fmt.Sprintf("category_saved_%d", category.ID))
	return utils.Render(c, categories_admin_templates.AddCategoryForm(categories))
}

func GetCategoriesPage(c echo.Context) error {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		log.Error(err)
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	categories, err := storage.StorageInstance.GetCategories(utils.GetOffsetAndLimit(page, models.CATEGORIES_PER_PAGE))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	next_page, err := utils.GetNextPage(page, storage.StorageInstance.GetCategoriesCount, models.CATEGORIES_PER_PAGE)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	return utils.Render(c, categories_admin_templates.Categories(categories, next_page))
}

func SearchCategoriesHandler(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	name := c.FormValue("search_name")
	var categories []*models.Category
	var err error
	if name != "" {
		categories, err = storage.StorageInstance.GetCategoriesByNameFuzzy(name, 0, models.CATEGORIES_PER_PAGE)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}

		return utils.Render(c, categories_admin_templates.Categories(categories, -1))

	}
	categories, err = storage.StorageInstance.GetCategories(utils.GetOffsetAndLimit(1, models.CATEGORIES_PER_PAGE))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	return utils.Render(c, categories_admin_templates.Categories(categories, 2))
}
