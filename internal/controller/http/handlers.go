package http

import (
	"github.com/labstack/echo/v4"
	"github.com/todo-enjoers/backend_v1/internal/model"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func (ctrl *Controller) HandleRegisterUser(c echo.Context) error {
	var req model.UserDTO
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err) //change
	}
	if ok, err := ValidateRequest(req); !ok {
		return err
	}

	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	req.Password = string(HashedPassword)

	err := ctrl.storage.InsertUser(c.Request().Context(), req)
	if err != nil {
		log.Printf("got unexpected error: %v\r\n", err)
		return c.String(http.StatusNotFound, "not found")
	}

	return c.JSON(http.StatusOK, result)
}

func (ctrl *Controller) HandleLoginUser(c echo.Context) error {
	// Handler: login
	return nil
}

/*
func (ctrl *Controller) HandleGetMe(c echo.Context) error {
	result, err := ctrl.storage.GetMe(c.Request().Context())
	if err != nil {
		log.Printf("got unexpected error: %v\r\n", err)
		return c.String(http.StatusBadRequest, "unexpected error")
	}
	return c.JSON(http.StatusOK, result)
}
*/

func (ctrl *Controller) HandleChangePasswordUser(c echo.Context) error {
	return nil
}

//func (ctrl *Controller) HandleGetAll(c echo.Context) error {
//	result, err := ctrl.storage.GetAll(c.Request().Context())
//	if err != nil {
//		log.Printf("got unexpected error: %v\r\n", err)
//		return c.String(http.StatusBadRequest, "unexpected error")
//	}
//	return c.JSON(http.StatusOK, result)
//}

func (ctrl *Controller) HandleRefreshToken(c echo.Context) error {
	return nil
}
