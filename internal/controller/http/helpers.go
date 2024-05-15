package http

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/todo-enjoers/backend_v1/internal/model"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func (ctrl *Controller) generateAccessAndRefreshTokenForUser(userID uuid.UUID) (access string, refreshToken string, err error) {
	//Creating access token for userID if isAccess is true
	access, err = ctrl.token.CreateTokenForUser(userID, true)
	if err != nil {
		return "", "", fmt.Errorf("error creating access token for userID: %w", err)
	}
	//Creating refresh token for userID if isAccess is false
	refreshToken, err = ctrl.token.CreateTokenForUser(userID, false)
	if err != nil {
		return "", "", fmt.Errorf("error creating refresh token for userID: %w", err)
	}
	return
}

func GetJWTFromBearerToken(raw string) (string, error) {
	splitToken := strings.Split(raw, "Bearer")
	if len(splitToken) != 2 {
		return "", fmt.Errorf("bearer token not in proper format")
	}
	reqToken := strings.TrimSpace(splitToken[1])
	return reqToken, nil
}

func (ctrl *Controller) getUserDataFromRequest(req *http.Request) (*model.UserDataInToken, error) {
	var (
		data     string
		token    string
		err      error
		userData *model.UserDataInToken
	)
	// Taking header for token from request
	data = req.Header.Get("Authorization")
	ctrl.log.Info("got authorization header", zap.Any("header_data", data))

	//Parsing token
	token, err = GetJWTFromBearerToken(data)
	if err != nil {
		return nil, fmt.Errorf("error while parsing token: %w", err)
	}

	//Getting userData
	userData, err = ctrl.token.GetDataFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("error while getting data from token: %w", err)
	}

	return userData, nil
}

func (ctrl *Controller) getUserIDFromRequest(req *http.Request) (uuid.UUID, error) {
	userData, err := ctrl.getUserDataFromRequest(req)
	if err != nil {
		return uuid.Nil, fmt.Errorf("unauthorized: error while getting user data from request: %w", err)
	}
	if !userData.IsAccess {
		return uuid.Nil, fmt.Errorf("unauthorized: error while checking user access: %w", err)
	}
	return userData.ID, nil
}

func (ctrl *Controller) PasswordToHash(raw string) ([]byte, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ctrl *Controller) CompareHashes(passReq []byte, passDB []byte) (err error) {
	if err = bcrypt.CompareHashAndPassword(passReq, passDB); err != nil {
		return fmt.Errorf("error while comparing hashes: %w", err)
	}
	return nil
}
