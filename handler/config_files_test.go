// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package handler

import (
	"bytes"
	"errors"
	"mime/multipart"
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

func Test_ConfigFileHandlerRegister(t *testing.T) {
	configFileHandler := &ConfigFileHandler{}
	rt := router.New()
	api := rt.Group("/api/v1")
	configFileHandler.Register(api)
}

func Test_CreateConfigFile(t *testing.T) {
	tests := map[string]func(t *testing.T) (int, *ConfigFileHandler, *bytes.Buffer, *multipart.Writer, string, *gomock.Controller){
		"success": func(*testing.T) (int, *ConfigFileHandler, *bytes.Buffer, *multipart.Writer, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			createConfigFileResponse := `{"id":"0","name":"abc"}`

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			writer.WriteField("name", "abc")
			part, _ := writer.CreateFormFile("file", "file.csv")
			part.Write([]byte(`sample random file contents`))
			writer.Close()

			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)
			configFileStore.EXPECT().Create(gomock.Any()).Times(1).Return(nil)

			handler := &ConfigFileHandler{configFileStore: configFileStore}
			return http.StatusCreated, handler, body, writer, createConfigFileResponse, ctrl
		},
		"error saving to database": func(*testing.T) (int, *ConfigFileHandler, *bytes.Buffer, *multipart.Writer, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			writer.WriteField("name", "abc")
			part, _ := writer.CreateFormFile("file", "file.csv")
			part.Write([]byte(`sample random file contents`))
			writer.Close()

			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)
			configFileStore.EXPECT().Create(gomock.Any()).Times(1).Return(errors.New("error saving to database"))

			handler := &ConfigFileHandler{configFileStore: configFileStore}
			return http.StatusInternalServerError, handler, body, writer, "", ctrl
		},
		"error empty cluster name": func(*testing.T) (int, *ConfigFileHandler, *bytes.Buffer, *multipart.Writer, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			writer.WriteField("name", "")
			part, _ := writer.CreateFormFile("file", "file.csv")
			part.Write([]byte(`sample random file contents`))
			writer.Close()

			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)
			handler := &ConfigFileHandler{configFileStore: configFileStore}
			return http.StatusBadRequest, handler, body, writer, "", ctrl
		},
		"error empty file upload": func(*testing.T) (int, *ConfigFileHandler, *bytes.Buffer, *multipart.Writer, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			writer.WriteField("name", "abc")
			writer.Close()

			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)
			handler := &ConfigFileHandler{configFileStore: configFileStore}
			return http.StatusBadRequest, handler, body, writer, "", ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			expectedStatus, handler, body, writer, expectedResponse, ctrl := tc(t)

			e := router.New()
			req := httptest.NewRequest(http.MethodPost, "/", body)
			req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			assert.NoError(t, handler.CreateConfigFile(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}

func Test_UpdateConfigFile(t *testing.T) {
	tests := map[string]func(t *testing.T) (int, *ConfigFileHandler, *bytes.Buffer, *multipart.Writer, string, string, *gomock.Controller){
		"success": func(*testing.T) (int, *ConfigFileHandler, *bytes.Buffer, *multipart.Writer, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			createConfigFileResponse := `null`

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			writer.WriteField("name", "new-file-name")
			part, _ := writer.CreateFormFile("file", "file.csv")
			part.Write([]byte(`sample random file contents`))
			writer.Close()

			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)

			configFileStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&model.ConfigFile{Name: "old-file-name"}, nil)
			configFileStore.EXPECT().Update(gomock.Any()).Times(1).Return(nil)

			handler := &ConfigFileHandler{configFileStore: configFileStore}
			return http.StatusNoContent, handler, body, writer, "1", createConfigFileResponse, ctrl
		},
		"error looking up cluster in database": func(*testing.T) (int, *ConfigFileHandler, *bytes.Buffer, *multipart.Writer, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			writer.WriteField("name", "new-file-name")
			part, _ := writer.CreateFormFile("file", "file.csv")
			part.Write([]byte(`sample random file contents`))
			writer.Close()

			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)
			configFileStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, errors.New("cluster not found"))

			handler := &ConfigFileHandler{configFileStore: configFileStore}
			return http.StatusInternalServerError, handler, body, writer, "1", "", ctrl
		},
		"error cluster not found in database": func(*testing.T) (int, *ConfigFileHandler, *bytes.Buffer, *multipart.Writer, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			writer.WriteField("name", "new-file-name")
			part, _ := writer.CreateFormFile("file", "file.csv")
			part.Write([]byte(`sample random file contents`))
			writer.Close()

			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)
			configFileStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, nil)

			handler := &ConfigFileHandler{configFileStore: configFileStore}
			return http.StatusNotFound, handler, body, writer, "1", "", ctrl
		},
		"error saving to database": func(*testing.T) (int, *ConfigFileHandler, *bytes.Buffer, *multipart.Writer, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			writer.WriteField("name", "abc")
			part, _ := writer.CreateFormFile("file", "file.csv")
			part.Write([]byte(`sample random file contents`))
			writer.Close()

			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)
			configFileStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&model.ConfigFile{Name: "abc"}, nil)
			configFileStore.EXPECT().Update(gomock.Any()).Times(1).Return(errors.New("error saving to database"))

			handler := &ConfigFileHandler{configFileStore: configFileStore}
			return http.StatusInternalServerError, handler, body, writer, "1", "", ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			expectedStatus, handler, body, writer, cfID, expectedResponse, ctrl := tc(t)

			e := router.New()
			req := httptest.NewRequest(http.MethodPatch, "/", body)
			req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/configuration-files/:id")
			c.SetParamNames("id")
			c.SetParamValues(cfID)

			assert.NoError(t, handler.UpdateConfigFile(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}

func Test_GetConfigFile(t *testing.T) {
	tests := map[string]func(t *testing.T) (int, *ConfigFileHandler, string, string, *gomock.Controller){
		"success": func(*testing.T) (int, *ConfigFileHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			getStorageSystemResponseJSON := `{"id":"1024","name":"file-name-1"}`

			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)
			cf := model.ConfigFile{
				Name: "file-name-1",
			}
			cf.ID = 1024
			configFileStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&cf, nil)

			handler := &ConfigFileHandler{configFileStore: configFileStore}
			return http.StatusOK, handler, "1", getStorageSystemResponseJSON, ctrl
		},
		"nil result from db": func(*testing.T) (int, *ConfigFileHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)
			configFileStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, nil)

			handler := &ConfigFileHandler{configFileStore: configFileStore}
			return http.StatusNotFound, handler, "1", "", ctrl
		},
		"error querying db": func(*testing.T) (int, *ConfigFileHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)
			configFileStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, errors.New("error"))

			handler := &ConfigFileHandler{configFileStore: configFileStore}
			return http.StatusInternalServerError, handler, "1", "", ctrl
		},
		"id is not numeric": func(*testing.T) (int, *ConfigFileHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)

			handler := &ConfigFileHandler{configFileStore: configFileStore}
			return http.StatusUnprocessableEntity, handler, "abc", "", ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			expectedStatus, handler, cfID, expectedResponse, ctrl := tc(t)

			e := router.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/configuration-files/:id")
			c.SetParamNames("id")
			c.SetParamValues(cfID)

			assert.NoError(t, handler.GetConfigFile(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}

func Test_ListConfigFiles(t *testing.T) {
	tests := map[string]func(t *testing.T) (int, *ConfigFileHandler, string, map[string]string, *gomock.Controller){
		"success": func(*testing.T) (int, *ConfigFileHandler, string, map[string]string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			listconfigFileResponseJSON := "[{\"id\":\"0\",\"name\":\"file-name-1\"},{\"id\":\"0\",\"name\":\"file-name-2\"}]"

			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)

			cfs := make([]model.ConfigFile, 0)
			cfs = append(cfs, model.ConfigFile{
				Name: "file-name-1",
			})
			cfs = append(cfs, model.ConfigFile{
				Name: "file-name-2",
			})
			configFileStore.EXPECT().GetAll().Times(1).Return(cfs, nil)

			handler := &ConfigFileHandler{configFileStore: configFileStore}
			return http.StatusOK, handler, listconfigFileResponseJSON, nil, ctrl
		},
		"success listing by config_name": func(*testing.T) (int, *ConfigFileHandler, string, map[string]string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			listStorageSystemResponseJSON := "[{\"id\":\"0\",\"name\":\"file-name-1\"}]"

			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)

			cfs := make([]model.ConfigFile, 0)
			cfs = append(cfs, model.ConfigFile{
				Name: "file-name-1",
			})

			configFileStore.EXPECT().GetAllByName(gomock.Any()).Times(1).Return(cfs, nil)

			handler := &ConfigFileHandler{configFileStore: configFileStore}
			return http.StatusOK, handler, listStorageSystemResponseJSON, map[string]string{"config_name": "file-name-1"}, ctrl
		},
		"error querying database": func(*testing.T) (int, *ConfigFileHandler, string, map[string]string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)
			configFileStore.EXPECT().GetAll().Times(1).Return(nil, errors.New("error"))
			handler := &ConfigFileHandler{configFileStore: configFileStore}
			return http.StatusInternalServerError, handler, "", nil, ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			expectedStatus, handler, expectedResponse, queryParams, ctrl := tc(t)

			e := router.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			q := req.URL.Query()
			for key, value := range queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			assert.NoError(t, handler.ListConfigFiles(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}

func Test_DeleteConfigFile(t *testing.T) {
	tests := map[string]func(t *testing.T) (int, *ConfigFileHandler, string, *gomock.Controller){
		"success": func(*testing.T) (int, *ConfigFileHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)
			configFileStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&model.ConfigFile{}, nil)
			configFileStore.EXPECT().Delete(gomock.Any()).Times(1)
			handler := &ConfigFileHandler{configFileStore: configFileStore}
			return http.StatusNoContent, handler, "1", ctrl
		},
		"nil result from db": func(*testing.T) (int, *ConfigFileHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)
			configFileStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, nil)
			handler := &ConfigFileHandler{configFileStore: configFileStore}
			return http.StatusNotFound, handler, "1", ctrl
		},
		"error getting from db": func(*testing.T) (int, *ConfigFileHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)
			configFileStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, errors.New("error"))
			handler := &ConfigFileHandler{configFileStore: configFileStore}
			return http.StatusInternalServerError, handler, "1", ctrl
		},
		"error deleting from db": func(*testing.T) (int, *ConfigFileHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)
			configFileStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&model.ConfigFile{}, nil)
			configFileStore.EXPECT().Delete(gomock.Any()).Times(1).Return(errors.New("error"))
			handler := &ConfigFileHandler{configFileStore: configFileStore}
			return http.StatusInternalServerError, handler, "1", ctrl
		},
		"id is not numeric": func(*testing.T) (int, *ConfigFileHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)
			handler := &ConfigFileHandler{configFileStore: configFileStore}
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
			c.SetPath("/configuration-files/:id")
			c.SetParamNames("id")
			c.SetParamValues(storageSystemID)

			assert.NoError(t, handler.DeleteConfigFile(c))
			assert.Equal(t, expectedStatus, rec.Code)
			ctrl.Finish()
		})
	}
}
