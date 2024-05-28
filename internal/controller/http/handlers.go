package http

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/todo-enjoers/backend_v1/internal/controller"
	"github.com/todo-enjoers/backend_v1/internal/model"
	"github.com/todo-enjoers/backend_v1/internal/storage"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// ./api/users
func (ctrl *Controller) HandleRegister(c echo.Context) error {
	var request model.UserRegisterRequest

	// Binding request
	if err := c.Bind(&request); err != nil {
		ctrl.log.Error("could not bind request", zap.Error(err))
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: controller.ErrBindingRequest.Error(),
			},
		)
	}

	// Validate request
	if ok, err := request.Validate(); !ok {
		ctrl.log.Error("invalid request", zap.Error(err))
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: storage.ErrBadRegisterRequest.Error(),
			},
		)
	}

	// Hashing password from request
	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.ErrorResponse{
				Error: storage.ErrHashingPassword.Error(),
			},
		)
	}

	// Taking user from UserDTO with new data
	user := &model.UserDTO{
		ID:       uuid.New(),
		Login:    request.Login,
		Password: string(HashedPassword),
	}
	ctrl.log.Info("got user", zap.Any("user", user))

	// Inserting in DB the user
	err = ctrl.store.User().Create(c.Request().Context(), user)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			ctrl.log.Error("user already exists", zap.Error(err))
			return c.JSON(
				http.StatusConflict,
				model.ErrorResponse{
					Error: storage.ErrAlreadyExists.Error(),
				},
			)
		}
		ctrl.log.Error("got error while creating user", zap.Error(err))
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: storage.ErrCreateUser.Error(),
			},
		)
	}
	ctrl.log.Info("successfully created user")

	// Generating token's for the user
	accessToken, refreshToken, err := ctrl.generateAccessAndRefreshTokenForUser(user.ID)
	if err != nil {
		ctrl.log.Error("got error while creating tokens", zap.Error(err))
		return c.JSON(
			http.StatusInternalServerError,
			model.ErrorResponse{
				Error: storage.ErrCreateToken.Error(),
			},
		)
	}

	response := model.UserRegisterResponse{
		ID:           user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return c.JSON(http.StatusCreated, response)
}

func (ctrl *Controller) HandleLogin(c echo.Context) error {
	var request model.UserLoginRequest

	// Binding request
	if err := c.Bind(&request); err != nil {
		ctrl.log.Error("error while binding request", zap.Error(err))
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
			},
		)
	}

	// Getting the "User" from DB
	user, err := ctrl.store.User().GetByLogin(c.Request().Context(), request.Login)
	if err != nil {
		ctrl.log.Error("error while getting user by login from DB", zap.Error(err))
		return c.JSON(
			http.StatusNoContent,
			model.ErrorResponse{
				Error: storage.ErrGetByLogin.Error(),
			},
		)
	}

	// Compare hashed password from request and from DB
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		ctrl.log.Error("invalid password", zap.Error(controller.InvalidPassword))
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: storage.ErrComparingPasswords.Error(),
			},
		)
	}

	// Generating access, refresh tokens for logged user
	access, refresh, err := ctrl.generateAccessAndRefreshTokenForUser(user.ID)
	if err != nil {
		ctrl.log.Error("error while creating tokens", zap.Error(err))
		return c.JSON(
			http.StatusInternalServerError,
			model.ErrorResponse{
				Error: storage.ErrCreateToken.Error(),
			},
		)
	}

	response := &model.UserLoginResponse{
		ID:           user.ID,
		AccessToken:  access,
		RefreshToken: refresh,
	}
	return c.JSON(http.StatusOK, response)
}

