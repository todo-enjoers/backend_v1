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

func (ctrl *Controller) HandleRegister(c echo.Context) error {
	var request model.UserRegisterRequest

	// Binding request
	if err := c.Bind(&request); err != nil {
		ctrl.log.Error("could not bind request", zap.Error(err))
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()}) //change
	}

	// Validate request
	if ok, err := request.Validate(); ok {
		ctrl.log.Error("invalid request", zap.Error(err))
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Hashing password from request
	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
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
			return c.JSON(http.StatusConflict, echo.Map{"error": err.Error()})
		}
		ctrl.log.Error("got error while creating user", zap.Error(err))
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	ctrl.log.Info("successfully created user")

	// Generating token's for the user
	accessToken, refreshToken, err := ctrl.generateAccessAndRefreshTokenForUser(user.ID)
	if err != nil {
		ctrl.log.Error("error while creating tokens", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
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
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()}) //change
	}

	// Getting the "User" from DB
	user, err := ctrl.store.User().GetByLogin(c.Request().Context(), request.Login)
	if err != nil {
		ctrl.log.Error("error while getting user by login from DB", zap.Error(err))
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	// Compare hashed password from request and from DB
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		ctrl.log.Error("invalid password", zap.Error(controller.InvalidPassword))
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()}) // StatusUnauthorized or StatusBadRequest
	}

	// Generating access, refresh tokens for logged user
	access, refresh, err := ctrl.generateAccessAndRefreshTokenForUser(user.ID)
	if err != nil {
		ctrl.log.Error("error while creating tokens", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
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
	id, err := ctrl.getUserIDFromAccessToken(c)
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}
	ctrl.log.Info("HandleChangePassword : logged in", zap.String("user_id", id.String()))

	// Binding request
	if err := c.Bind(&request); err != nil {
		ctrl.log.Error("error while binding request", zap.Error(err))
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Getting the "User" from DB
	user, err := ctrl.store.User().GetByID(c.Request().Context(), id)
	if err != nil {
		ctrl.log.Error("error while getting user by login from DB", zap.Error(err))
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	// Compare hashed password from request and from DB
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword)); err != nil {
		ctrl.log.Error("invalid password", zap.Error(controller.InvalidPassword)) // return controller.InvalidPassword or echo.map
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()}) // StatusUnauthorized or StatusBadRequest
	}

	// Compare NewPassword and  NewPasswordAgain
	if request.NewPassword != request.NewPasswordAgain {
		ctrl.log.Error("password are not equal", zap.Error(controller.ErrPasswordAreNotEqual))// return controller.ErrPasswordAreNotEqual or echo.map
		return c.JSON(http.StatusConflict, controller.ErrPasswordAreNotEqual) // StatusConflict or what?
	}

	// Hashing NewPassword from request
	HashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Inserting NewPassword in DB
	if err = ctrl.store.User().ChangePassword(c.Request().Context(), string(HashedNewPassword), user.ID) // ???
	if err != nil {
		ctrl.log.Error("error while inserting in DB changed password", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err)
	}

	response := &model.UserChangePasswordResponse{
		OldPassword: request.OldPassword,
		NewPassword: HashedNewPassword,
		NewPasswordAgain: HashedNewPassword,
	}

	return c.JSON(http.StatusOK, response)

}

func (ctrl *Controller) HandleGetMe(c echo.Context) error {

	// Validate user with Token returning id
	requestUserID, err := ctrl.getUserIDFromAccessToken(c)
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()}) // here another question
	}
	ctrl.log.Info("HandleGetMe: logged in", zap.String("user_id", requestUserID.String()))

	// Getting "Me" from DB
	me, err := ctrl.store.User().GetByID(c.Request().Context(), requestUserID)
	if err != nil {
		ctrl.log.Error("error while getting user by id from DB", zap.Error(err))
		return c.JSON(http.StatusConflict, echo.Map{"error": err.Error()}) // question status
	}

	response := &model.UserGetMeResponse{
		ID:   me.ID,
		Name: me.Login,
	}

	return c.JSON(http.StatusOK, response)
}

func (ctrl *Controller) HandleGetAll(c echo.Context) error {
	var list []model.UserDTO

	// Validate user with Token returning id
	requestUserID, err := ctrl.getUserIDFromAccessToken(c)
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(controller.ErrValidationToken))
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()}) // here another question
	}
	ctrl.log.Info("HandleGetMe: logged in", zap.String("user_id", requestUserID.String()))

	list, err = ctrl.store.User().GetAll(c.Request().Context())
	if err != nil {
		ctrl.log.Error("error while getting users by id from DB", zap.Error(err))
	}

	// response ???

	return c.JSON(http.StatusOK, list)
}
