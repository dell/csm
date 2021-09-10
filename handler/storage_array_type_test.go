// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

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

func Test_StorageArrayTypeHandlerRegister(t *testing.T) {
	storageArrayTypeHandler := &StorageArrayTypeHandler{}
	rt := router.New()
	api := rt.Group("/api/v1")
	storageArrayTypeHandler.Register(api)
}
func Test_GetStorageArrayType(t *testing.T) {
	tests := map[string]func(t *testing.T) (int, *StorageArrayTypeHandler, string, string, *gomock.Controller){
		"success": func(*testing.T) (int, *StorageArrayTypeHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			getStorageArrayTypeResponseJSON := `{"id":"1","name":"type-1"}`

			storageArrayTypeStore := mocks.NewMockStorageArrayTypeStoreInterface(ctrl)
			storgeArrayType := model.StorageArrayType{
				Name: "type-1",
			}
			storgeArrayType.ID = 1

			storageArrayTypeStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&storgeArrayType, nil)
			handler := &StorageArrayTypeHandler{storageArrayTypeStore}
			return http.StatusOK, handler, "1", getStorageArrayTypeResponseJSON, ctrl
		},
		"nil result from db": func(*testing.T) (int, *StorageArrayTypeHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			storageArrayTypeStore := mocks.NewMockStorageArrayTypeStoreInterface(ctrl)
			storageArrayTypeStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, nil)
			handler := &StorageArrayTypeHandler{storageArrayTypeStore}
			return http.StatusNotFound, handler, "1", "", ctrl
		},
		"error querying db": func(*testing.T) (int, *StorageArrayTypeHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			storageArrayTypeStore := mocks.NewMockStorageArrayTypeStoreInterface(ctrl)
			storageArrayTypeStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, errors.New("error"))
			handler := &StorageArrayTypeHandler{storageArrayTypeStore}
			return http.StatusInternalServerError, handler, "1", "", ctrl
		},
		"id is not numeric": func(*testing.T) (int, *StorageArrayTypeHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			storageArrayTypeStore := mocks.NewMockStorageArrayTypeStoreInterface(ctrl)
			handler := &StorageArrayTypeHandler{storageArrayTypeStore}
			return http.StatusUnprocessableEntity, handler, "abc", "", ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			expectedStatus, handler, storageArrayTypeID, expectedResponse, ctrl := tc(t)

			e := router.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/storage-array-types/:id")
			c.SetParamNames("id")
			c.SetParamValues(storageArrayTypeID)

			assert.NoError(t, handler.GetStorageArrayType(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}
func Test_ListStorageArrayType(t *testing.T) {
	tests := map[string]func(t *testing.T) (int, *StorageArrayTypeHandler, string, *gomock.Controller){
		"success": func(*testing.T) (int, *StorageArrayTypeHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			listStorageArrayTypeResponseJSON := `[{"id":"0","name":"type-1"},{"id":"0","name":"type-2"}]`

			storageArrayTypeStore := mocks.NewMockStorageArrayTypeStoreInterface(ctrl)

			storageArrayTypes := make([]model.StorageArrayType, 0)
			storageArrayTypes = append(storageArrayTypes, model.StorageArrayType{
				Name: "type-1",
			})
			storageArrayTypes = append(storageArrayTypes, model.StorageArrayType{
				Name: "type-2",
			})
			storageArrayTypeStore.EXPECT().GetAll().Times(1).Return(storageArrayTypes, nil)

			handler := &StorageArrayTypeHandler{storageArrayTypeStore}

			return http.StatusOK, handler, listStorageArrayTypeResponseJSON, ctrl
		},
		"error querying database": func(*testing.T) (int, *StorageArrayTypeHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			storageArrayTypeStore := mocks.NewMockStorageArrayTypeStoreInterface(ctrl)
			storageArrayTypeStore.EXPECT().GetAll().Times(1).Return(nil, errors.New("error"))
			handler := &StorageArrayTypeHandler{storageArrayTypeStore}

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

			assert.NoError(t, handler.ListStorageArrayTypes(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}
