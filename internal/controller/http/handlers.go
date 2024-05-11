package http

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/todo-enjoers/backend_v1/internal/model"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func (ctrl *Controller) HandleRegisterUser(c echo.Context) error {
	var request model.UserRegisterRequest
	if err := c.Bind(&request); err != nil {
		ctrl.log.Error("could not bind request", zap.Error(err))
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()}) //change
	}
	if ok, err := ValidateRequest(request); !ok {
		return err
	}

	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &model.UserDTO{
		ID:       uuid.New(),
		Login:    request.Login,
		Password: string(HashedPassword),
	}
	ctrl.log.Info("got user", zap.Any("user", user))

	err := ctrl.store.InsertUser(c.Request().Context(), request)
	if err != nil {
		log.Printf("got unexpected error: %v\r\n", err)
		return c.String(http.StatusNotFound, "not found")
	}

	return c.JSON(http.StatusOK, nil)
}

func (ctrl *Controller) HandleLoginUser(c echo.Context) error {
	// Handler: login
	return nil
}

func (ctrl *Controller) HandleGetMe(c echo.Context) error {
	result, err := ctrl.storage.GetMe(c.Request().Context())
	if err != nil {
		log.Printf("got unexpected error: %v\r\n", err)
		return c.String(http.StatusBadRequest, "unexpected error")
	}
	return c.JSON(http.StatusOK, result)
}

func (ctrl *Controller) HandleChangePasswordUser(c echo.Context) error {
	return nil
}

func (ctrl *Controller) HandleGetAll(c echo.Context) error {
	result, err := ctrl.storage.GetAll(c.Request().Context())
	if err != nil {
		log.Printf("got unexpected error: %v\r\n", err)
		return c.String(http.StatusBadRequest, "unexpected error")
	}
	return c.JSON(http.StatusOK, result)
}

func (ctrl *Controller) HandleRefreshToken(c echo.Context) error {
	return nil
}

func (ctrl *Controller) HandleGetTodos(c echo.Context) error {}
