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
func (ctrl *Controller) HandleCreateGroup(c echo.Context) error {
	var (
		request   model.GroupDTO
		createdBy uuid.UUID
		err       error
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
	createdBy, err = ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(
			http.StatusUnauthorized,
			model.ErrorResponse{
				Error: controller.ErrValidationToken.Error(),
			},
		)
	}

	ctrl.log.Info("HandleCreateGroup : got user", zap.String("created_by_user_id", createdBy.String()))

	// Taking group from GroupDTO with new data
	group := &model.GroupDTO{
		ID:        uuid.New(),
		Name:      request.Name,
		CreatedBy: createdBy,
	}
	ctrl.log.Info("got group", zap.Any("group", group))

	// Inserting in DB the group
	err = ctrl.store.Group().Create(c.Request().Context(), group)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			ctrl.log.Error("group already exists", zap.Error(err))
			return c.JSON(
				http.StatusConflict,
				model.ErrorResponse{
					Error: storage.ErrAlreadyExists.Error(),
				},
			)
		}
		ctrl.log.Error("got error while creating group", zap.Error(err))
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: storage.ErrCreateGroup.Error(),
			},
		)
	}

	response := model.GroupResponse{
		ID:        group.ID,
		Name:      group.Name,
		CreatedBy: group.CreatedBy,
	}

	return c.JSON(http.StatusCreated, response)
}

func (ctrl *Controller) HandleGetGroupByID(c echo.Context) error {
	var (
		group         *model.GroupDTO
		err           error
		requestUserID uuid.UUID
		requestID     uuid.UUID
	)

	// Binding requestID
	if err = c.Bind(&requestID); err != nil {
		ctrl.log.Error("error while binding request", zap.Error(err))
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
	// Todo: add this log example to all handlers
	ctrl.log.Info("HandleGetGroup: got user id", zap.String("user_id", requestUserID.String()))

	// Getting "Group" from DB
	group, err = ctrl.store.Group().GetByID(c.Request().Context(), requestID)
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

func (ctrl *Controller) HandleCreateInvite(c echo.Context) error {
	return nil
}

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

// ./api/todos
func (ctrl *Controller) HandleCreateTodo(c echo.Context) error {
	var request model.TodoCreateRequest
	user, err := ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		return err
	}

	if err := c.Bind(&request); err != nil {
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
		IsCompleted: false,
		CreatedBy:   user,
	}

	err = ctrl.store.Todo().Create(c.Request().Context(), todo)
	if errors.Is(err, storage.ErrAlreadyExists) {
		ctrl.log.Error("user already exists", zap.Error(err))
		return c.JSON(
			http.StatusConflict,
			model.ErrorResponse{
				Error: storage.ErrAlreadyExists.Error(),
			},
		)
	}
	ctrl.log.Info("successfully created todo", zap.Any("todo", todo))
	return c.JSON(http.StatusCreated, todo)

}

func (ctrl *Controller) HandleGetTodosById(c echo.Context) error {
	id := c.Param("id")
	todoID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: "Invalid todo ID",
		})
	}

	todo, err := ctrl.store.Todo().GetByID(c.Request().Context(), todoID)
	if err != nil {
		if errors.Is(err, storage.ErrNotAccessible) {
			return c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: "Todo not found",
			})
		}
		return err
	}
	return c.JSON(http.StatusOK, todo)
}

func (ctrl *Controller) HandleUpdateTodo(c echo.Context) error {
	var request model.TodoUpdateRequest
	todoIDStr := c.Param("id")
	todoID, err := uuid.Parse(todoIDStr)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: "invalid todo ID format",
			},
		)
	}

	user, err := ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		return err
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			model.ErrorResponse{
				Error: controller.ErrBindingRequest.Error(),
			},
		)
	}

	todo, err := ctrl.store.Todo().GetByID(c.Request().Context(), todoID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return c.JSON(
				http.StatusNotFound,
				model.ErrorResponse{
					Error: storage.ErrNotAccessible.Error(),
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

	if todo.CreatedBy != user {
		return c.JSON(
			http.StatusForbidden,
			model.ErrorResponse{
				Error: "you do not have permission to update this todo",
			},
		)
	}

	todo.Name = request.Name
	todo.Description = request.Description
	todo.IsCompleted = request.IsCompleted

	err = ctrl.store.Todo().Update(c.Request().Context(), todo)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			model.ErrorResponse{
				Error: err.Error(),
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
				Error: "invalid todo ID format",
			},
		)
	}

	user, err := ctrl.getUserIDFromRequest(c.Request())
	if err != nil {
		return err
	}

	err = ctrl.store.Todo().DeleteTodos(c.Request().Context(), todoID, user)
	if err != nil {
		if err == storage.ErrNotFound {
			return c.JSON(
				http.StatusNotFound,
				model.ErrorResponse{
					Error: "todo not found",
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
