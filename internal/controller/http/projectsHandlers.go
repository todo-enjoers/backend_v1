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

func (ctrl *Controller) HandleCreateProject(c echo.Context) error {
	var request model.ProjectRequest

	// Validate user with Token returning userID
	userID, err := ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(errPkg.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: errPkg.ErrValidationToken.Error(),
			},
		)
	}
	ctrl.log.Info("HandleCreateProject: logged in", zap.String("user_id", userID.String()))

	if err = c.Bind(&request); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: errPkg.ErrBindingRequest.Error(),
			},
		)
	}

	project := &model.ProjectDTO{
		ID:        uuid.New(),
		Name:      request.Name,
		CreatedBy: userID,
	}

	err = ctrl.store.Project().Create(c.Request().Context(), project)
	if errors.Is(err, errPkg.ErrAlreadyExists) {
		ctrl.log.Error("project already exists", zap.Error(err))
		return c.JSON(
			http.StatusConflict,
			model.ErrorResponse{
				Error: errPkg.ErrAlreadyExists.Error(),
			},
		)
	}
	response := model.ProjectResponse{
		ID:        project.ID,
		Name:      project.Name,
		CreatedBy: project.CreatedBy,
	}
	ctrl.log.Info("successfully created new project", zap.Any("project", response))
	return c.JSON(http.StatusCreated, response)
}

func (ctrl *Controller) HandleDeleteProject(c echo.Context) error {
	var (
		projectIDStr string
		projectID    uuid.UUID
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
	ctrl.log.Info("HandleDeleteProject: logged in", zap.String("user_id", userID.String()))

	projectIDStr = c.Param("id")
	projectID, err = uuid.Parse(projectIDStr)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: errPkg.ErrBadRequestId.Error(),
			},
		)
	}

	err = ctrl.store.Project().Delete(c.Request().Context(), projectID)
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

	ctrl.log.Info("successfully deleted project", zap.String("id", projectID.String()))
	return c.NoContent(http.StatusNoContent)
}

func (ctrl *Controller) HandleUpdateProject(c echo.Context) error {
	var (
		request      model.ProjectRequest
		projectIDStr string
		projectID    uuid.UUID
		userID       uuid.UUID
		err          error
	)

	// Validate user with Token returning userID
	userID, err = ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: errPkg.Unauthenticated.Error(),
			},
		)
	}
	ctrl.log.Info("HandleUpdateProject: logged in", zap.String("user_id", userID.String()))

	projectIDStr = c.Param("id")
	projectID, err = uuid.Parse(projectIDStr)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: errPkg.ErrBadRequestId.Error(),
			},
		)
	}

	if err = c.Bind(&request); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: errPkg.ErrBadRequestId.Error(),
			},
		)
	}

	gotProject, err := ctrl.store.Project().GetByID(c.Request().Context(), projectID)
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			model.ErrorResponse{
				Error: err.Error(),
			},
		)
	}
	gotProject.Name = request.Name

	err = ctrl.store.Project().UpdateName(c.Request().Context(), request.Name, projectID)
	if err != nil {
		switch {
		case errors.Is(err, errPkg.ErrNotFound):
			{
				return c.JSON(
					http.StatusNotFound,
					model.ErrorResponse{
						Error: errPkg.ErrNotFound.Error(),
					},
				)
			}
		case errors.Is(err, errPkg.ErrAlreadyExists):
			{
				return c.JSON(
					http.StatusConflict,
					model.ErrorResponse{
						Error: errPkg.ErrAlreadyExists.Error(),
					},
				)
			}
		case err != nil:
			return c.JSON(
				http.StatusInternalServerError,
				model.ErrorResponse{
					Error: errPkg.ErrInternalServer.Error(),
				},
			)
		}
	}

	ctrl.log.Info("successfully updated project", zap.Any("project", gotProject))
	return c.JSON(http.StatusCreated, gotProject)
}

func (ctrl *Controller) HandleGetMyProject(c echo.Context) error {
	var myProjects []model.ProjectDTO

	// Validate user with Token returning userID
	userID, err := ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(errPkg.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: errPkg.ErrValidationToken.Error(),
			},
		)
	}
	ctrl.log.Info("HandleGetMyProject: logged in", zap.String("user_id", userID.String()))

	myProjects, err = ctrl.store.Project().GetMyProjects(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, errPkg.ErrNotAccessible) {
			return c.JSON(
				http.StatusNoContent, model.ErrorResponse{
					Error: errPkg.ErrNoContent.Error(),
				},
			)
		}
		return err
	}
	return c.JSON(http.StatusOK, myProjects)
}

func (ctrl *Controller) HandleGetMyProjectById(c echo.Context) error {
	var (
		projectIDStr string
		userID       uuid.UUID
		projectID    uuid.UUID
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
	ctrl.log.Info("HandleGetMyProjectById: logged in", zap.String("user_id", userID.String()))

	projectIDStr = c.Param("id")
	projectID, err = uuid.Parse(projectIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: errPkg.ErrBadRequestId.Error(),
		},
		)
	}

	project, err := ctrl.store.Project().GetByID(c.Request().Context(), projectID)
	if err != nil {
		if errors.Is(err, errPkg.ErrNotAccessible) {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: errPkg.ErrNotFound.Error(),
			},
			)
		}
		return err
	}
	return c.JSON(http.StatusOK, project)

}
