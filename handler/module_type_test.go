package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/router"
	"github.com/dell/csm-deployment/store/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_GetModuleType(t *testing.T) {
	tests := map[string]func(t *testing.T) (int, *ModuleTypeHandler, string, string, *gomock.Controller){
		"success": func(*testing.T) (int, *ModuleTypeHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			getModuleTypeResponseJSON := "{\"id\":1,\"name\":\"module-1\",\"version\":\"1.2.3\",\"standalone\":false}"

			moduleTypeStore := mocks.NewMockModuleStoreInterface(ctrl)
			moduleType := model.ModuleType{
				Version: "1.2.3",
				Name:    "module-1",
			}
			moduleType.ID = 1

			moduleTypeStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&moduleType, nil)
			handler := &ModuleTypeHandler{moduleTypeStore}
			return http.StatusOK, handler, "1", getModuleTypeResponseJSON, ctrl
		},
		"nil result from db": func(*testing.T) (int, *ModuleTypeHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			moduleTypeStore := mocks.NewMockModuleStoreInterface(ctrl)
			moduleTypeStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, nil)
			handler := &ModuleTypeHandler{moduleTypeStore}
			return http.StatusNotFound, handler, "1", "", ctrl
		},
		"error querying db": func(*testing.T) (int, *ModuleTypeHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			moduleTypeStore := mocks.NewMockModuleStoreInterface(ctrl)
			moduleTypeStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, errors.New("error"))
			handler := &ModuleTypeHandler{moduleTypeStore}
			return http.StatusInternalServerError, handler, "1", "", ctrl
		},
		"id is not numeric": func(*testing.T) (int, *ModuleTypeHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			moduleTypeStore := mocks.NewMockModuleStoreInterface(ctrl)
			handler := &ModuleTypeHandler{moduleTypeStore}
			return http.StatusUnprocessableEntity, handler, "abc", "", ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			expectedStatus, handler, moduleID, expectedResponse, ctrl := tc(t)

			e := router.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/module-types/:id")
			c.SetParamNames("id")
			c.SetParamValues(moduleID)

			assert.NoError(t, handler.GetModuleType(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}
func Test_ListModuleType(t *testing.T) {
	tests := map[string]func(t *testing.T) (int, *ModuleTypeHandler, string, *gomock.Controller){
		"success": func(*testing.T) (int, *ModuleTypeHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			listModuleTypeResponseJSON := "[{\"id\":0,\"name\":\"module-1\",\"version\":\"1.2.3\",\"standalone\":false},{\"id\":0,\"name\":\"module-2\",\"version\":\"1.2.3\",\"standalone\":false}]"

			moduleTypeStore := mocks.NewMockModuleStoreInterface(ctrl)

			moduleTypeArrays := make([]model.ModuleType, 0)
			moduleTypeArrays = append(moduleTypeArrays, model.ModuleType{
				Version: "1.2.3",
				Name:    "module-1",
			})
			moduleTypeArrays = append(moduleTypeArrays, model.ModuleType{
				Version: "1.2.3",
				Name:    "module-2",
			})
			moduleTypeStore.EXPECT().GetAll().Times(1).Return(moduleTypeArrays, nil)

			handler := &ModuleTypeHandler{moduleTypeStore}

			return http.StatusOK, handler, listModuleTypeResponseJSON, ctrl
		},
		"error querying database": func(*testing.T) (int, *ModuleTypeHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			moduleTypeStore := mocks.NewMockModuleStoreInterface(ctrl)
			moduleTypeStore.EXPECT().GetAll().Times(1).Return(nil, errors.New("error"))
			handler := &ModuleTypeHandler{moduleTypeStore}

			return http.StatusInternalServerError, handler, "", ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			expectedStatus, handler, expectedResponse, ctrl := tc(t)

			e := router.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			assert.NoError(t, handler.ListModuleType(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}
