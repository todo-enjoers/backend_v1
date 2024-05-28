package http

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/todo-enjoers/backend_v1/internal/model"
	errPkg "github.com/todo-enjoers/backend_v1/internal/pkg/errors"
	"go.uber.org/zap"
	"net/http"
)

func (ctrl *Controller) HandleCreateColumn(c echo.Context) error {
	var (
		request model.ColumRequest
		userID  uuid.UUID
		err     error
	)

	// Validate user with Token returning userID
	userID, err = ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(errPkg.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: errPkg.ErrValidationToken.Error(),
			},
		)
	}
	ctrl.log.Info("HandleCreateColumn: logged in", zap.String("user_id", userID.String()))

	if err = c.Bind(&request); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: errPkg.ErrBindingRequest.Error(),
			},
		)
	}

	column := &model.ColumDTO{
		ProjectId: request.ProjectId,
		Name:      request.Name,
		Order:     request.Order,
	}

	err = ctrl.store.Column().CreateColumn(c.Request().Context(), column)
	if errors.Is(err, errPkg.ErrAlreadyExists) {
		ctrl.log.Error("user already exists", zap.Error(err))
		return c.JSON(
			http.StatusConflict,
			model.ErrorResponse{
				Error: errPkg.ErrAlreadyExists.Error(),
			},
		)
	}
	response := model.ColumResponse{
		ProjectId: column.ProjectId,
		Name:      column.Name,
		Order:     column.Order,
	}
	ctrl.log.Info("successfully created new todo", zap.Any("todo", response))
	return c.JSON(http.StatusCreated, response)
}

func (ctrl *Controller) HandleDeleteColumn(c echo.Context) error {
	var (
		projectID   string
		projectUUID uuid.UUID
		columnName  string
		userID      uuid.UUID
		err         error
	)

	// Validate user with Token returning userID
	userID, err = ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(errPkg.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: errPkg.ErrValidationToken.Error(),
			},
		)
	}
	ctrl.log.Info("HandleDeleteColumn: logged in", zap.String("user_id", userID.String()))

	projectID = c.Param("id")
	projectUUID, err = uuid.Parse(projectID)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: errPkg.ErrBadRequestId.Error(),
			})
	}
	columnName = c.Param("name")

	err = ctrl.store.Column().DeleteColumn(c.Request().Context(), columnName, projectUUID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.ErrorResponse{
				Error: errPkg.ErrInternalServer.Error(),
			})
	}
	return c.NoContent(http.StatusNoContent)
}

func (ctrl *Controller) HandleGetColumnByName(c echo.Context) error {
	var (
		projectUUID uuid.UUID
		projectID   string
		userID      uuid.UUID
		columnName  string
		err         error
	)

	// Validate user with Token returning userID
	userID, err = ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(errPkg.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: errPkg.ErrValidationToken.Error(),
			},
		)
	}
	ctrl.log.Info("HandleGetColumnByName: logged in", zap.String("user_id", userID.String()))

	projectID = c.Param("id")
	projectUUID, err = uuid.Parse(projectID)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: errPkg.ErrBadRequestId.Error(),
			})
	}
	columnName = c.Param("name")
	column, err := ctrl.store.Column().GetColumnByName(c.Request().Context(), columnName, projectUUID)
	if err != nil {
		if errors.Is(err, errPkg.ErrNotAccessible) {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: errPkg.ErrNotFound.Error(),
			})
		}
		return err
	}
	return c.JSON(http.StatusOK, column)
}

func (ctrl *Controller) HandleUpdateColumn(c echo.Context) error {
	var (
		request      model.ColumRequest
		projectIDStr string
		projectUUID  uuid.UUID
		columnName   string
		userID       uuid.UUID
		err          error
	)

	// Validate user with Token returning userID
	userID, err = ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(errPkg.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: errPkg.ErrValidationToken.Error(),
			},
		)
	}
	ctrl.log.Info("HandleUpdateColumn: logged in", zap.String("user_id", userID.String()))

	projectIDStr = c.Param("id")
	projectUUID, err = uuid.Parse(projectIDStr)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: errPkg.ErrBadRequestId.Error(),
			})
	}

	columnName = c.Param("name")
	if err = c.Bind(&request); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: errPkg.ErrBadRequestId.Error(),
			},
		)
	}
	column, err := ctrl.store.Column().GetColumnByName(c.Request().Context(), columnName, projectUUID)
	if err != nil {
		if errors.Is(err, errPkg.ErrNotFound) {
			return c.JSON(
				http.StatusNotFound,
				model.ErrorResponse{
					Error: errPkg.ErrNotFound.Error(),
				},
			)
		}
		return c.JSON(
			http.StatusInternalServerError,
			model.ErrorResponse{
				Error: errPkg.ErrInternalServer.Error(),
			},
		)
	}

	column.Name = request.Name

	err = ctrl.store.Column().UpdateColumn(c.Request().Context(), column, columnName, projectUUID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.ErrorResponse{
				Error: errPkg.ErrInternalServer.Error(),
			},
		)
	}

	ctrl.log.Info("successfully updated column", zap.Any("todo", column))
	return c.JSON(http.StatusCreated, column)
}

func (ctrl *Controller) HandleGetAllColumn(c echo.Context) error {
	var (
		listColumns  []model.ColumDTO
		userID       uuid.UUID
		projectIDStr string
		projectUUID  uuid.UUID
		err          error
	)

	// Taking a userID from request
	userID, err = ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(errPkg.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: errPkg.ErrValidationToken.Error(),
			},
		)
	}
	ctrl.log.Info("HandleGetAllColumn: logged in", zap.String("user_id", userID.String()))

	projectIDStr = c.Param("id")
	projectUUID, err = uuid.Parse(projectIDStr)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: errPkg.ErrBadRequestId.Error(),
			})
	}

	listColumns, err = ctrl.store.Column().GetAllColumns(c.Request().Context(), projectUUID)
	if err != nil {
		ctrl.log.Error("error while getting group by id from DB", zap.Error(err))
		return c.JSON(
			http.StatusNoContent,
			model.ErrorResponse{
				Error: errPkg.ErrGetByID.Error(),
			},
		)
	}

	return c.JSON(http.StatusOK, listColumns)
}
