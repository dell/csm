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

func Test_GetDriverType(t *testing.T) {
	tests := map[string]func(t *testing.T) (int, *DriverTypeHandler, string, string, *gomock.Controller){
		"success": func(*testing.T) (int, *DriverTypeHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			getDriverTypeResponseJSON := `{"id":1,"storage_array_type_ID":1,"version":"1.2.3"}`

			driverTypeStore := mocks.NewMockDriverTypeStoreInterface(ctrl)
			driverType := model.DriverType{
				Version:            "1.2.3",
				StorageArrayTypeID: 1,
			}
			driverType.ID = 1

			driverTypeStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&driverType, nil)
			handler := &DriverTypeHandler{driverTypeStore}
			return http.StatusOK, handler, "1", getDriverTypeResponseJSON, ctrl
		},
		"nil result from db": func(*testing.T) (int, *DriverTypeHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			driverTypeStore := mocks.NewMockDriverTypeStoreInterface(ctrl)
			driverTypeStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, nil)
			handler := &DriverTypeHandler{driverTypeStore}
			return http.StatusNotFound, handler, "1", "", ctrl
		},
		"error querying db": func(*testing.T) (int, *DriverTypeHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			driverTypeStore := mocks.NewMockDriverTypeStoreInterface(ctrl)
			driverTypeStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, errors.New("error"))
			handler := &DriverTypeHandler{driverTypeStore}
			return http.StatusInternalServerError, handler, "1", "", ctrl
		},
		"id is not numeric": func(*testing.T) (int, *DriverTypeHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			driverTypeStore := mocks.NewMockDriverTypeStoreInterface(ctrl)
			handler := &DriverTypeHandler{driverTypeStore}
			return http.StatusUnprocessableEntity, handler, "abc", "", ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			expectedStatus, handler, driverID, expectedResponse, ctrl := tc(t)

			e := router.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/driver-types/:id")
			c.SetParamNames("id")
			c.SetParamValues(driverID)

			assert.NoError(t, handler.GetDriverType(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}
func Test_ListDriverType(t *testing.T) {
	tests := map[string]func(t *testing.T) (int, *DriverTypeHandler, string, *gomock.Controller){
		"success": func(*testing.T) (int, *DriverTypeHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			listDriverTypeResponseJSON := "[{\"id\":0,\"storage_array_type_ID\":1,\"version\":\"1.2.3\"},{\"id\":0,\"storage_array_type_ID\":2,\"version\":\"1.2.0\"}]"

			driverTypeStore := mocks.NewMockDriverTypeStoreInterface(ctrl)

			driverTypeArrays := make([]model.DriverType, 0)
			driverTypeArrays = append(driverTypeArrays, model.DriverType{
				Version:            "1.2.3",
				StorageArrayTypeID: 1,
			})
			driverTypeArrays = append(driverTypeArrays, model.DriverType{
				Version:            "1.2.0",
				StorageArrayTypeID: 2,
			})
			driverTypeStore.EXPECT().GetAll().Times(1).Return(driverTypeArrays, nil)

			handler := &DriverTypeHandler{driverTypeStore}

			return http.StatusOK, handler, listDriverTypeResponseJSON, ctrl
		},
		"error querying database": func(*testing.T) (int, *DriverTypeHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			driverTypeStore := mocks.NewMockDriverTypeStoreInterface(ctrl)
			driverTypeStore.EXPECT().GetAll().Times(1).Return(nil, errors.New("error"))
			handler := &DriverTypeHandler{driverTypeStore}

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

			assert.NoError(t, handler.ListDriverType(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}
