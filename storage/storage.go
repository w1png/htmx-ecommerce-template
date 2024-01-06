package storage

import (
	"github.com/w1png/go-htmx-ecommerce-template/config"
	"github.com/w1png/go-htmx-ecommerce-template/errors"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage/sqlite"
)

type Storage interface {
	CreateUser(user *models.User) error
	GetUserById(id uint) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	GetUsers(offset, limit int) ([]*models.User, error)
	GetAllUsersByUsernameFuzzy(username string) ([]*models.User, error)
	GetUsersByUsernameFuzzy(username string, offset, limit int) ([]*models.User, error)
	GetUsersCount() (int, error)
	UpdateUser(user *models.User) error
	DeleteUserById(id uint) error

	CreateCategory(category *models.Category) error
	GetCategoryById(id uint) (*models.Category, error)
	GetCategoryBySlug(slug string) (*models.Category, error)
	GetAllCategories() ([]*models.Category, error)
	GetMainCategories() ([]*models.Category, error)
	GetCategoryChildren(id uint) ([]*models.Category, error)
	GetCategoryProducts(id uint) ([]*models.Product, error)
	GetCategories(offset, limit int) ([]*models.Category, error)
	GetCategoriesCount() (int, error)
	GetAllCategoriesByNameFuzzy(name string) ([]*models.Category, error)
	GetCategoriesByNameFuzzy(name string, offset, limit int) ([]*models.Category, error)
	UpdateCategory(category *models.Category) error
	DeleteCategoryById(id uint) error

	CreateProduct(product *models.Product) error
	GetProductById(id uint) (*models.Product, error)
	GetProductBySlug(slug string) (*models.Product, error)
	GetAllProducts() ([]*models.Product, error)
	GetEnabledProducts(offset, limit int) ([]*models.Product, error)
	GetFeaturedProducts() ([]*models.Product, error)
	GetProducts(offset, limit int) ([]*models.Product, error)
	GetProductsCount() (int, error)
	GetAllProductsByNameFuzzy(name string) ([]*models.Product, error)
	GetProductsByNameFuzzy(name string, offset, limit int) ([]*models.Product, error)
	GetProductsByTags(search_tag string, offset, limit int) ([]*models.Product, error)
	UpdateProduct(product *models.Product) error
	DeleteProductById(id uint) error

	CreateCart(cart *models.Cart) error
	GetCartById(id uint) (*models.Cart, error)
	GetCartByUUID(uuid string) (*models.Cart, error)

	CreateCartProduct(cartProduct *models.CartProduct) error
	GetCartProductById(id uint) (*models.CartProduct, error)
	GetCartProductByProductIdAndCartID(id uint, cart_id uint) (*models.CartProduct, error)
	UpdateCartProduct(cartProduct *models.CartProduct) error
	DeleteCartProductById(id uint) error

	CreateOrder(order *models.Order) error
	GetOrderById(id uint) (*models.Order, error)
	GetOrders(status models.OrderStatus, offset, limit int) ([]*models.Order, error)
	GetOrdersCount(status models.OrderStatus) (int, error)
	UpdateOrder(order *models.Order) error

	CreateOrderProduct(order_product *models.OrderProduct) error
}

var StorageInstance Storage

func InitStorage() error {
	var err error
	switch config.ConfigInstance.StorageType {
	case "sqlite":
		if StorageInstance, err = sqlite.NewSQLiteStorage(); err != nil {
			return err
		}
	default:
		return errors.NewUnknownDatabaseTypeError(config.ConfigInstance.StorageType)
	}

	return nil
}
