package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
)

type userCredentials struct {
	username string
	password string
}

func getCredentials(basicAuth []string) (*userCredentials, error) {

	if len(basicAuth) != 1 {
		return nil, errors.New("ambiguity: basic token not found or it in proper format")
	}
	token := strings.Split(basicAuth[0], "Basic")
	if len(token) != 2 {
		return nil, errors.New("basic token not in proper format")
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(token[1]))
	if err != nil {
		return nil, fmt.Errorf("decode error: %v", err)
	}
	result := strings.Split((string)(decoded), ":")
	if len(result) != 2 {
		return nil, errors.New("basic token not in proper format")
	}

	return &userCredentials{
		username: result[0],
		password: result[1],
	}, nil
}

func (h *UserHandler) authenticateLogin(c echo.Context) (*model.User, int, utils.ErrorResponse) {
	creds, err := getCredentials(c.Request().Header.Values("authorization"))
	if err != nil {
		return nil, http.StatusUnauthorized, utils.NewErrorResponse(http.StatusUnauthorized, utils.ErrorSeverity, "parsing token", err)
	}

	u, err := h.userStore.GetByUsername(creds.username)
	if err != nil {
		return nil, http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "encountered error looking for the user", err)
	}
	if u == nil {
		return nil, http.StatusForbidden, utils.NewErrorResponse(http.StatusForbidden, utils.CriticalSeverity, "invalid username or password", err)
	}
	if u.Password != creds.password {
		return nil, http.StatusForbidden, utils.NewErrorResponse(http.StatusForbidden, utils.CriticalSeverity, "invalid username or password", err)
	}
	return u, http.StatusOK, utils.ErrorResponse{}

}

// Login godoc
// @Summary Login for existing user
// @Description Login for existing user
// @ID login
// @Tags user
// @Accept  json
// @Produce  json
// @Security BasicAuth
// @Success 200 {string} string "Bearer Token for Logged in User"
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /users/login [post]
func (h *UserHandler) Login(c echo.Context) error {
	u, code, err := h.authenticateLogin(c)
	if u != nil {
		return c.JSON(http.StatusOK, newUserResponse(u))
	}

	return c.JSON(code, err)
}

// ChangePasword godoc
// @Summary Change password for existing user
// @Description Change password for existing user
// @ID change-password
// @Tags user
// @Accept  json
// @Produce  json
// @Security BasicAuth
// @Param password query string true "Enter New Password" format(password)
// @Success 204 "No Content"
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /users/change-password [patch]
func (h *UserHandler) ChangePasword(c echo.Context) error {
	u, code, err := h.authenticateLogin(c)
	if u == nil {
		return c.JSON(code, err)
	}

	u.Password = c.QueryParam("password")
	if err := h.userStore.Update(u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.ErrorSeverity, "", err))
	}

	return c.JSON(http.StatusNoContent, nil)
}
