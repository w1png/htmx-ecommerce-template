package admin_handlers

import (
	"fmt"
	"math"
	"net/http"
	"reflect"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/w1png/htmx-template/errors"
	"github.com/w1png/htmx-template/models"
	"github.com/w1png/htmx-template/storage"
	admin_users_templates "github.com/w1png/htmx-template/templates/admin/users"
	"github.com/w1png/htmx-template/utils"
)

func getNextUsersPage(page int) (int, error) {
	users_count, err := storage.StorageInstance.GetUsersCount()
	if err != nil {
		return -1, err
	}

	total_pages := int(math.Ceil(float64(users_count) / float64(models.USERS_PER_PAGE)))
	if total_pages <= page {
		return -1, nil
	}

	return page + 1, nil
}

func UserIndexHandler(c echo.Context) error {
	users, err := storage.StorageInstance.GetUsers(utils.GetOffsetAndLimit(1, models.USERS_PER_PAGE))
	if err != nil {
		return err
	}

	next_page, err := getNextUsersPage(1)
	if err != nil {
		return err
	}

	return utils.Render(c, admin_users_templates.Index(users, next_page))
}

func UserIndexApiHandler(c echo.Context) error {
	users, err := storage.StorageInstance.GetAllUsers()
	if err != nil {
		return err
	}

	next_page, err := getNextUsersPage(1)
	if err != nil {
		return err
	}

	return utils.Render(c, admin_users_templates.IndexApi(users, next_page))
}

func GetUserHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	user, err := storage.StorageInstance.GetUserById(uint(id))
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&errors.ObjectNotFoundError{}) {
			return c.String(http.StatusNotFound, "Пользователь не найден")
		}
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	return utils.Render(c, admin_users_templates.User(user))
}

func PostUserHandler(c echo.Context) error {
	c.Response().Header().Set("HX-Reswap", "innerHTML")

	if err := c.Request().ParseForm(); err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	username := c.FormValue("username")
	password := c.FormValue("password")
	password_repeat := c.FormValue("password_repeat")
	is_admin := c.FormValue("is_admin") == "true"

	if username == "" || password == "" || password_repeat == "" {
		return c.String(http.StatusBadRequest, "Поля не могут быть пустыми")
	}

	if !models.IsUsernameValid(username) {
		return c.String(http.StatusBadRequest, models.GetUsernameRules())
	}

	if password != password_repeat {
		return c.String(http.StatusBadRequest, "Пароли не совпадают")
	}

	if !models.IsPasswordValid(password) {
		return c.String(http.StatusBadRequest, models.GetPasswordRules())
	}

	if _, err := storage.StorageInstance.GetUserByUsername(username); err == nil {
		return c.String(http.StatusBadRequest, "Пользователь с таким именем уже существует")
	} else {
		if err != nil && reflect.TypeOf(err) != reflect.TypeOf(&errors.ObjectNotFoundError{}) {
			log.Error(err)
			return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
	}

	user, err := models.NewUser(username, password, is_admin)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}
	if err := storage.StorageInstance.CreateUser(user); err != nil {
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	c.Response().Header().Del("HX-Reswap")

	return utils.Render(c, admin_users_templates.User(user))
}

func PutUserHandler(c echo.Context) error {
	c.Response().Header().Set("HX-Reswap", "innerHTML")

	if err := c.Request().ParseForm(); err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	username := c.FormValue("username")
	password := c.FormValue("password")
	password_repeat := c.FormValue("password_repeat")
	is_admin := c.FormValue("is_admin") == "true"

	if username == "" {
		return c.String(http.StatusBadRequest, "Имя пользователя не может быть пустым")
	}

	if !models.IsUsernameValid(username) {
		return c.String(http.StatusBadRequest, models.GetUsernameRules())
	}

	user, err := storage.StorageInstance.GetUserById(uint(id))
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&errors.ObjectNotFoundError{}) {
			return c.String(http.StatusNotFound, "Пользователь не найден")
		}
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	user.Username = username

	if password != "" {
		if password != password_repeat {
			return c.String(http.StatusBadRequest, "Пароли не совпадают")
		}

		if !models.IsPasswordValid(password) {
			return c.String(http.StatusBadRequest, models.GetPasswordRules())
		}

		user.PasswordHash, err = user.HashPassword(password)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
	}
	if user.ID != 1 {
		user.IsAdmin = is_admin
	}

	if err := storage.StorageInstance.UpdateUser(user); err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	c.Response().Header().Del("HX-Reswap")
	c.Response().Header().Set("HX-Trigger", fmt.Sprintf("user_saved_%d", user.ID))

	return utils.Render(c, admin_users_templates.AddUserForm())
}

func EditUserHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	user, err := storage.StorageInstance.GetUserById(uint(id))
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&errors.ObjectNotFoundError{}) {
			return c.String(http.StatusNotFound, "Пользователь не найден")
		}
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	return utils.Render(c, admin_users_templates.UserEdit(user))
}

func GetAddUserHandler(c echo.Context) error {
	return utils.Render(c, admin_users_templates.AddUserForm())
}

func DeleteUserHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	if err := storage.StorageInstance.DeleteUserById(uint(id)); err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&errors.ObjectNotFoundError{}) {
			return c.String(http.StatusNotFound, "Пользователь не найден")
		}
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	return c.HTMLBlob(http.StatusOK, []byte(""))
}

func SearchUsersHandler(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return err
	}

	username := c.FormValue("search_username")

	var users []*models.User
	var err error
	if username == "" {
		if users, err = storage.StorageInstance.GetAllUsers(); err != nil {
			return err
		}
	} else {
		if users, err = storage.StorageInstance.GetAllUsersByUsernameFuzzy(username); err != nil {
			return err
		}
	}

	return utils.Render(c, admin_users_templates.Users(users[:models.USERS_PER_PAGE], -1))
}

func GetUsersPage(c echo.Context) error {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		return err
	}

	users, err := storage.StorageInstance.GetUsers(utils.GetOffsetAndLimit(page, models.USERS_PER_PAGE))
	if err != nil {
		return err
	}

	next_page, err := getNextUsersPage(page)
	if err != nil {
		return err
	}

	return utils.Render(c, admin_users_templates.Users(users, next_page))
}
