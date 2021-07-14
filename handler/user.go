package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
)

type userCredentials struct {
	Username string
	Password string
}

func getCredentials(s string) (*userCredentials, error) {
	c := &userCredentials{}

	token := strings.Split(s, "Basic")
	if len(token) != 2 {
		return nil, errors.New("basic token not in proper format")
	}
	tokenStr := strings.TrimSpace(token[1])
	decoded, err := base64.StdEncoding.DecodeString(tokenStr)

	if err != nil {
		return nil, fmt.Errorf("decode error: %v", err)
	}
	result := strings.Split((string)(decoded), ":")
	if len(result) != 2 {
		return nil, errors.New("basic token not in proper format")
	}

	return c, nil
}

// Login godoc
// @Summary Login for existing user
// @Description Login for existing user
// @ID login
// @Tags user
// @Accept  json
// @Produce  json
// @Security BasicAuth
// @Param Authorization header string true "Basic access authentication" default(Basic <base64 of username:password>)
// @Success 200 {string} string
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 422 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /users/login [post]
func (h *Handler) Login(c echo.Context) error {
	basicAuth := c.Request().Header.Values("Authorization")
	if len(basicAuth) != 0 {
		t, err := getCredentials(basicAuth[0])
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, utils.NewErrorResponse(http.StatusBadRequest, utils.ErrorSeverity, "parsing token", err))
		}
		c.Logger().Print(t)
		return c.JSON(http.StatusOK, basicAuth[0])
	}
	req := &userLoginRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.ErrorSeverity, "", err))
	}
	u, err := h.userStore.GetByUsername(req.Username)
	if err != nil {
		return c.JSON(http.StatusForbidden, utils.NewErrorResponse(http.StatusForbidden, utils.CriticalSeverity, "", err))
	}
	if u == nil {
		return c.JSON(http.StatusForbidden, utils.NewErrorResponse(http.StatusForbidden, utils.CriticalSeverity, "", err))
	}
	if !u.CheckPassword(req.Password) {
		return c.JSON(http.StatusForbidden, utils.NewErrorResponse(http.StatusForbidden, utils.CriticalSeverity, "", err))
	}
	return c.JSON(http.StatusOK, newUserResponse(u))
}
