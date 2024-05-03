package http

import (
	"github.com/labstack/echo/v4"
	"github.com/todo-enjoers/backend_v1/internal/storage/model"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func (ctrl *Controller) HandleRegisterUser(c echo.Context) error {
	Login := c.FormValue("email")
	Password := c.FormValue("password")

	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	user := &model.UserDTO{
		Login:    Login,
		Password: string(HashedPassword),
	}
	return c.JSON(http.StatusOK)
}

func (ctrl *Controller) HandleGetAll(c echo.Context) error {
	result, err := ctrl.storage.GetAll(c.Request().Context())
	if err != nil {
		log.Printf("got unexpected error: %v\r\n", err)
		return c.String(http.StatusBadRequest, "unexpected error")
	}
	return c.JSON(http.StatusOK, result)
}