func (ctrl *Controller) HandleChangePassword(c echo.Context) error {
	var request model.UserChangePasswordRequest

	// Validate user with Token returning userID
	userID, err := ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
			},
		)
	}
	ctrl.log.Info("HandleChangePassword : logged in", zap.String("user_id", userID.String()))

	// Binding request
	if err = c.Bind(&request); err != nil {
		ctrl.log.Error("error while binding request", zap.Error(err))
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: controller.ErrBindingRequest.Error(),
			},
		)
	}

	// Getting the "User" from DB
	user, err := ctrl.store.User().GetByID(c.Request().Context(), userID)
	if err != nil {
		ctrl.log.Error("error while getting user by login from DB", zap.Error(err))
		return c.JSON(
			http.StatusNoContent,
			model.ErrorResponse{
				Error: storage.ErrGetByLogin.Error(),
			},
		)
	}

	// Compare hashed password from request and from DB
	err = ctrl.CompareHashes([]byte(user.Password), []byte(request.OldPassword))
	if err != nil {
		ctrl.log.Error("invalid password", zap.Error(controller.InvalidPassword))
		return c.JSON(
			http.StatusInternalServerError,
			model.ErrorResponse{
				Error: controller.InvalidPassword.Error(),
			},
		)
	}

	// Compare NewPassword and  NewPasswordAgain
	if request.NewPassword != request.NewPasswordAgain {
		ctrl.log.Error("password are not equal", zap.Error(controller.ErrPasswordAreNotEqual))
		return c.JSON(
			http.StatusNotAcceptable,
			model.ErrorResponse{
				Error: controller.ErrPasswordAreNotEqual.Error(),
			},
		)
	}

	// Hashing NewPassword from request
	newHashedPassword, err := ctrl.PasswordToHash(request.NewPassword)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.ErrorResponse{
				Error: storage.ErrHashingPassword.Error(),
			},
		)
	}

	// Inserting NewPassword in DB
	err = ctrl.store.User().ChangePassword(c.Request().Context(), string(newHashedPassword), user.ID) // ???
	if err != nil {
		ctrl.log.Error("error while inserting in DB changed password", zap.Error(err))
		return c.JSON(
			http.StatusConflict,
			model.ErrorResponse{
				Error: controller.ErrInsertingInDB.Error(),
			},
		)
	}

	return c.NoContent(http.StatusOK)

}

func (ctrl *Controller) HandleGetMe(c echo.Context) error {
	var (
		me     *model.UserDTO
		err    error
		userID uuid.UUID
	)

	// Validate user with Token returning userID
	userID, err = ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
			},
		)
	}
	ctrl.log.Info("HandleGetMe: logged in", zap.String("user_id", userID.String()))

	// Getting "Me" from DB
	me, err = ctrl.store.User().GetByID(c.Request().Context(), userID)
	if err != nil {
		ctrl.log.Error("error while getting user by id from DB", zap.Error(err))
		return c.JSON(
			http.StatusNoContent,
			model.ErrorResponse{
				Error: storage.ErrGetByID.Error(),
			},
		)
	}

	response := &model.UserGetMeResponse{
		ID:   me.ID,
		Name: me.Login,
	}

	return c.JSON(http.StatusOK, response)
}

func (ctrl *Controller) HandleGetAll(c echo.Context) error {
	var list []model.UserDTO

	// Validate user with Token returning userID
	userID, err := ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
			},
		)
	}
	ctrl.log.Info("HandleGetAll: logged in", zap.String("user_id", userID.String()))

	list, err = ctrl.store.User().GetAll(c.Request().Context())
	if err != nil {
		ctrl.log.Error("error while getting users by id from DB", zap.Error(err))
		return c.JSON(
			http.StatusNoContent,
			model.ErrorResponse{
				Error: storage.ErrGetByID.Error(),
			},
		)
	}

	return c.JSON(http.StatusOK, list)
}

func (ctrl *Controller) HandleRefreshToken(c echo.Context) error {
	var (
		request model.UserCoupleTokensRequest
		userID  uuid.UUID
		err     error
	)

	// Validate user with Token returning userID
	userID, err = ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
			},
		)
	}
	ctrl.log.Info("HandleRefreshToken: logged in", zap.String("user_id", userID.String()))

	// Binding request
	if err = c.Bind(&request); err != nil {
		ctrl.log.Error("could not bind request", zap.Error(err))
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: controller.ErrBindingRequest.Error(),
			},
		)
	}

	// Generating token's for the user
	accessToken, refreshToken, err := ctrl.generateAccessAndRefreshTokenForUser(userID)
	if err != nil {
		ctrl.log.Error("got error while creating tokens", zap.Error(err))
		return c.JSON(
			http.StatusInternalServerError,
			model.ErrorResponse{
				Error: storage.ErrCreateToken.Error(),
			},
		)
	}

	response := &model.UserCoupleTokensResponse{
		ID:           userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return c.JSON(http.StatusCreated, response)
}

