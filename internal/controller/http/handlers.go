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
			http.StatusBadRequest,
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
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: storage.ErrGetByLogin.Error(),
			},
		)
	}

	// Compare hashed password from request and from DB
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		ctrl.log.Error("invalid password", zap.Error(controller.InvalidPassword))
		return c.JSON(
			http.StatusUnauthorized,
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

	// Validate user with Token returning id
	id, err := ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
			},
		)
	}
	ctrl.log.Info("HandleChangePassword : logged in", zap.String("user_id", id.String()))

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
	user, err := ctrl.store.User().GetByID(c.Request().Context(), id)
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
				Error: controller.ErrPasswordAreNotEqual.Error(), // StatusConflict or what?
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
		me            *model.UserDTO
		err           error
		requestUserID uuid.UUID
	)

	// Taking a UserID from request
	requestUserID, err = ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
			},
		)
	}
	ctrl.log.Info("HandleGetMe: logged in", zap.String("user_id", requestUserID.String()))

	// Getting "Me" from DB
	me, err = ctrl.store.User().GetByID(c.Request().Context(), requestUserID)
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

	// Taking a UserID from request
	requestUserID, err := ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
			},
		)
	}
	ctrl.log.Info("HandleGetAll: logged in", zap.String("user_id", requestUserID.String()))

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
		//refreshToken  string
		requestUserID uuid.UUID
		err           error
	)
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

	// Taking a UserID from request
	requestUserID, err = ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
			},
		)
	}

	// Generating token's for the user
	accessToken, refreshToken, err := ctrl.generateAccessAndRefreshTokenForUser(requestUserID)
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
		ID:           requestUserID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return c.JSON(http.StatusCreated, response)
}

// ./api/groups
func (ctrl *Controller) HandleCreateInvite(c echo.Context) error {
	var (
		request model.GroupRequest
		err     error
	)

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

	// Validate user with Token returning id
	requestUserID, err := ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
			},
		)
	}

	ctrl.log.Info("HandleAddInGroup : got user invitor", zap.String("added_by_user_id", requestUserID.String()))

	// Taking group from GroupDTO with new data
	user := &model.GroupDTO{
		UserID:    requestUserID,
		ProjectID: request.ProjectID,
	}
	ctrl.log.Info("got user", zap.Any("user", user))

	// Inserting in DB the group
	err = ctrl.store.Group().CreateGroup(c.Request().Context(), user)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			ctrl.log.Error("user already in group", zap.Error(err))
			return c.JSON(
				http.StatusConflict,
				model.ErrorResponse{
					Error: storage.ErrAlreadyExists.Error(),
				},
			)
		}
		ctrl.log.Error("got error while adding in group", zap.Error(err))
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: storage.ErrCreateGroup.Error(),
			},
		)
	}

	response := model.GroupResponse{
		UserID:    user.UserID,
		ProjectID: user.ProjectID,
	}

	return c.JSON(http.StatusCreated, response)
}

func (ctrl *Controller) HandleGetGroupByID(c echo.Context) error {
	var (
		group         []model.GroupDTO
		err           error
		requestUserID uuid.UUID
	)

	// Taking a UserID from request
	requestUserID, err = ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
			},
		)
	}
	// Todo: add this log example to all handlers
	ctrl.log.Info("HandleGetGroup: got user id", zap.String("user_id", requestUserID.String()))

	// Getting "Group" from DB
	group, err = ctrl.store.Group().GetUsersInProjectByProjectID(c.Request().Context())
	if err != nil {
		ctrl.log.Error("error while getting group by id from DB", zap.Error(err))
		return c.JSON(
			http.StatusNoContent,
			model.ErrorResponse{
				Error: storage.ErrGetByID.Error(),
			},
		)
	}

	response := &model.GroupResponse{
		ID:        group.ID,
		Name:      group.Name,
		CreatedBy: group.CreatedBy,
	}

	return c.JSON(http.StatusOK, response)
}

//func (ctrl *Controller) HandleCreateInvite(c echo.Context) error {
//	url := fmt.Sprintf("?user_id=%s&project_id=%s")
//}

