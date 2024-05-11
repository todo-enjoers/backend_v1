package http

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/todo-enjoers/backend_v1/internal/config"
	"github.com/todo-enjoers/backend_v1/internal/controller"
	"github.com/todo-enjoers/backend_v1/internal/storage"
	"go.uber.org/zap"
)

// Checking whether the interface "Controller" implements the structure "Controller"
var _ controller.Controller = (*Controller)(nil)

type Controller struct {
	server  *echo.Echo
	storage storage.Interface
	cfg     *config.Config
	log     *zap.Logger
}

func New(store storage.Interface, cfg *config.Config, log *zap.Logger) (*Controller, error) {
	log.Info("init controller")
	ctrl := &Controller{
		server:  echo.New(),
		storage: store,
		cfg:     cfg,
		log:     log,
	}
	if err := ctrl.configure(); err != nil {
		return nil, err
	}
	return ctrl, nil
}

func (ctrl *Controller) configure() error {
	ctrl.configureMiddlewares()
	ctrl.configureRoutes()
	return nil
}

func (ctrl *Controller) configureRoutes() {
	log.Info("configuring routes")
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

func (ctrl *Controller) configureMiddlewares(c echo.Context) error {
	var middlewares = []echo.MiddlewareFunc{
		middleware.Gzip(),
		middleware.Recover(),
		middleware.RequestIDWithConfig(middleware.RequestIDConfig{
			Skipper:      middleware.DefaultSkipper,
			Generator:    uuid.NewString,
			TargetHeader: echo.HeaderXRequestID,
		}),
		middleware.Logger(),
	}
	ctrl.server.Use(middlewares...)
}

func (ctrl *Controller) Run(_ context.Context, log *zap.Logger) error {
	log.Info("starting HTTP server on address: %s", zap.String("url", ctrl.cfg.BindAddr))
	return ctrl.server.Start(ctrl.cfg.BindAddr)
}

func (ctrl *Controller) Shutdown(ctx context.Context) error {
	return ctrl.server.Shutdown(ctx)
}