// ./api/projects

func (ctrl *Controller) HandleCreateProject(c echo.Context) error {
	var request model.ProjectRequest

	// Validate user with Token returning userID
	userID, err := ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
			},
		)
	}
	ctrl.log.Info("HandleCreateProject: logged in", zap.String("user_id", userID.String()))

	if err = c.Bind(&request); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: controller.ErrBindingRequest.Error(),
			},
		)
	}

	project := &model.ProjectDTO{
		ID:        uuid.New(),
		Name:      request.Name,
		CreatedBy: userID,
	}

	err = ctrl.store.Project().Create(c.Request().Context(), project)
	if errors.Is(err, storage.ErrAlreadyExists) {
		ctrl.log.Error("project already exists", zap.Error(err))
		return c.JSON(
			http.StatusConflict,
			model.ErrorResponse{
				Error: storage.ErrAlreadyExists.Error(),
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
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
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
				Error: storage.ErrBadRequestId.Error(),
			},
		)
	}

	err = ctrl.store.Project().Delete(c.Request().Context(), projectID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return c.JSON(
				http.StatusNotFound,
				model.ErrorResponse{
					Error: storage.ErrNotFound.Error(),
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
				Error: controller.Unauthenticated.Error(),
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
				Error: storage.ErrBadRequestId.Error(),
			},
		)
	}

	if err = c.Bind(&request); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: storage.ErrBadRequestId.Error(),
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
		case errors.Is(err, storage.ErrNotFound):
			{
				return c.JSON(
					http.StatusNotFound,
					model.ErrorResponse{
						Error: storage.ErrNotFound.Error(),
					},
				)
			}
		case errors.Is(err, storage.ErrAlreadyExists):
			{
				return c.JSON(
					http.StatusConflict,
					model.ErrorResponse{
						Error: storage.ErrAlreadyExists.Error(),
					},
				)
			}
		case err != nil:
			return c.JSON(
				http.StatusInternalServerError,
				model.ErrorResponse{
					Error: storage.ErrInternalServer.Error(),
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
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
			},
		)
	}
	ctrl.log.Info("HandleGetMyProject: logged in", zap.String("user_id", userID.String()))

	myProjects, err = ctrl.store.Project().GetMyProjects(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, storage.ErrNotAccessible) {
			return c.JSON(
				http.StatusNoContent, model.ErrorResponse{
					Error: controller.ErrNoContent.Error(),
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
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
			},
		)
	}
	ctrl.log.Info("HandleGetMyProjectById: logged in", zap.String("user_id", userID.String()))

	projectIDStr = c.Param("id")
	projectID, err = uuid.Parse(projectIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: storage.ErrBadRequestId.Error(),
		},
		)
	}

	project, err := ctrl.store.Project().GetByID(c.Request().Context(), projectID)
	if err != nil {
		if errors.Is(err, storage.ErrNotAccessible) {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: storage.ErrNotFound.Error(),
			},
			)
		}
		return err
	}
	return c.JSON(http.StatusOK, project)

}

// ./api/todos

func (ctrl *Controller) HandleCreateTodo(c echo.Context) error {
	var (
		request model.TodoCreateRequest
		userID  uuid.UUID
	)

	// Validate user with Token returning userID
	userId, err := ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
			},
		)
	}
	ctrl.log.Info("HandleCreateTodo: logged in", zap.String("user_id", userID.String()))

	if err = c.Bind(&request); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: controller.ErrBindingRequest.Error(),
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
	if errors.Is(err, storage.ErrAlreadyExists) {
		ctrl.log.Error("project already exists", zap.Error(err))
		return c.JSON(
			http.StatusConflict,
			model.ErrorResponse{
				Error: storage.ErrAlreadyExists.Error(),
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
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
			},
		)
	}
	ctrl.log.Info("HandleGetTodosById: logged in", zap.String("user_id", userID.String()))

	id = c.Param("id")
	todoID, err = uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: storage.ErrBadRequestId.Error(),
		},
		)
	}

	todo, err := ctrl.store.Todo().GetByID(c.Request().Context(), todoID)
	if err != nil {
		if errors.Is(err, storage.ErrNotAccessible) {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: storage.ErrNotFound.Error(),
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
				Error: storage.ErrBadRequestId.Error(),
			},
		)
	}

	//get request
	if err := c.Bind(&request); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: storage.ErrBadRequestId.Error(),
			},
		)
	}

	todo, err := ctrl.store.Todo().GetByID(c.Request().Context(), todoID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return c.JSON(
				http.StatusNotFound,
				model.ErrorResponse{
					Error: storage.ErrNotFound.Error(),
				},
			)
		}
		return c.JSON(
			http.StatusInternalServerError,
			model.ErrorResponse{
				Error: storage.ErrInternalServer.Error(),
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
				Error: storage.ErrInternalServer.Error(),
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
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
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
				Error: storage.ErrBadRequestId.Error(),
			},
		)
	}

	err = ctrl.store.Todo().Delete(c.Request().Context(), todoID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return c.JSON(
				http.StatusNotFound,
				model.ErrorResponse{
					Error: storage.ErrNotFound.Error(),
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
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
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
				Error: storage.ErrGetByID.Error(),
			},
		)
	}

	return c.JSON(http.StatusOK, listTodos)
}

