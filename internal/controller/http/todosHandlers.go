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

func (ctrl *Controller) HandleCreateTodo(c echo.Context) error {
	var (
		request model.TodoCreateRequest
		userID  uuid.UUID
	)

	// Validate user with Token returning userID
	userId, err := ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(errPkg.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: errPkg.ErrValidationToken.Error(),
			},
		)
	}
	ctrl.log.Info("HandleCreateTodo: logged in", zap.String("user_id", userID.String()))

	if err = c.Bind(&request); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: errPkg.ErrBindingRequest.Error(),
			},
		)
	}

	todo := &model.TodoDTO{
		ID:          uuid.New(),
		Name:        request.Name,
		Description: request.Description,
		IsCompleted: request.IsCompleted,
		ProjectID:   request.ProjectID,
		CreatedBy:   userId,
		Column:      request.Column,
	}

	err = ctrl.store.Todo().Create(c.Request().Context(), todo)
	if errors.Is(err, errPkg.ErrAlreadyExists) {
		ctrl.log.Error("project already exists", zap.Error(err))
		return c.JSON(
			http.StatusConflict,
			model.ErrorResponse{
				Error: errPkg.ErrAlreadyExists.Error(),
			},
		)
	}
	response := model.TodoCreateResponse{
		ID:          todo.ID,
		Name:        todo.Name,
		Description: todo.Description,
		IsCompleted: todo.IsCompleted,
		ProjectID:   todo.ProjectID,
		CreatedBy:   todo.CreatedBy,
		Column:      todo.Column,
	}
	ctrl.log.Info("successfully created new todo", zap.Any("todo", response))
	return c.JSON(http.StatusCreated, response)

}

func (ctrl *Controller) HandleGetTodosById(c echo.Context) error {
	var (
		id     string
		todoID uuid.UUID
		userID uuid.UUID
		err    error
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
	ctrl.log.Info("HandleGetTodosById: logged in", zap.String("user_id", userID.String()))

	id = c.Param("id")
	todoID, err = uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: errPkg.ErrBadRequestId.Error(),
		},
		)
	}

	todo, err := ctrl.store.Todo().GetByID(c.Request().Context(), todoID)
	if err != nil {
		if errors.Is(err, errPkg.ErrNotAccessible) {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: errPkg.ErrNotFound.Error(),
			},
			)
		}
		return err
	}
	return c.JSON(http.StatusOK, todo)
}

func (ctrl *Controller) HandleChangeTodo(c echo.Context) error {
	var (
		request   model.TodoUpdateRequest
		todoIDStr string
		todoID    uuid.UUID
		userID    uuid.UUID
		err       error
	)

	// Validate user with Token returning userID
	userID, err = ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: "Unauthorized",
			},
		)
	}
	ctrl.log.Info("HandleChangeTodo: logged in", zap.String("user_id", userID.String()))

	// parse id
	todoIDStr = c.Param("id")
	todoID, err = uuid.Parse(todoIDStr)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: errPkg.ErrBadRequestId.Error(),
			},
		)
	}

	//get request
	if err := c.Bind(&request); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: errPkg.ErrBadRequestId.Error(),
			},
		)
	}

	todo, err := ctrl.store.Todo().GetByID(c.Request().Context(), todoID)
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

	//change todos
	todo.Name = request.Name
	todo.Description = request.Description
	todo.IsCompleted = request.IsCompleted

	//work with db
	err = ctrl.store.Todo().Update(c.Request().Context(), todo, todoID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.ErrorResponse{
				Error: errPkg.ErrInternalServer.Error(),
			},
		)
	}

	ctrl.log.Info("successfully updated todo", zap.Any("todo", todo))
	return c.JSON(http.StatusCreated, todo)
}

func (ctrl *Controller) HandleDeleteTodo(c echo.Context) error {
	var (
		todoIDStr string
		todoID    uuid.UUID
		userID    uuid.UUID
		err       error
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

	todoIDStr = c.Param("id")
	todoID, err = uuid.Parse(todoIDStr)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: errPkg.ErrBadRequestId.Error(),
			},
		)
	}

	err = ctrl.store.Todo().Delete(c.Request().Context(), todoID)
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
				Error: err.Error(),
			},
		)
	}

	ctrl.log.Info("successfully deleted todo", zap.String("id", todoID.String()))
	return c.NoContent(http.StatusNoContent)
}

func (ctrl *Controller) HandleGetAllTodos(c echo.Context) error {
	var (
		listTodos []model.TodoDTO
		err       error
		userID    uuid.UUID
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
	ctrl.log.Info("HandleGetAllTodos: logged in", zap.String("user_id", userID.String()))

	listTodos, err = ctrl.store.Todo().GetAll(c.Request().Context(), userID)
	if err != nil {
		ctrl.log.Error("error while getting todos by id from DB", zap.Error(err))
		return c.JSON(
			http.StatusNoContent,
			model.ErrorResponse{
				Error: errPkg.ErrGetByID.Error(),
			},
		)
	}

	return c.JSON(http.StatusOK, listTodos)
}
