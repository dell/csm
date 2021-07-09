package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dell/csm-deployment/router"
	"github.com/dell/csm-deployment/store/mocks"
	"github.com/golang/mock/gomock"

	"github.com/dell/csm-deployment/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_CreateStorageArray(t *testing.T) {

	tests := map[string]func(t *testing.T) (int, *StorageArrayHandler, string, string, *gomock.Controller){
		"success": func(*testing.T) (int, *StorageArrayHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			createStorageSystemRequestJSON := `{"storage_array_type":"powerflex", "unique_id":"1", "username":"admin", "password":"password", "management_endpoint":"http://localhost"}`
			createStorageSystemResponseJSON := `{"id":0,"storage_array_type_id":0,"unique_id":"1","username":"admin","management_endpoint":"http://localhost"}`

			storageArrayStore := mocks.NewMockStorageArrayStoreInterface(ctrl)
			storageArrayStore.EXPECT().GetTypeByTypeName("powerflex").Times(1).Return(&model.StorageArrayType{Name: "powerflex"}, nil)
			storageArrayStore.EXPECT().Create(gomock.Any()).Times(1)
			handler := &StorageArrayHandler{storageArrayStore}
			return http.StatusCreated, handler, createStorageSystemRequestJSON, createStorageSystemResponseJSON, ctrl
		},
		"invalid payload": func(*testing.T) (int, *StorageArrayHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			createStorageSystemRequestJSON := `invalid-payload`

			storageArrayStore := mocks.NewMockStorageArrayStoreInterface(ctrl)
			handler := &StorageArrayHandler{storageArrayStore}
			return http.StatusUnprocessableEntity, handler, createStorageSystemRequestJSON, "", ctrl
		},
		"error getting array type": func(*testing.T) (int, *StorageArrayHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			createStorageSystemRequestJSON := `{"storage_array_type":"powerflex", "unique_id":"1", "username":"admin", "password":"password", "management_endpoint":"http://localhost"}`

			storageArrayStore := mocks.NewMockStorageArrayStoreInterface(ctrl)
			storageArrayStore.EXPECT().GetTypeByTypeName("powerflex").Times(1).Return(nil, errors.New("error"))
			handler := &StorageArrayHandler{storageArrayStore}
			return http.StatusUnprocessableEntity, handler, createStorageSystemRequestJSON, "", ctrl
		},
		"error persisting to database": func(*testing.T) (int, *StorageArrayHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			createStorageSystemRequestJSON := `{"storage_array_type":"powerflex", "unique_id":"1", "username":"admin", "password":"password", "management_endpoint":"http://localhost"}`

			storageArrayStore := mocks.NewMockStorageArrayStoreInterface(ctrl)
			storageArrayStore.EXPECT().GetTypeByTypeName("powerflex").Times(1).Return(&model.StorageArrayType{Name: "powerflex"}, nil)
			storageArrayStore.EXPECT().Create(gomock.Any()).Return(errors.New("error")).Times(1)
			handler := &StorageArrayHandler{storageArrayStore}
			return http.StatusInternalServerError, handler, createStorageSystemRequestJSON, "", ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			expectedStatus, handler, createRequest, expectedResponse, ctrl := tc(t)

			e := router.New()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(createRequest))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			assert.NoError(t, handler.CreateStorageArray(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}

func Test_UpdateStorageArray(t *testing.T) {

	tests := map[string]func(t *testing.T) (int, *StorageArrayHandler, string, string, *gomock.Controller){
		"success": func(*testing.T) (int, *StorageArrayHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			updateStorageSystemRequestJSON := `{"storage_array_type":"powerflex", "unique_id":"1", "username":"admin", "password":"password", "management_endpoint":"http://localhost"}`
			updateStorageSystemResponseJSON := `{"id":0,"storage_array_type_id":0,"unique_id":"1","username":"admin","management_endpoint":"http://localhost"}`

			storageArrayStore := mocks.NewMockStorageArrayStoreInterface(ctrl)
			storageArrayStore.EXPECT().GetTypeByTypeName("powerflex").Times(1).Return(&model.StorageArrayType{Name: "powerflex"}, nil)
			storageArrayStore.EXPECT().Update(gomock.Any()).Times(1)
			handler := &StorageArrayHandler{storageArrayStore}
			return http.StatusOK, handler, updateStorageSystemRequestJSON, updateStorageSystemResponseJSON, ctrl
		},
		"invalid request": func(*testing.T) (int, *StorageArrayHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			updateStorageSystemRequestJSON := `invalid-request`

			storageArrayStore := mocks.NewMockStorageArrayStoreInterface(ctrl)
			handler := &StorageArrayHandler{storageArrayStore}
			return http.StatusUnprocessableEntity, handler, updateStorageSystemRequestJSON, "", ctrl
		},
		"error getting array type": func(*testing.T) (int, *StorageArrayHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			updateStorageSystemRequestJSON := `{"storage_array_type":"powerflex", "unique_id":"1", "username":"admin", "password":"password", "management_endpoint":"http://localhost"}`

			storageArrayStore := mocks.NewMockStorageArrayStoreInterface(ctrl)
			storageArrayStore.EXPECT().GetTypeByTypeName("powerflex").Times(1).Return(nil, errors.New("error"))
			handler := &StorageArrayHandler{storageArrayStore}
			return http.StatusUnprocessableEntity, handler, updateStorageSystemRequestJSON, "", ctrl
		},
		"error persisting to database": func(*testing.T) (int, *StorageArrayHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			updateStorageSystemRequestJSON := `{"storage_array_type":"powerflex", "unique_id":"1", "username":"admin", "password":"password", "management_endpoint":"http://localhost"}`

			storageArrayStore := mocks.NewMockStorageArrayStoreInterface(ctrl)
			storageArrayStore.EXPECT().GetTypeByTypeName("powerflex").Times(1).Return(&model.StorageArrayType{Name: "powerflex"}, nil)
			storageArrayStore.EXPECT().Update(gomock.Any()).Return(errors.New("error")).Times(1)
			handler := &StorageArrayHandler{storageArrayStore}
			return http.StatusInternalServerError, handler, updateStorageSystemRequestJSON, "", ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			expectedStatus, handler, request, expectedResponse, ctrl := tc(t)

			e := router.New()
			req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(request))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/storage-arrays/:id")
			c.SetParamNames("id")
			c.SetParamValues("1")

			assert.NoError(t, handler.UpdateStorageArray(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}

var ()

func Test_GetStorageArray(t *testing.T) {

	tests := map[string]func(t *testing.T) (int, *StorageArrayHandler, string, string, *gomock.Controller){
		"success": func(*testing.T) (int, *StorageArrayHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			getStorageSystemResponseJSON := `{"id":1,"storage_array_type_id":1,"unique_id":"def321","username":"user","management_endpoint":"http://localhost:4321"}`

			storageArrayStore := mocks.NewMockStorageArrayStoreInterface(ctrl)
			storageArray := model.StorageArray{
				UniqueID:           "def321",
				Username:           "user",
				ManagementEndpoint: "http://localhost:4321",
				StorageArrayTypeID: 1,
			}
			storageArray.ID = 1
			storageArrayStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&storageArray, nil)
			handler := &StorageArrayHandler{storageArrayStore}
			return http.StatusOK, handler, "1", getStorageSystemResponseJSON, ctrl
		},
		"nil result from db": func(*testing.T) (int, *StorageArrayHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			storageArrayStore := mocks.NewMockStorageArrayStoreInterface(ctrl)
			storageArrayStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, nil)
			handler := &StorageArrayHandler{storageArrayStore}
			return http.StatusNotFound, handler, "1", "", ctrl
		},
		"error querying db": func(*testing.T) (int, *StorageArrayHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			storageArrayStore := mocks.NewMockStorageArrayStoreInterface(ctrl)
			storageArrayStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, errors.New("error"))
			handler := &StorageArrayHandler{storageArrayStore}
			return http.StatusInternalServerError, handler, "1", "", ctrl
		},
		"id is not numeric": func(*testing.T) (int, *StorageArrayHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			storageArrayStore := mocks.NewMockStorageArrayStoreInterface(ctrl)
			handler := &StorageArrayHandler{storageArrayStore}
			return http.StatusUnprocessableEntity, handler, "abc", "", ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			expectedStatus, handler, storageSystemID, expectedResponse, ctrl := tc(t)

			e := router.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/storage-arrays/:id")
			c.SetParamNames("id")
			c.SetParamValues(storageSystemID)

			assert.NoError(t, handler.GetStorageArray(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}

func Test_ListStorageArrays(t *testing.T) {

	tests := map[string]func(t *testing.T) (int, *StorageArrayHandler, string, *gomock.Controller){
		"success": func(*testing.T) (int, *StorageArrayHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			listStorageSystemResponseJSON := "[{\"id\":0,\"storage_array_type_id\":1,\"unique_id\":\"abc123\",\"username\":\"admin\",\"management_endpoint\":\"http://localhost:1234\"},{\"id\":0,\"storage_array_type_id\":2,\"unique_id\":\"def321\",\"username\":\"user\",\"management_endpoint\":\"http://localhost:4321\"}]"

			storageArrayStore := mocks.NewMockStorageArrayStoreInterface(ctrl)

			storageArrays := make([]model.StorageArray, 0)
			storageArrays = append(storageArrays, model.StorageArray{
				UniqueID:           "abc123",
				Username:           "admin",
				ManagementEndpoint: "http://localhost:1234",
				StorageArrayTypeID: 1,
			})
			storageArrays = append(storageArrays, model.StorageArray{
				UniqueID:           "def321",
				Username:           "user",
				ManagementEndpoint: "http://localhost:4321",
				StorageArrayTypeID: 2,
			})
			storageArrayStore.EXPECT().GetAll().Times(1).Return(storageArrays, nil)
			handler := &StorageArrayHandler{storageArrayStore}
			return http.StatusOK, handler, listStorageSystemResponseJSON, ctrl
		},
		"error querying database": func(*testing.T) (int, *StorageArrayHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			storageArrayStore := mocks.NewMockStorageArrayStoreInterface(ctrl)
			storageArrayStore.EXPECT().GetAll().Times(1).Return(nil, errors.New("error"))
			handler := &StorageArrayHandler{storageArrayStore}
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

			assert.NoError(t, handler.ListStorageArrays(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}

func Test_DeleteStorageArray(t *testing.T) {

	tests := map[string]func(t *testing.T) (int, *StorageArrayHandler, string, *gomock.Controller){
		"success": func(*testing.T) (int, *StorageArrayHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			storageArrayStore := mocks.NewMockStorageArrayStoreInterface(ctrl)
			storageArrayStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&model.StorageArray{
				UniqueID:           "def321",
				Username:           "user",
				ManagementEndpoint: "http://localhost:4321",
				StorageArrayTypeID: 2,
			}, nil)
			storageArrayStore.EXPECT().Delete(gomock.Any()).Times(1)
			handler := &StorageArrayHandler{storageArrayStore}
			return http.StatusOK, handler, "1", ctrl
		},
		"nil result from db": func(*testing.T) (int, *StorageArrayHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			storageArrayStore := mocks.NewMockStorageArrayStoreInterface(ctrl)
			storageArrayStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, nil)
			handler := &StorageArrayHandler{storageArrayStore}
			return http.StatusNotFound, handler, "1", ctrl
		},
		"error getting from db": func(*testing.T) (int, *StorageArrayHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			storageArrayStore := mocks.NewMockStorageArrayStoreInterface(ctrl)
			storageArrayStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, errors.New("error"))
			handler := &StorageArrayHandler{storageArrayStore}
			return http.StatusInternalServerError, handler, "1", ctrl
		},
		"error deleting from db": func(*testing.T) (int, *StorageArrayHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			storageArrayStore := mocks.NewMockStorageArrayStoreInterface(ctrl)
			storageArrayStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&model.StorageArray{
				UniqueID:           "def321",
				Username:           "user",
				ManagementEndpoint: "http://localhost:4321",
				StorageArrayTypeID: 2,
			}, nil)
			storageArrayStore.EXPECT().Delete(gomock.Any()).Times(1).Return(errors.New("error"))
			handler := &StorageArrayHandler{storageArrayStore}
			return http.StatusInternalServerError, handler, "1", ctrl
		},
		"id is not numeric": func(*testing.T) (int, *StorageArrayHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			storageArrayStore := mocks.NewMockStorageArrayStoreInterface(ctrl)

			handler := &StorageArrayHandler{storageArrayStore}
			return http.StatusUnprocessableEntity, handler, "abc", ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			expectedStatus, handler, storageSystemID, ctrl := tc(t)

			e := router.New()
			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/storage-arrays/:id")
			c.SetParamNames("id")
			c.SetParamValues(storageSystemID)

			assert.NoError(t, handler.DeleteStorageArray(c))
			assert.Equal(t, expectedStatus, rec.Code)
			ctrl.Finish()
		})
	}
}