func (ctrl *Controller) HandleGetMyGroups(c echo.Context) error {
	var (
		listGroups []model.GroupDTO
		err        error
		UserID     uuid.UUID
	)

	// Taking a UserID from request
	UserID, err = ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
			},
		)
	}
	ctrl.log.Info("HandleGetAll: logged in", zap.String("user_id", UserID.String()))

	// Taking a UserID from request
	UserID, err = ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
			},
		)
	}
	// Todo: add this log example to all handlers
	ctrl.log.Info("HandleGetGroup: got user id", zap.String("user_id", UserID.String()))

	// Getting list of "Groups" from DB
	listGroups, err = ctrl.store.Group().GetMyGroups(c.Request().Context(), UserID)
	if err != nil {
		ctrl.log.Error("error while getting group by id from DB", zap.Error(err))
		return c.JSON(
			http.StatusNoContent,
			model.ErrorResponse{
				Error: storage.ErrGetByID.Error(),
			},
		)
	}

	return c.JSON(http.StatusOK, listGroups)
}

// ./api/projects

// ./api/todos
func (ctrl *Controller) HandleCreateTodo(c echo.Context) error {
	var request model.TodoCreateRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: controller.ErrBindingRequest.Error(),
			},
		)
	}
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

	todo := &model.TodoDTO{
		ID:          uuid.New(),
		Name:        request.Name,
		Description: request.Description,
		IsCompleted: request.IsCompleted,
		ProjectID:   request.ProjectID,
		CreatedBy:   userId,
		Column:      request.Column,
	}

	err = ctrl.store.Todo().Create(c.Request().Context(), todo, userId, todo.ProjectID, todo.Column)
	if errors.Is(err, storage.ErrAlreadyExists) {
		ctrl.log.Error("user already exists", zap.Error(err))
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
	return c.JSON(http.StatusCreated, todo)

}

func (ctrl *Controller) HandleGetTodosById(c echo.Context) error {
	id := c.Param("id")
	todoID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: storage.ErrBadRequestId.Error(),
		})
	}
	_, err = ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
			},
		)
	}

	todo, err := ctrl.store.Todo().GetByID(c.Request().Context(), todoID)
	if err != nil {
		if errors.Is(err, storage.ErrNotAccessible) {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: storage.ErrNotFound.Error(),
			})
		}
		return err
	}
	return c.JSON(http.StatusOK, todo)
}

func (ctrl *Controller) HandleChangeTodo(c echo.Context) error {
	var request model.TodoUpdateRequest

	// get user id
	user, err := ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: "Unauthorized",
			},
		)
	}

	// parse id
	todoIDStr := c.Param("id")
	todoID, err := uuid.Parse(todoIDStr)
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

	// get user by id
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
	//checking
	if todo.CreatedBy != user {
		return c.JSON(
			http.StatusForbidden,
			model.ErrorResponse{
				Error: storage.ErrForbidden.Error(),
			},
		)
	}
	//change todos
	todo.Name = request.Name
	todo.Description = request.Description
	todo.IsCompleted = request.IsCompleted

	//work with db
	err = ctrl.store.Todo().Update(c.Request().Context(), todo, user)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.ErrorResponse{
				Error: storage.ErrInternalServer.Error(),
			},
		)
	}

	ctrl.log.Info("successfully updated todo", zap.Any("todo", todo))
	return c.JSON(http.StatusOK, todo)
}

func (ctrl *Controller) HandleDeleteTodo(c echo.Context) error {
	todoIDStr := c.Param("id")
	todoID, err := uuid.Parse(todoIDStr)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: storage.ErrBadRequestId.Error(),
			},
		)
	}

	_, err = ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
			},
		)
	}

	err = ctrl.store.Todo().DeleteTodos(c.Request().Context(), todoID)
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
		UserID    uuid.UUID
		Columns   string
	)

	// Taking a UserID from request
	UserID, err = ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
			},
		)
	}
	ctrl.log.Info("HandleGetAll: logged in", zap.String("user_id", UserID.String()))

	// Taking a UserID from request
	UserID, err = ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
			},
		)
	}
	// Todo: add this log example to all handlers
	ctrl.log.Info("HandleGetAllTodos: got user id", zap.String("user_id", UserID.String()))

	// Getting list of "Groups" from DB
	listTodos, err = ctrl.store.Todo().GetAll(c.Request().Context(), UserID)
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
