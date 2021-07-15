package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/router"
	"github.com/dell/csm-deployment/store/mocks"
	"github.com/dell/csm-deployment/utils"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func generateBasicAuth(authType, username, password string) string {
	return fmt.Sprintf("%s %s", authType, base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password))))
}
func Test_Login(t *testing.T) {

	tests := map[string]func(t *testing.T) (int, *UserHandler, string, string, *gomock.Controller){
		"success": func(*testing.T) (int, *UserHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			username := "admin"
			password := "password123"
			expectedResponse := fmt.Sprintf(`"%s"`, utils.GenerateJWT(username))

			userStore := mocks.NewMockUserStoreInterface(ctrl)
			user := model.User{
				Username: username,
				Password: password,
			}
			userStore.EXPECT().GetByUsername(gomock.Any()).Times(1).Return(&user, nil)

			handler := &UserHandler{userStore}

			return http.StatusOK, handler, generateBasicAuth("Basic", username, password), expectedResponse, ctrl
		},
		"wrong password": func(*testing.T) (int, *UserHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			username := "admin"
			password := "wrong-password123"
			expectedResponse := "{\"http_status_code\":403,\"messages\":[{\"code\":403,\"message\":\"the password, wrong-password123, is forbidden\",\"message_l10n\":null,\"Arguments\":null,\"severity\":\"CRITICAL\"}]}"

			userStore := mocks.NewMockUserStoreInterface(ctrl)
			user := model.User{
				Username: username,
				Password: "password123",
			}
			userStore.EXPECT().GetByUsername(gomock.Any()).Times(1).Return(&user, nil)

			handler := &UserHandler{userStore}

			return http.StatusForbidden, handler, generateBasicAuth("Basic", username, password), expectedResponse, ctrl
		},
		"nil result from db": func(*testing.T) (int, *UserHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			username := "admin"
			password := "password123"
			expectedResponse := ""

			userStore := mocks.NewMockUserStoreInterface(ctrl)
			userStore.EXPECT().GetByUsername(gomock.Any()).Times(1).Return(nil, nil)
			handler := &UserHandler{userStore}
			return http.StatusForbidden, handler, generateBasicAuth("Basic", username, password), expectedResponse, ctrl
		},
		"error querying db": func(*testing.T) (int, *UserHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			username := "admin"
			password := "password123"
			expectedResponse := ""

			userStore := mocks.NewMockUserStoreInterface(ctrl)
			userStore.EXPECT().GetByUsername(gomock.Any()).Times(1).Return(nil, errors.New("error"))
			handler := &UserHandler{userStore}
			return http.StatusInternalServerError, handler, generateBasicAuth("Basic", username, password), expectedResponse, ctrl
		},
		"bad authorization type": func(*testing.T) (int, *UserHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			username := "admin"
			password := "password123"
			expectedResponse := "{\"http_status_code\":401,\"messages\":[{\"code\":401,\"message\":\"parsing token\",\"message_l10n\":\"basic token not in proper format\",\"Arguments\":null,\"severity\":\"ERROR\"}]}"
			handler := &UserHandler{}
			return http.StatusUnauthorized, handler, generateBasicAuth("Bearer", username, password), expectedResponse, ctrl
		},
		"error decoding": func(*testing.T) (int, *UserHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			expectedResponse := "{\"http_status_code\":401,\"messages\":[{\"code\":401,\"message\":\"parsing token\",\"message_l10n\":\"decode error: illegal base64 data at input byte 5\",\"Arguments\":null,\"severity\":\"ERROR\"}]}"
			handler := &UserHandler{}
			return http.StatusUnauthorized, handler, "Basic admin:password", expectedResponse, ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			expectedStatus, handler, basicAuth, expectedResponse, ctrl := tc(t)

			e := router.New()
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Add("authorization", basicAuth)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			assert.NoError(t, handler.Login(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}
func Test_ChangePasword(t *testing.T) {
	tests := map[string]func(t *testing.T) (int, *UserHandler, string, string, *gomock.Controller){
		"success": func(*testing.T) (int, *UserHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			username := "admin"
			password := "password123"
			expectedResponse := "null"

			userStore := mocks.NewMockUserStoreInterface(ctrl)
			user := model.User{
				Username: username,
				Password: password,
			}
			userStore.EXPECT().GetByUsername(gomock.Any()).Times(1).Return(&user, nil)
			userStore.EXPECT().Update(gomock.Any()).Times(1).Return(nil)

			handler := &UserHandler{userStore}

			return http.StatusNoContent, handler, generateBasicAuth("Basic", username, password), expectedResponse, ctrl
		},
		"error updating db": func(*testing.T) (int, *UserHandler, string, string, *gomock.Controller) {

			ctrl := gomock.NewController(t)

			username := "admin"
			password := "password123"
			expectedResponse := "{\"http_status_code\":422,\"messages\":[{\"code\":422,\"message\":\"Unprocessable Entity\",\"message_l10n\":\"error\",\"Arguments\":null,\"severity\":\"ERROR\"}]}"

			userStore := mocks.NewMockUserStoreInterface(ctrl)
			user := model.User{
				Username: username,
				Password: password,
			}
			userStore.EXPECT().GetByUsername(gomock.Any()).Times(1).Return(&user, nil)
			userStore.EXPECT().Update(gomock.Any()).Times(1).Return(errors.New("error"))

			handler := &UserHandler{userStore}

			return http.StatusUnprocessableEntity, handler, generateBasicAuth("Basic", username, password), expectedResponse, ctrl
		},
		"error decoding": func(*testing.T) (int, *UserHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			expectedResponse := "{\"http_status_code\":401,\"messages\":[{\"code\":401,\"message\":\"parsing token\",\"message_l10n\":\"decode error: illegal base64 data at input byte 5\",\"Arguments\":null,\"severity\":\"ERROR\"}]}"
			handler := &UserHandler{}
			return http.StatusUnauthorized, handler, "Basic admin:password", expectedResponse, ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			expectedStatus, handler, basicAuth, expectedResponse, ctrl := tc(t)

			e := router.New()
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Add("authorization", basicAuth)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.QueryParams().Add("password", "password")

			assert.NoError(t, handler.ChangePasword(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}
