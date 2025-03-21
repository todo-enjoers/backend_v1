package http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"

	"github.com/todo-enjoers/backend_v1/internal/config"
	"github.com/todo-enjoers/backend_v1/internal/controller"
	"github.com/todo-enjoers/backend_v1/internal/pkg/token"
	"github.com/todo-enjoers/backend_v1/internal/storage"
)

// Checking whether the interface "Controller" implements the structure "Controller"
var _ controller.Controller = (*Controller)(nil)

type Controller struct {
	server *echo.Echo
	log    *zap.Logger
	cfg    *config.Config
	token  token.Provider
	store  storage.Interface
}

func New(
	store storage.Interface,
	log *zap.Logger,
	cfg *config.Config,
	tokenProvider token.Provider,
) (*Controller, error) {
	log.Info("initialize controller")
	ctrl := &Controller{
		server: echo.New(),
		store:  store,
		cfg:    cfg,
		log:    log,
		token:  tokenProvider,
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
	api := ctrl.server.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("/register", ctrl.HandleRegister)
			users.POST("/login", ctrl.HandleLogin)
			users.GET("/me", ctrl.HandleGetMe)
			users.GET("/all", ctrl.HandleGetAll)
			users.POST("/change-password", ctrl.HandleChangePassword)
			users.POST("/refresh-token", ctrl.HandleRefreshToken)
		}

		todos := api.Group("/todos")
		{
			todos.GET("/", ctrl.HandleGetAllTodos)
			todos.GET("/:id", ctrl.HandleGetTodosById)
			todos.POST("/", ctrl.HandleCreateTodo)
			todos.PUT("/:id", ctrl.HandleChangeTodo)
			todos.DELETE("/:id", ctrl.HandleDeleteTodo)
		}

		projects := api.Group("/projects")
		{
			projects.POST("/create", ctrl.HandleCreateProject)
			projects.DELETE("/delete/:id", ctrl.HandleDeleteProject)
			projects.PUT("/update/:id", ctrl.HandleUpdateProject)
			projects.GET("/", ctrl.HandleGetMyProject)
			projects.GET("/:id", ctrl.HandleGetMyProjectById)
		}
		columns := api.Group("/columns")
		{
			columns.POST("/", ctrl.HandleCreateColumn)
			columns.DELETE("/:id/:name", ctrl.HandleDeleteColumn)
			columns.PUT("/:id/:name", ctrl.HandleUpdateColumn)
			columns.GET("/:id/:name", ctrl.HandleGetColumnByName)
			columns.GET("/:id/", ctrl.HandleGetAllColumn)
		}
	}
}

func (ctrl *Controller) configureMiddlewares() {
	var middlewares = []echo.MiddlewareFunc{
		middleware.Gzip(),
		middleware.Recover(),
		middleware.RequestIDWithConfig(middleware.RequestIDConfig{
			Skipper:      middleware.DefaultSkipper,
			Generator:    uuid.NewString,
			TargetHeader: echo.HeaderXRequestID,
		}),
		middleware.Logger(),
		middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogValuesFunc: ctrl.logValuesFunc,
			LogLatency:    true,
			LogRequestID:  true,
			LogMethod:     true,
			LogURI:        true,
		}),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		}),
	}
	ctrl.server.Use(middlewares...)
}

func (ctrl *Controller) logValuesFunc(_ echo.Context, v middleware.RequestLoggerValues) error {
	ctrl.log.Info("Request",
		zap.String("uri", v.URI),
		zap.String("method", v.Method),
		zap.Duration("duration", v.Latency),
		zap.String("request-id", v.RequestID),
	)
	return nil
}

func (ctrl *Controller) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	//  goroutine of starting HTTP server
	go func() {
		ctrl.log.Info("starting HTTP server on address", zap.String("", ctrl.cfg.Controller.GetBindAddress()))
		err := ctrl.server.Start(ctrl.cfg.Controller.GetBindAddress())
		if err != nil {
			cancel()
		}
	}()
	return ctx.Err()
}

func (ctrl *Controller) Shutdown(ctx context.Context) error {
	return ctrl.server.Shutdown(ctx)
}
