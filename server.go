package main

import (
	"fmt"

	"github.com/labstack/echo"
	echoMiddleware "github.com/labstack/echo/middleware"
	"github.com/w1png/htmx-template/config"
	admin_handlers "github.com/w1png/htmx-template/handlers/admin"
	user_handlers "github.com/w1png/htmx-template/handlers/user"
	"github.com/w1png/htmx-template/middleware"
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
}

func (s *HTTPServer) gatherUserApiRoutes() {
	api_group := s.echo.Group("/api")
	api_group.GET("/health", user_handlers.HealthHandler)
	api_group.GET("/index", user_handlers.IndexHandler)
	api_group.GET("/admin_login", user_handlers.LoginPageApiHandler)

	api_group.POST("/admin_login", user_handlers.PostLoginHandler)
}

func (s *HTTPServer) gatherAdminRoutes(g *echo.Group) {
	g.GET("/health", user_handlers.HealthHandler)
	g.GET("", admin_handlers.AdminIndexHandler)

	g.GET("/users", admin_handlers.UserIndexHandler)
	g.GET("/categories", admin_handlers.CategoriesIndexHandler)
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

	api_group.POST("/categories", admin_handlers.PostCategoryHandler)
}

func (s *HTTPServer) Run() error {
	return s.echo.Start(fmt.Sprintf(":%s", config.ConfigInstance.Port))
}
