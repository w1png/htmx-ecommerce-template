package main

import (
	"fmt"

	"github.com/labstack/echo"
	echoMiddleware "github.com/labstack/echo/middleware"
	"github.com/w1png/go-htmx-ecommerce-template/config"
	admin_handlers "github.com/w1png/go-htmx-ecommerce-template/handlers/admin"
	user_handlers "github.com/w1png/go-htmx-ecommerce-template/handlers/user"
	"github.com/w1png/go-htmx-ecommerce-template/middleware"
)

type HTTPServer struct {
	echo *echo.Echo
}

func NewHTTPServer() *HTTPServer {
	server := &HTTPServer{
		echo: echo.New(),
	}

	server.echo.Use(echoMiddleware.Logger())
	server.echo.Use(echoMiddleware.Recover())
	server.echo.Use(middleware.UseAuth)
	server.echo.Use(middleware.UseUrl)
	server.echo.Use(middleware.UseCategories)
	server.echo.Use(middleware.UseCart)

	server.echo.Static("/static", "static")

	server.gatherUserApiRoutes()
	server.gatherUserRoutes()

	admin_group := server.echo.Group("/admin")
	admin_group.Use(middleware.UseAdmin)

	server.gatherAdminApiRoutes(admin_group)
	server.gatherAdminRoutes(admin_group)

	return server
}

func (s *HTTPServer) gatherUserRoutes() {
	s.echo.GET("/health", user_handlers.HealthHandler)

	s.echo.GET("/", user_handlers.IndexHandler)
	s.echo.GET("/admin_login", user_handlers.LoginPageHandler)
	s.echo.GET("/categories/:slug", user_handlers.CategoryHandler)
	s.echo.GET("/products/:slug", user_handlers.ProductHandler)
	s.echo.GET("/checkout", user_handlers.CheckoutHandler)
}

func (s *HTTPServer) gatherUserApiRoutes() {
	api_group := s.echo.Group("/api")
	api_group.GET("/health", user_handlers.HealthHandler)
	api_group.GET("/index", user_handlers.IndexHandler)
	api_group.GET("/admin_login", user_handlers.LoginPageApiHandler)
	api_group.GET("/categories/:slug", user_handlers.CategoryApiHandler)
	api_group.GET("/products/:slug", user_handlers.ProductApiHandler)

	api_group.POST("/admin_login", user_handlers.PostLoginHandler)

	api_group.GET("/checkout", user_handlers.CheckoutApiHandler)
	api_group.GET("/cart", user_handlers.GetCartHandler)
	api_group.PUT("/cart/change_quantity/:product_id", user_handlers.ChangeCartProductQuantityHandler)

	api_group.GET("/checkout/delivery_type_form", user_handlers.GetDeliveryTypeForm)

	api_group.POST("/checkout", user_handlers.PostOrderHandler)
}

func (s *HTTPServer) gatherAdminRoutes(g *echo.Group) {
	g.GET("/health", user_handlers.HealthHandler)
	g.GET("", admin_handlers.AdminIndexHandler)

	g.GET("/users", admin_handlers.UserIndexHandler)
	g.GET("/categories", admin_handlers.CategoriesIndexHandler)
	g.GET("/products", admin_handlers.ProductsIndexHandler)
	g.GET("/orders", admin_handlers.OrdersIndexHandler)
}

func (s *HTTPServer) gatherAdminApiRoutes(g *echo.Group) {
	api_group := g.Group("/api")
	api_group.GET("/index", admin_handlers.AdminApiIndexHandler)
	api_group.GET("/health", user_handlers.HealthHandler)

	api_group.GET("/users", admin_handlers.UserIndexApiHandler)
	api_group.GET("/users/:id", admin_handlers.GetUserHandler)
	api_group.POST("/users", admin_handlers.PostUserHandler)
	api_group.GET("/users/:id/edit", admin_handlers.EditUserHandler)
	api_group.PUT("/users/:id", admin_handlers.PutUserHandler)
	api_group.GET("/users/add", admin_handlers.GetAddUserHandler)
	api_group.POST("/users/search", admin_handlers.SearchUsersHandler)
	api_group.DELETE("/users/:id", admin_handlers.DeleteUserHandler)
	api_group.GET("/users/page/:page", admin_handlers.GetUsersPage)

	api_group.GET("/categories", admin_handlers.CategoriesIndexHandler)
	api_group.GET("/categories/:id", admin_handlers.GetCategoryHandler)
	api_group.GET("/categories/:id/edit", admin_handlers.EditCategoryHandler)
	api_group.GET("/categories/add", admin_handlers.GetAddCategoryHandler)
	api_group.DELETE("/categories/:id", admin_handlers.DeleteCategoryHandler)
	api_group.PUT("/categories/:id", admin_handlers.PutCategoryHandler)
	api_group.GET("/categories/page/:page", admin_handlers.GetCategoriesPage)
	api_group.POST("/categories/search", admin_handlers.SearchCategoriesHandler)

	api_group.GET("/products", admin_handlers.ProductsIndexApiHandler)
	api_group.POST("/products", admin_handlers.PostProductHandler)
	api_group.DELETE("/products/:id", admin_handlers.DeleteProductHandler)
	api_group.GET("/products/add", admin_handlers.GetAddProductFormHandler)
	api_group.GET("/products/:id", admin_handlers.GetProductHandler)
	api_group.GET("/products/page/:page", admin_handlers.GetProductsPage)
	api_group.GET("/products/:id/edit", admin_handlers.GetEditProductFormHandler)
	api_group.PUT("/products/:id", admin_handlers.PutProductHandler)

	api_group.POST("/categories", admin_handlers.PostCategoryHandler)

	api_group.GET("/orders", admin_handlers.OrdersIndexApiHandler)
	api_group.GET("/orders/:id/modal", admin_handlers.GetOrderModalHandler)
	api_group.PUT("/orders/:id", admin_handlers.UpdateOrderStatusHandler)
	api_group.GET("/orders/:id/status", admin_handlers.GetOrderStatusHandler)
	api_group.GET("/orders/page/:page", admin_handlers.GetOrdersPageHandler)
}

func (s *HTTPServer) Run() error {
	return s.echo.Start(fmt.Sprintf(":%s", config.ConfigInstance.Port))
}
