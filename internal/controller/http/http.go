package http

import (
	"github.com/labstack/echo/v4"
	"github.com/todo-enjoers/backend_v1/config"
	"github.com/todo-enjoers/backend_v1/internal/storage"
	"log"
)

type Controller struct {
	echo    *echo.Echo
	storage *storage.Storage
	cfg     *config.Config
}

func New(repo *storage.Storage, cfg *config.Config) *Controller {
	log.Println("init controller")
	ctrl := &Controller{
		echo:    echo.New(),
		storage: repo,
		cfg:     cfg,
	}
	ctrl.configureRoutes()
	return ctrl
}

func (ctrl *Controller) configureRoutes() {
	log.Println("configuring routes")
	router := ctrl.echo
	/*
		router.GET("/todos", ctrl.HandleGetTodos)
		router.GET("/todos/{id}", ctrl.HandleGetTodoByID)
		router.POST("/todos", ctrl.HandleCreateTodo)
		router.PUT("/todos/{id}", ctrl.HandleUpdateTodo)
		router.DELETE("/todos/{id}", ctrl.HandleDeleteTodo)
	*/
	router.POST("/users/register", ctrl.HandleRegisterUser)
	router.POST("/users/login", ctrl.HandleLoginUser)
	router.GET("/users/me", ctrl.HandleGetMe)
	router.POST("/users/change-password", ctrl.HandleChangePasswordUser)
	router.POST("/users/refresh-token", ctrl.HandleRefreshToken)
}

func (ctrl *Controller) Run() error {
	log.Printf("starting http server on address: %s", ctrl.cfg.BindAddr)
	return ctrl.echo.Start(ctrl.cfg.BindAddr)
}
