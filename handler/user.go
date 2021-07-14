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

func (h *Handler) authenticateLogin(c echo.Context) (*model.User, int, utils.ErrorResponse) {
	creds, err := getCredentials(c.Request().Header.Values("authorization"))
	if err != nil {
		return nil, http.StatusUnauthorized, utils.NewErrorResponse(http.StatusUnauthorized, utils.ErrorSeverity, "parsing token", err)
	}

	u, err := h.userStore.GetByUsername(creds.username)
	if err != nil {
		v := fmt.Sprintf("encountered error looking for the user: %s", creds.username)
		return nil, http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, v, err)
	}
	if u == nil {
		v := fmt.Sprintf("the user, %s, is forbidden", creds.username)
		return nil, http.StatusForbidden, utils.NewErrorResponse(http.StatusForbidden, utils.CriticalSeverity, v, err)
	}
	if u.Password != creds.password {
		v := fmt.Sprintf("the password, %s, is forbidden", creds.password)
		return nil, http.StatusForbidden, utils.NewErrorResponse(http.StatusForbidden, utils.CriticalSeverity, v, err)
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
// @Success 200 {object} userResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /users/login [post]
func (h *Handler) Login(c echo.Context) error {
	u, code, err := h.authenticateLogin(c)
	if u != nil {
		return c.JSON(http.StatusOK, newUserResponse(u))
	}

	return c.JSON(code, err)
}

// UpdateUser godoc
// @Summary Update current user
// @Description Update user information for current user
// @ID update-user
// @Tags user
// @Accept  json
// @Produce  json
// @Security BasicAuth
// @Param user body userUpdateRequest true "User details to update. At least **one** field is required."
// @Success 204 "No Content"
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /users/update [patch]
func (h *Handler) UpdateUser(c echo.Context) error {
	u, code, err := h.authenticateLogin(c)
	if u == nil {
		return c.JSON(code, err)
	}
	req := new(userUpdateRequest)
	req.populate(u)
	if err := req.bind(c, u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.ErrorSeverity, "", err))
	}
	if err := h.userStore.Update(u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.ErrorSeverity, "", err))
	}

	return c.JSON(http.StatusNoContent, nil)

}