// ./api/columns

func (ctrl *Controller) HandleCreateColumn(c echo.Context) error {
	var (
		request model.ColumRequest
		userID  uuid.UUID
		err     error
	)

	// Validate user with Token returning userID
	userID, err = ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
			},
		)
	}
	ctrl.log.Info("HandleCreateColumn: logged in", zap.String("user_id", userID.String()))

	if err = c.Bind(&request); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: controller.ErrBindingRequest.Error(),
			},
		)
	}

	column := &model.ColumDTO{
		ProjectId: request.ProjectId,
		Name:      request.Name,
		Order:     request.Order,
	}

	err = ctrl.store.Column().CreateColumn(c.Request().Context(), column)
	if errors.Is(err, storage.ErrAlreadyExists) {
		ctrl.log.Error("user already exists", zap.Error(err))
		return c.JSON(
			http.StatusConflict,
			model.ErrorResponse{
				Error: storage.ErrAlreadyExists.Error(),
			},
		)
	}
	response := model.ColumRequest{
		ProjectId: request.ProjectId,
		Name:      request.Name,
		Order:     request.Order,
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
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
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
				Error: storage.ErrBadRequestId.Error(),
			})
	}
	columnName = c.Param("name")

	err = ctrl.store.Column().DeleteColumn(c.Request().Context(), columnName, projectUUID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.ErrorResponse{
				Error: storage.ErrInternalServer.Error(),
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
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
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
				Error: storage.ErrBadRequestId.Error(),
			})
	}
	columnName = c.Param("name")
	column, err := ctrl.store.Column().GetColumnByName(c.Request().Context(), columnName, projectUUID)
	if err != nil {
		if errors.Is(err, storage.ErrNotAccessible) {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: storage.ErrNotFound.Error(),
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
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
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
				Error: storage.ErrBadRequestId.Error(),
			})
	}

	columnName = c.Param("name")
	if err = c.Bind(&request); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: storage.ErrBadRequestId.Error(),
			},
		)
	}
	column, err := ctrl.store.Column().GetColumnByName(c.Request().Context(), columnName, projectUUID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return c.JSON(
				http.StatusNotFound,
				model.ErrorResponse{
					Error: storage.ErrNotFound.Error(),
				},
			)
		}
		return c.JSON(
			http.StatusInternalServerError,
			model.ErrorResponse{
				Error: storage.ErrInternalServer.Error(),
			},
		)
	}

	column.Name = request.Name

	err = ctrl.store.Column().UpdateColumn(c.Request().Context(), column, columnName, projectUUID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.ErrorResponse{
				Error: storage.ErrInternalServer.Error(),
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
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
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
				Error: storage.ErrBadRequestId.Error(),
			})
	}

	listColumns, err = ctrl.store.Column().GetAllColumns(c.Request().Context(), projectUUID)
	if err != nil {
		ctrl.log.Error("error while getting group by id from DB", zap.Error(err))
		return c.JSON(
			http.StatusNoContent,
			model.ErrorResponse{
				Error: storage.ErrGetByID.Error(),
			},
		)
	}

	return c.JSON(http.StatusOK, listColumns)
}
