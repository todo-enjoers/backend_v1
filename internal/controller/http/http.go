package http

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/todo-enjoers/backend_v1/config"
	"github.com/todo-enjoers/backend_v1/internal/controller"
	"github.com/todo-enjoers/backend_v1/internal/storage"
	"log"
)

var _ controller.Controller = (*Controller)(nil)

type Controller struct {
	server  *echo.Echo
	storage storage.Interface
	cfg     *config.Config
}

func New(repo storage.Interface, cfg *config.Config) *Controller {
	log.Println("init controller")
	ctrl := &Controller{
		server:  echo.New(),
		storage: repo,
		cfg:     cfg,
	}
	ctrl.configureRoutes()
	return ctrl
}

func (ctrl *Controller) configure() {
	//ctrl.configureMiddlewares()
	ctrl.configureRoutes()

}

func (ctrl *Controller) configureRoutes() {
	log.Println("configuring routes")
	router := ctrl.server
	api := router.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("/register", ctrl.HandleRegisterUser)
			users.POST("/login", ctrl.HandleLoginUser)
			//users.GET("/me", ctrl.HandleGetMe)
			users.POST("/change-password", ctrl.HandleChangePasswordUser)
			users.POST("/refresh-token", ctrl.HandleRefreshToken)
		}
		todos := api.Group("/todos")
		{
			todos.GET("/", ctrl.HandleGetTodos)
			todos.GET("/{id}", ctrl.HandleGetTodoByID)
			todos.POST("/", ctrl.HandleCreateTodo)
			todos.PUT("/{id}", ctrl.HandleUpdateTodo)
			todos.DELETE("/{id}", ctrl.HandleDeleteTodo)
		}
	}
}

/*
func (ctrl *Controller) configureMiddlewares(c echo.Context) error {

}
*/

func (ctrl *Controller) Run(ctx context.Context) error {
	log.Printf("starting HTTP server on address: %s", ctrl.cfg.BindAddr)
	return ctrl.server.Start(ctrl.cfg.BindAddr)
}

func (ctrl *Controller) Shutdown(ctx context.Context) error {
	return ctrl.server.Shutdown(ctx)
}
