// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	handlerMocks "github.com/dell/csm-deployment/handler/mocks"
	"github.com/dell/csm-deployment/k8s"
	k8sMocks "github.com/dell/csm-deployment/k8s/mocks"
	kappMocks "github.com/dell/csm-deployment/kapp/mocks"
	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/prechecks"
	"github.com/dell/csm-deployment/router"
	"github.com/dell/csm-deployment/store/mocks"
	"github.com/dell/csm-deployment/ytt"
	yttMocks "github.com/dell/csm-deployment/ytt/mocks"
	"github.com/golang/mock/gomock"
	"github.com/k14s/ytt/pkg/cmd/template"
	"github.com/k14s/ytt/pkg/files"
	"github.com/labstack/echo/v4"

	"github.com/stretchr/testify/assert"
)

func Test_ApplicationHandlerRegister(t *testing.T) {
	applicationHandler := &ApplicationHandler{}
	rt := router.New()
	api := rt.Group("/api/v1")
	applicationHandler.Register(api)
}

func Test_CreateApplication(t *testing.T) {

	tests := map[string]func(t *testing.T) (int, *ApplicationHandler, string, string, *gomock.Controller){
		"success": func(*testing.T) (int, *ApplicationHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			request := `{"name": "app1", "cluster_id": "1", "driver_type_id": "2", "module_types": ["1"], "storage_arrays": ["1"]}`
			response := ""

			applicationStore := mocks.NewMockApplicationStoreInterface(ctrl)
			arrayStore := mocks.NewMockStorageArrayStoreInterface(ctrl)
			moduleStore := mocks.NewMockModuleTypeStoreInterface(ctrl)
			taskStore := mocks.NewMockTaskStoreInterface(ctrl)
			precheckHandler := handlerMocks.NewMockPrecheckHandlerInterface(ctrl)

			applicationStore.EXPECT().Create(gomock.Any()).Times(1)
			arrayStore.EXPECT().GetAllByID(gomock.Any()).Times(1).Return([]model.StorageArray{}, nil)
			moduleStore.EXPECT().GetAllByID(gomock.Any()).Times(1).Return([]model.ModuleType{{Name: "sample-module"}}, nil)
			moduleStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&model.ModuleType{}, nil)
			taskStore.EXPECT().Create(gomock.Any()).Times(1).Return(nil)
			precheckHandler.EXPECT().Precheck(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

			handler := &ApplicationHandler{
				applicationStore: applicationStore,
				arrayStore:       arrayStore,
				ModuleTypeStore:  moduleStore,
				taskStore:        taskStore,
				precheckHandler:  precheckHandler,
				SkipGoRoutine:    true,
			}
			return http.StatusAccepted, handler, request, response, ctrl
		},
		"error with no configuration for replication module": func(*testing.T) (int, *ApplicationHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			request := `{"name": "app1", "cluster_id": "1", "driver_type_id": "2", "module_types": ["1"], "storage_arrays": ["1"]}`
			response := ""

			applicationStore := mocks.NewMockApplicationStoreInterface(ctrl)
			arrayStore := mocks.NewMockStorageArrayStoreInterface(ctrl)
			moduleStore := mocks.NewMockModuleTypeStoreInterface(ctrl)
			taskStore := mocks.NewMockTaskStoreInterface(ctrl)
			precheckHandler := handlerMocks.NewMockPrecheckHandlerInterface(ctrl)

			arrayStore.EXPECT().GetAllByID(gomock.Any()).Times(1).Return([]model.StorageArray{}, nil)
			moduleStore.EXPECT().GetAllByID(gomock.Any()).Times(1).Return([]model.ModuleType{{Name: model.ModuleTypeReplication}}, nil)
			moduleStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&model.ModuleType{Name: model.ModuleTypeReplication}, nil)

			handler := &ApplicationHandler{
				applicationStore: applicationStore,
				arrayStore:       arrayStore,
				ModuleTypeStore:  moduleStore,
				taskStore:        taskStore,
				precheckHandler:  precheckHandler,
				SkipGoRoutine:    true,
			}
			return http.StatusUnprocessableEntity, handler, request, response, ctrl
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

			assert.NoError(t, handler.CreateApplication(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}

func Test_CreateApplicationWithReplication(t *testing.T) {

	tests := map[string]func(t *testing.T) (int, *ApplicationHandler, string, string, *gomock.Controller){
		"success": func(*testing.T) (int, *ApplicationHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			request := `{"name": "app1", "cluster_id": "1", "driver_type_id": "2", "storage_arrays": ["1"], "module_configuration": ["target_cluster:2"], "module_types":["4"]}`
			response := ""

			module := model.ModuleType{
				Name: "replication",
			}
			module.ID = 4

			applicationStore := mocks.NewMockApplicationStoreInterface(ctrl)
			arrayStore := mocks.NewMockStorageArrayStoreInterface(ctrl)
			moduleStore := mocks.NewMockModuleTypeStoreInterface(ctrl)
			taskStore := mocks.NewMockTaskStoreInterface(ctrl)
			precheckHandler := handlerMocks.NewMockPrecheckHandlerInterface(ctrl)

			applicationStore.EXPECT().Create(gomock.Any()).Times(1)
			arrayStore.EXPECT().GetAllByID(gomock.Any()).Times(1).Return([]model.StorageArray{}, nil)
			moduleStore.EXPECT().GetAllByID(gomock.Any()).Times(1).Return([]model.ModuleType{}, nil)
			moduleStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&module, nil)
			taskStore.EXPECT().Create(gomock.Any()).Times(1).Return(nil)
			precheckHandler.EXPECT().Precheck(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

			handler := &ApplicationHandler{
				applicationStore: applicationStore,
				arrayStore:       arrayStore,
				ModuleTypeStore:  moduleStore,
				taskStore:        taskStore,
				precheckHandler:  precheckHandler,
				SkipGoRoutine:    true,
			}
			return http.StatusAccepted, handler, request, response, ctrl
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

			assert.NoError(t, handler.CreateApplication(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}

func Test_GetApplication(t *testing.T) {

	tests := map[string]func(t *testing.T) (int, *ApplicationHandler, string, string, *gomock.Controller){
		"success": func(*testing.T) (int, *ApplicationHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			getStorageSystemResponseJSON := `{"id":"1","name":"app-1","cluster_id":"0","driver_type_id":"0","module_types":null,"storage_arrays":null,"driver_configuration":[""],"module_configuration":[""],"application_output":""}`

			applicationStore := mocks.NewMockApplicationStoreInterface(ctrl)
			application := model.Application{
				Name: "app-1",
			}
			application.ID = 1
			applicationStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&application, nil)
			handler := &ApplicationHandler{applicationStore: applicationStore}
			return http.StatusOK, handler, "1", getStorageSystemResponseJSON, ctrl
		},
		"nil result from db": func(*testing.T) (int, *ApplicationHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			applicationStore := mocks.NewMockApplicationStoreInterface(ctrl)
			applicationStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, nil)
			handler := &ApplicationHandler{applicationStore: applicationStore}
			return http.StatusNotFound, handler, "1", "", ctrl
		},
		"error querying db": func(*testing.T) (int, *ApplicationHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			applicationStore := mocks.NewMockApplicationStoreInterface(ctrl)
			applicationStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, errors.New("error"))
			handler := &ApplicationHandler{applicationStore: applicationStore}
			return http.StatusInternalServerError, handler, "1", "", ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			expectedStatus, handler, clusterID, expectedResponse, ctrl := tc(t)

			e := router.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/clusters/:id")
			c.SetParamNames("id")
			c.SetParamValues(clusterID)

			assert.NoError(t, handler.GetApplication(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}

func Test_ListApplications(t *testing.T) {

	tests := map[string]func(t *testing.T) (int, *ApplicationHandler, string, map[string]string, *gomock.Controller){
		"success": func(*testing.T) (int, *ApplicationHandler, string, map[string]string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			response := "[{\"id\":\"0\",\"name\":\"app-1\",\"cluster_id\":\"5\",\"driver_type_id\":\"3\",\"module_types\":null,\"storage_arrays\":null,\"driver_configuration\":[\"\"],\"module_configuration\":[\"\"],\"application_output\":\"\"},{\"id\":\"0\",\"name\":\"app-2\",\"cluster_id\":\"6\",\"driver_type_id\":\"4\",\"module_types\":null,\"storage_arrays\":null,\"driver_configuration\":[\"\"],\"module_configuration\":[\"\"],\"application_output\":\"\"}]"

			applicationStore := mocks.NewMockApplicationStoreInterface(ctrl)

			applications := make([]model.Application, 0)
			applications = append(applications, model.Application{
				Name:         "app-1",
				ClusterID:    5,
				DriverTypeID: 3,
			})
			applications = append(applications, model.Application{
				Name:         "app-2",
				ClusterID:    6,
				DriverTypeID: 4,
			})
			applicationStore.EXPECT().GetAll().Times(1).Return(applications, nil)
			handler := &ApplicationHandler{applicationStore: applicationStore}
			return http.StatusOK, handler, response, nil, ctrl
		},
		"success listing by name": func(*testing.T) (int, *ApplicationHandler, string, map[string]string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			response := "[{\"id\":\"0\",\"name\":\"app-1\",\"cluster_id\":\"5\",\"driver_type_id\":\"3\",\"module_types\":null,\"storage_arrays\":null,\"driver_configuration\":[\"\"],\"module_configuration\":[\"\"],\"application_output\":\"\"}]"

			applicationStore := mocks.NewMockApplicationStoreInterface(ctrl)

			application := model.Application{
				Name:         "app-1",
				ClusterID:    5,
				DriverTypeID: 3,
			}
			applicationStore.EXPECT().GetByName(gomock.Any()).Times(1).Return(&application, nil)
			handler := &ApplicationHandler{applicationStore: applicationStore}
			return http.StatusOK, handler, response, map[string]string{"name": "app-1"}, ctrl
		},
		"error querying database": func(*testing.T) (int, *ApplicationHandler, string, map[string]string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			applicationStore := mocks.NewMockApplicationStoreInterface(ctrl)
			applicationStore.EXPECT().GetAll().Times(1).Return(nil, errors.New("error"))
			handler := &ApplicationHandler{applicationStore: applicationStore}
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

			assert.NoError(t, handler.ListApplications(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}

func Test_DeleteApplication(t *testing.T) {

	tests := map[string]func(t *testing.T) (int, *ApplicationHandler, string, *gomock.Controller){
		"success": func(*testing.T) (int, *ApplicationHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			clusterStore.EXPECT().GetByID(gomock.Any()).Return(&model.Cluster{ConfigFileData: []byte("config")}, nil)

			kappClient := kappMocks.NewMockInterface(ctrl)
			kappClient.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil)

			runtimeClient := k8sMocks.NewMockControllerRuntimeInterface(ctrl)
			runtimeClient.EXPECT().DeleteNameSpaceByName(gomock.Any(), gomock.Any()).Return(nil)

			applicationStore := mocks.NewMockApplicationStoreInterface(ctrl)
			application := model.Application{
				Name: "app-1",
			}
			application.ID = 1
			applicationStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&application, nil)
			applicationStore.EXPECT().Delete(gomock.Any()).Times(1).Return(nil)
			handler := &ApplicationHandler{
				applicationStore: applicationStore,
				clusterStore:     clusterStore,
				kappClient:       kappClient,
				runtimeClientFunc: func(data []byte, logger echo.Logger) (k8s.ControllerRuntimeInterface, error) {
					return runtimeClient, nil
				},
			}

			return http.StatusNoContent, handler, "1", ctrl
		},
		"error deleting from the database": func(*testing.T) (int, *ApplicationHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			clusterStore.EXPECT().GetByID(gomock.Any()).Return(&model.Cluster{ConfigFileData: []byte("config")}, nil)

			kappClient := kappMocks.NewMockInterface(ctrl)
			kappClient.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil)

			runtimeClient := k8sMocks.NewMockControllerRuntimeInterface(ctrl)

			applicationStore := mocks.NewMockApplicationStoreInterface(ctrl)
			application := model.Application{
				Name: "app-1",
			}
			application.ID = 1
			applicationStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&application, nil)
			applicationStore.EXPECT().Delete(gomock.Any()).Times(1).Return(errors.New("error"))
			handler := &ApplicationHandler{
				applicationStore: applicationStore,
				clusterStore:     clusterStore,
				kappClient:       kappClient,
				runtimeClientFunc: func(data []byte, logger echo.Logger) (k8s.ControllerRuntimeInterface, error) {
					return runtimeClient, nil
				},
			}
			return http.StatusInternalServerError, handler, "1", ctrl
		},
		"nil result from db": func(*testing.T) (int, *ApplicationHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			applicationStore := mocks.NewMockApplicationStoreInterface(ctrl)
			applicationStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, nil)
			handler := &ApplicationHandler{applicationStore: applicationStore}
			return http.StatusNotFound, handler, "1", ctrl
		},
		"error querying db": func(*testing.T) (int, *ApplicationHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			applicationStore := mocks.NewMockApplicationStoreInterface(ctrl)
			applicationStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, errors.New("error"))
			handler := &ApplicationHandler{applicationStore: applicationStore}
			return http.StatusInternalServerError, handler, "1", ctrl
		},
		"error getting cluster from db": func(*testing.T) (int, *ApplicationHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			clusterStore.EXPECT().GetByID(gomock.Any()).Return(nil, errors.New("error"))

			applicationStore := mocks.NewMockApplicationStoreInterface(ctrl)
			application := model.Application{
				Name: "app-1",
			}
			application.ID = 1
			applicationStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&application, nil)
			handler := &ApplicationHandler{
				applicationStore: applicationStore,
				clusterStore:     clusterStore,
			}
			return http.StatusInternalServerError, handler, "1", ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			expectedStatus, handler, clusterID, ctrl := tc(t)

			e := router.New()
			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/clusters/:id")
			c.SetParamNames("id")
			c.SetParamValues(clusterID)

			assert.NoError(t, handler.DeleteApplication(c))
			assert.Equal(t, expectedStatus, rec.Code)
			ctrl.Finish()
		})
	}
}

type EmptyValidator struct{}

func (e EmptyValidator) Validate() error {
	return nil
}

func Test_PrecheckHandler(t *testing.T) {

	tests := map[string]func(t *testing.T) (bool, *PrecheckHandler, *gomock.Controller){

		"success": func(*testing.T) (bool, *PrecheckHandler, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			driverStore := mocks.NewMockDriverTypeStoreInterface(ctrl)
			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)
			precheckGetter := handlerMocks.NewMockPrecheckGetterInterface(ctrl)

			clusterStore.EXPECT().GetByID(gomock.Any()).Return(&model.Cluster{}, nil)
			configFileStore.EXPECT().GetAll().Return([]model.ConfigFile{}, nil)
			driverStore.EXPECT().GetByID(gomock.Any()).Return(&model.DriverType{StorageArrayType: model.StorageArrayType{Name: "powerflex"}}, nil)
			precheckGetter.EXPECT().GetDriverPrechecks(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]prechecks.Validator{EmptyValidator{}})
			precheckGetter.EXPECT().GetModuleTypePrechecks(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]prechecks.Validator{EmptyValidator{}})

			handler := &PrecheckHandler{
				driverStore:     driverStore,
				clusterStore:    clusterStore,
				configFileStore: configFileStore,
				precheckGetter:  precheckGetter,
			}
			return false, handler, ctrl
		},
		"error getting cluster": func(*testing.T) (bool, *PrecheckHandler, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			driverStore := mocks.NewMockDriverTypeStoreInterface(ctrl)
			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)
			precheckGetter := handlerMocks.NewMockPrecheckGetterInterface(ctrl)

			clusterStore.EXPECT().GetByID(gomock.Any()).Return(&model.Cluster{}, errors.New("error"))

			handler := &PrecheckHandler{
				driverStore:     driverStore,
				clusterStore:    clusterStore,
				configFileStore: configFileStore,
				precheckGetter:  precheckGetter,
			}
			return true, handler, ctrl
		},
		"error can't find cluster": func(*testing.T) (bool, *PrecheckHandler, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			driverStore := mocks.NewMockDriverTypeStoreInterface(ctrl)
			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)
			precheckGetter := handlerMocks.NewMockPrecheckGetterInterface(ctrl)

			clusterStore.EXPECT().GetByID(gomock.Any()).Return(nil, nil)

			handler := &PrecheckHandler{
				driverStore:     driverStore,
				clusterStore:    clusterStore,
				configFileStore: configFileStore,
				precheckGetter:  precheckGetter,
			}
			return true, handler, ctrl
		},
		"error getting config files": func(*testing.T) (bool, *PrecheckHandler, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			driverStore := mocks.NewMockDriverTypeStoreInterface(ctrl)
			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)
			precheckGetter := handlerMocks.NewMockPrecheckGetterInterface(ctrl)

			clusterStore.EXPECT().GetByID(gomock.Any()).Return(&model.Cluster{}, nil)
			configFileStore.EXPECT().GetAll().Return([]model.ConfigFile{}, errors.New("error"))

			handler := &PrecheckHandler{
				driverStore:     driverStore,
				clusterStore:    clusterStore,
				configFileStore: configFileStore,
				precheckGetter:  precheckGetter,
			}
			return true, handler, ctrl
		},
		"error getting driver": func(*testing.T) (bool, *PrecheckHandler, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			driverStore := mocks.NewMockDriverTypeStoreInterface(ctrl)
			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)
			precheckGetter := handlerMocks.NewMockPrecheckGetterInterface(ctrl)

			clusterStore.EXPECT().GetByID(gomock.Any()).Return(&model.Cluster{}, nil)
			configFileStore.EXPECT().GetAll().Return([]model.ConfigFile{}, nil)
			driverStore.EXPECT().GetByID(gomock.Any()).Return(&model.DriverType{}, errors.New("error"))

			handler := &PrecheckHandler{
				driverStore:     driverStore,
				clusterStore:    clusterStore,
				configFileStore: configFileStore,
				precheckGetter:  precheckGetter,
			}
			return true, handler, ctrl
		},
		"error finding driver": func(*testing.T) (bool, *PrecheckHandler, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			driverStore := mocks.NewMockDriverTypeStoreInterface(ctrl)
			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)
			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)
			precheckGetter := handlerMocks.NewMockPrecheckGetterInterface(ctrl)

			clusterStore.EXPECT().GetByID(gomock.Any()).Return(&model.Cluster{}, nil)
			configFileStore.EXPECT().GetAll().Return([]model.ConfigFile{}, nil)
			driverStore.EXPECT().GetByID(gomock.Any()).Return(nil, nil)

			handler := &PrecheckHandler{
				driverStore:     driverStore,
				clusterStore:    clusterStore,
				configFileStore: configFileStore,
				precheckGetter:  precheckGetter,
			}
			return true, handler, ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			expectError, handler, ctrl := tc(t)

			err := handler.Precheck(echo.New().NewContext(nil, nil), 1, 2, []model.ModuleType{{Name: "sample-module"}}, "")

			if expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			ctrl.Finish()
		})
	}
}

func Test_CreateReplicationNamespace(t *testing.T) {

	tests := map[string]func(t *testing.T) (*ApplicationHandler, *gomock.Controller){

		"success": func(*testing.T) (*ApplicationHandler, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			runtimeClient := k8sMocks.NewMockControllerRuntimeInterface(ctrl)
			runtimeClient.EXPECT().CreateNameSpaceFromName(gomock.Any(), gomock.Any()).Return(nil)

			handler := &ApplicationHandler{
				runtimeClientFunc: func(data []byte, logger echo.Logger) (k8s.ControllerRuntimeInterface, error) {
					return runtimeClient, nil
				},
			}
			return handler, ctrl
		},
		"error creating namespace": func(*testing.T) (*ApplicationHandler, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			runtimeClient := k8sMocks.NewMockControllerRuntimeInterface(ctrl)
			runtimeClient.EXPECT().CreateNameSpaceFromName(gomock.Any(), gomock.Any()).Return(errors.New("error"))

			handler := &ApplicationHandler{
				runtimeClientFunc: func(data []byte, logger echo.Logger) (k8s.ControllerRuntimeInterface, error) {
					return runtimeClient, nil
				},
			}
			return handler, ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			handler, ctrl := tc(t)
			handler.createReplicationNamespace([]byte{}, echo.New().NewContext(nil, nil))
			ctrl.Finish()
		})
	}
}

func Test_CreateReplicationSecrets(t *testing.T) {

	tests := map[string]func(t *testing.T) (*ApplicationHandler, *gomock.Controller){

		"success": func(*testing.T) (*ApplicationHandler, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			runtimeClient := k8sMocks.NewMockControllerRuntimeInterface(ctrl)
			runtimeClient.EXPECT().CreateNameSpaceFromName(gomock.Any(), gomock.Any()).Return(nil)
			runtimeClient.EXPECT().CreateSecretFromName(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			handler := &ApplicationHandler{
				runtimeClientFunc: func(data []byte, logger echo.Logger) (k8s.ControllerRuntimeInterface, error) {
					return runtimeClient, nil
				},
			}
			return handler, ctrl
		},
		"error creating namespace": func(*testing.T) (*ApplicationHandler, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			runtimeClient := k8sMocks.NewMockControllerRuntimeInterface(ctrl)
			runtimeClient.EXPECT().CreateNameSpaceFromName(gomock.Any(), gomock.Any()).Return(errors.New("error"))

			handler := &ApplicationHandler{
				runtimeClientFunc: func(data []byte, logger echo.Logger) (k8s.ControllerRuntimeInterface, error) {
					return runtimeClient, nil
				},
			}
			return handler, ctrl
		},
		"error creating secret": func(*testing.T) (*ApplicationHandler, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			runtimeClient := k8sMocks.NewMockControllerRuntimeInterface(ctrl)
			runtimeClient.EXPECT().CreateNameSpaceFromName(gomock.Any(), gomock.Any()).Return(nil)
			runtimeClient.EXPECT().CreateSecretFromName(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))

			handler := &ApplicationHandler{
				runtimeClientFunc: func(data []byte, logger echo.Logger) (k8s.ControllerRuntimeInterface, error) {
					return runtimeClient, nil
				},
			}
			return handler, ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			handler, ctrl := tc(t)
			handler.createReplicationSecrets([]byte{}, []byte{}, echo.New().NewContext(nil, nil))
			ctrl.Finish()
		})
	}
}

func Test_CreateReverseProxySecrets(t *testing.T) {

	tests := map[string]func(t *testing.T) (*ApplicationHandler, *gomock.Controller){
		"success": func(*testing.T) (*ApplicationHandler, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			runtimeClient := k8sMocks.NewMockControllerRuntimeInterface(ctrl)
			runtimeClient.EXPECT().CreateTLSSecretFromName(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)
			configFileStore.EXPECT().GetAllByName(gomock.Any()).Times(2).Return([]model.ConfigFile{{ConfigFileData: []byte("config")}}, nil)

			handler := &ApplicationHandler{
				configFileStore: configFileStore,
				runtimeClientFunc: func(data []byte, logger echo.Logger) (k8s.ControllerRuntimeInterface, error) {
					return runtimeClient, nil
				},
			}
			return handler, ctrl
		},
		"error getting config by name": func(*testing.T) (*ApplicationHandler, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			runtimeClient := k8sMocks.NewMockControllerRuntimeInterface(ctrl)

			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)
			configFileStore.EXPECT().GetAllByName(gomock.Any()).Return([]model.ConfigFile{{ConfigFileData: []byte("config")}}, errors.New("error"))

			handler := &ApplicationHandler{
				configFileStore: configFileStore,
				runtimeClientFunc: func(data []byte, logger echo.Logger) (k8s.ControllerRuntimeInterface, error) {
					return runtimeClient, nil
				},
			}
			return handler, ctrl
		},
		"error creating tls secret": func(*testing.T) (*ApplicationHandler, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			runtimeClient := k8sMocks.NewMockControllerRuntimeInterface(ctrl)
			runtimeClient.EXPECT().CreateTLSSecretFromName(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))

			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)
			configFileStore.EXPECT().GetAllByName(gomock.Any()).Times(2).Return([]model.ConfigFile{{ConfigFileData: []byte("config")}}, nil)

			handler := &ApplicationHandler{
				configFileStore: configFileStore,
				runtimeClientFunc: func(data []byte, logger echo.Logger) (k8s.ControllerRuntimeInterface, error) {
					return runtimeClient, nil
				},
			}
			return handler, ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			handler, ctrl := tc(t)
			handler.createReverseProxySecrets("", "", []byte{}, []byte{}, echo.New().NewContext(nil, nil))
			ctrl.Finish()
		})
	}
}

func Test_CaptureDiff(t *testing.T) {

	tests := map[string]func(t *testing.T) (*ApplicationHandler, *gomock.Controller){
		"success": func(*testing.T) (*ApplicationHandler, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			applicationStore := mocks.NewMockApplicationStoreInterface(ctrl)
			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)

			taskStore := mocks.NewMockTaskStoreInterface(ctrl)
			taskStore.EXPECT().Update(gomock.Any()).Return(nil)

			applicationStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&model.Application{
				ClusterID: 1,
				StorageArrays: []model.StorageArray{
					{
						StorageArrayType: model.StorageArrayType{
							Name: model.ArrayTypePowerFlex,
						},
					},
				},
			}, nil)

			clusterStore.EXPECT().GetByID(gomock.Any()).Return(&model.Cluster{ConfigFileData: []byte("config")}, nil)

			applicationStateChangeStore := mocks.NewMockApplicationStateChangeStoreInterface(ctrl)
			applicationStateChangeStore.EXPECT().Create(gomock.Any()).Times(2).Return(nil)

			yttClient := yttMocks.NewMockInterface(ctrl)
			yttClient.EXPECT().SetOptions(gomock.Any())
			yttOutput := ytt.Output{}
			yttOutput.Output = &template.Output{}
			yttOutput.Files = []files.OutputFile{}
			yttClient.EXPECT().NamespaceTemplateFromApplication(gomock.Any(), gomock.Any()).Return(yttOutput, nil)
			yttClient.EXPECT().GenerateDynamicSecret(gomock.Any(), gomock.Any(), gomock.Any()).Return(yttOutput, nil)
			yttClient.EXPECT().ConfigMapTemplateFromApplication(gomock.Any(), gomock.Any()).Return(yttOutput, nil)
			yttClient.EXPECT().TemplateFromApplication(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(yttOutput, nil)

			runtimeClient := k8sMocks.NewMockControllerRuntimeInterface(ctrl)
			runtimeClient.EXPECT().CreateNameSpace(gomock.Any(), gomock.Any()).Return(nil)
			runtimeClient.EXPECT().CreateConfigMap(gomock.Any(), gomock.Any()).Return(nil)

			kappClient := kappMocks.NewMockInterface(ctrl)
			kappClient.EXPECT().GetDeployDiff(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("kapp output", nil)

			handler := &ApplicationHandler{
				applicationStore:            applicationStore,
				applicationStateChangeStore: applicationStateChangeStore,
				clusterStore:                clusterStore,
				taskStore:                   taskStore,
				yttClient:                   yttClient,
				kappClient:                  kappClient,
				runtimeClientFunc: func(data []byte, logger echo.Logger) (k8s.ControllerRuntimeInterface, error) {
					return runtimeClient, nil
				},
			}
			return handler, ctrl
		},
		"success with unity": func(*testing.T) (*ApplicationHandler, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			applicationStore := mocks.NewMockApplicationStoreInterface(ctrl)
			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)

			taskStore := mocks.NewMockTaskStoreInterface(ctrl)
			taskStore.EXPECT().Update(gomock.Any()).Return(nil)

			applicationStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&model.Application{
				ClusterID: 1,
				StorageArrays: []model.StorageArray{
					{
						StorageArrayType: model.StorageArrayType{
							Name: model.ArrayTypeUnity,
						},
					},
				},
			}, nil)

			clusterStore.EXPECT().GetByID(gomock.Any()).Return(&model.Cluster{ConfigFileData: []byte("config")}, nil)

			applicationStateChangeStore := mocks.NewMockApplicationStateChangeStoreInterface(ctrl)
			applicationStateChangeStore.EXPECT().Create(gomock.Any()).Times(2).Return(nil)

			yttClient := yttMocks.NewMockInterface(ctrl)
			yttClient.EXPECT().SetOptions(gomock.Any())
			yttOutput := ytt.Output{}
			yttOutput.Output = &template.Output{}
			yttOutput.Files = []files.OutputFile{}
			yttClient.EXPECT().GetEmptySecret(gomock.Any(), gomock.Any()).Return(yttOutput, nil)

			yttClient.EXPECT().NamespaceTemplateFromApplication(gomock.Any(), gomock.Any()).Return(yttOutput, nil)
			yttClient.EXPECT().GenerateDynamicSecret(gomock.Any(), gomock.Any(), gomock.Any()).Return(yttOutput, nil)
			yttClient.EXPECT().ConfigMapTemplateFromApplication(gomock.Any(), gomock.Any()).Return(yttOutput, nil)
			yttClient.EXPECT().TemplateFromApplication(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(yttOutput, nil)

			runtimeClient := k8sMocks.NewMockControllerRuntimeInterface(ctrl)
			runtimeClient.EXPECT().CreateNameSpace(gomock.Any(), gomock.Any()).Return(nil)
			runtimeClient.EXPECT().CreateConfigMap(gomock.Any(), gomock.Any()).Return(nil)
			runtimeClient.EXPECT().CreateSecret(gomock.Any(), gomock.Any()).Return(nil)

			kappClient := kappMocks.NewMockInterface(ctrl)
			kappClient.EXPECT().GetDeployDiff(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("kapp output", nil)

			handler := &ApplicationHandler{
				applicationStore:            applicationStore,
				applicationStateChangeStore: applicationStateChangeStore,
				clusterStore:                clusterStore,
				taskStore:                   taskStore,
				yttClient:                   yttClient,
				kappClient:                  kappClient,
				runtimeClientFunc: func(data []byte, logger echo.Logger) (k8s.ControllerRuntimeInterface, error) {
					return runtimeClient, nil
				},
			}

			return handler, ctrl
		},
		"success with replication module": func(*testing.T) (*ApplicationHandler, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			applicationStore := mocks.NewMockApplicationStoreInterface(ctrl)
			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)

			taskStore := mocks.NewMockTaskStoreInterface(ctrl)
			taskStore.EXPECT().Update(gomock.Any()).Return(nil)

			applicationStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&model.Application{
				ClusterID: 1,
				StorageArrays: []model.StorageArray{
					{
						StorageArrayType: model.StorageArrayType{
							Name: model.ArrayTypePowerFlex,
						},
					},
				},
				ModuleTypes: []model.ModuleType{
					{
						Name: model.ModuleTypeReplication,
					},
				},
				ModuleConfiguration: "target_cluster=1",
			}, nil)

			clusterStore.EXPECT().GetByID(gomock.Any()).AnyTimes().Return(&model.Cluster{ConfigFileData: []byte("config")}, nil)

			moduleTypeStore := mocks.NewMockModuleTypeStoreInterface(ctrl)
			moduleTypeStore.EXPECT().GetByID(gomock.Any()).Return(&model.ModuleType{Name: model.ModuleTypeReplication}, nil)

			applicationStateChangeStore := mocks.NewMockApplicationStateChangeStoreInterface(ctrl)
			applicationStateChangeStore.EXPECT().Create(gomock.Any()).Times(2).Return(nil)

			yttClient := yttMocks.NewMockInterface(ctrl)
			yttClient.EXPECT().SetOptions(gomock.Any())
			yttOutput := ytt.Output{}
			yttOutput.Output = &template.Output{}
			yttOutput.Files = []files.OutputFile{}
			yttClient.EXPECT().NamespaceTemplateFromApplication(gomock.Any(), gomock.Any()).Return(yttOutput, nil)
			yttClient.EXPECT().GenerateDynamicSecret(gomock.Any(), gomock.Any(), gomock.Any()).Return(yttOutput, nil)
			yttClient.EXPECT().ConfigMapTemplateFromApplication(gomock.Any(), gomock.Any()).Return(yttOutput, nil)
			yttClient.EXPECT().TemplateFromApplication(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(yttOutput, nil)

			runtimeClient := k8sMocks.NewMockControllerRuntimeInterface(ctrl)
			runtimeClient.EXPECT().CreateNameSpaceFromName(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			runtimeClient.EXPECT().CreateNameSpace(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			runtimeClient.EXPECT().CreateConfigMap(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			runtimeClient.EXPECT().CreateSecretFromName(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

			kappClient := kappMocks.NewMockInterface(ctrl)
			kappClient.EXPECT().GetDeployDiff(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("kapp output", nil)

			handler := &ApplicationHandler{
				applicationStore:            applicationStore,
				applicationStateChangeStore: applicationStateChangeStore,
				clusterStore:                clusterStore,
				taskStore:                   taskStore,
				yttClient:                   yttClient,
				kappClient:                  kappClient,
				ModuleTypeStore:             moduleTypeStore,
				runtimeClientFunc: func(data []byte, logger echo.Logger) (k8s.ControllerRuntimeInterface, error) {
					return runtimeClient, nil
				},
			}
			return handler, ctrl
		},
		"success with reverse proxy module": func(*testing.T) (*ApplicationHandler, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			applicationStore := mocks.NewMockApplicationStoreInterface(ctrl)
			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)

			taskStore := mocks.NewMockTaskStoreInterface(ctrl)
			taskStore.EXPECT().Update(gomock.Any()).Return(nil)

			applicationStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&model.Application{
				ClusterID: 1,
				StorageArrays: []model.StorageArray{
					{
						StorageArrayType: model.StorageArrayType{
							Name: model.ArrayTypePowerFlex,
						},
					},
				},
				ModuleTypes: []model.ModuleType{
					{
						Name: model.ModuleTypeReplication,
					},
				},
				ModuleConfiguration: "reverseProxy.tlsSecretKeyFile=keyfile reverseProxy.tlsSecretCertFile=certfile",
			}, nil)

			clusterStore.EXPECT().GetByID(gomock.Any()).AnyTimes().Return(&model.Cluster{ConfigFileData: []byte("config")}, nil)

			moduleTypeStore := mocks.NewMockModuleTypeStoreInterface(ctrl)
			moduleTypeStore.EXPECT().GetByID(gomock.Any()).Return(&model.ModuleType{Name: model.ModuleTypeReverseProxy}, nil)

			applicationStateChangeStore := mocks.NewMockApplicationStateChangeStoreInterface(ctrl)
			applicationStateChangeStore.EXPECT().Create(gomock.Any()).Times(2).Return(nil)

			yttClient := yttMocks.NewMockInterface(ctrl)
			yttClient.EXPECT().SetOptions(gomock.Any())
			yttOutput := ytt.Output{}
			yttOutput.Output = &template.Output{}
			yttOutput.Files = []files.OutputFile{}
			yttClient.EXPECT().NamespaceTemplateFromApplication(gomock.Any(), gomock.Any()).Return(yttOutput, nil)
			yttClient.EXPECT().GenerateDynamicSecret(gomock.Any(), gomock.Any(), gomock.Any()).Return(yttOutput, nil)
			yttClient.EXPECT().ConfigMapTemplateFromApplication(gomock.Any(), gomock.Any()).Return(yttOutput, nil)
			yttClient.EXPECT().TemplateFromApplication(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(yttOutput, nil)

			runtimeClient := k8sMocks.NewMockControllerRuntimeInterface(ctrl)
			runtimeClient.EXPECT().CreateNameSpaceFromName(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			runtimeClient.EXPECT().CreateNameSpace(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			runtimeClient.EXPECT().CreateConfigMap(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			runtimeClient.EXPECT().CreateSecretFromName(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
			runtimeClient.EXPECT().CreateTLSSecretFromName(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			kappClient := kappMocks.NewMockInterface(ctrl)
			kappClient.EXPECT().GetDeployDiff(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("kapp output", nil)

			configFileStore := mocks.NewMockConfigFileStoreInterface(ctrl)
			configFileStore.EXPECT().GetAllByName(gomock.Any()).Times(2).Return([]model.ConfigFile{{ConfigFileData: []byte("config")}}, nil)

			CreateReplicationController = func(sourceConfig string, targetConfig string, c echo.Context) error {
				return nil
			}
			handler := &ApplicationHandler{
				applicationStore:            applicationStore,
				applicationStateChangeStore: applicationStateChangeStore,
				clusterStore:                clusterStore,
				taskStore:                   taskStore,
				yttClient:                   yttClient,
				kappClient:                  kappClient,
				ModuleTypeStore:             moduleTypeStore,
				configFileStore:             configFileStore,
				runtimeClientFunc: func(data []byte, logger echo.Logger) (k8s.ControllerRuntimeInterface, error) {
					return runtimeClient, nil
				},
			}
			return handler, ctrl
		},
		"error getting deploy diff": func(*testing.T) (*ApplicationHandler, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			applicationStore := mocks.NewMockApplicationStoreInterface(ctrl)
			clusterStore := mocks.NewMockClusterStoreInterface(ctrl)

			taskStore := mocks.NewMockTaskStoreInterface(ctrl)
			taskStore.EXPECT().Update(gomock.Any()).Return(nil)

			applicationStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&model.Application{
				ClusterID: 1,
				StorageArrays: []model.StorageArray{
					{
						StorageArrayType: model.StorageArrayType{
							Name: model.ArrayTypePowerFlex,
						},
					},
				},
			}, nil)

			clusterStore.EXPECT().GetByID(gomock.Any()).Return(&model.Cluster{ConfigFileData: []byte("config")}, nil)

			applicationStateChangeStore := mocks.NewMockApplicationStateChangeStoreInterface(ctrl)
			applicationStateChangeStore.EXPECT().Create(gomock.Any()).Times(2).Return(nil)

			yttClient := yttMocks.NewMockInterface(ctrl)
			yttClient.EXPECT().SetOptions(gomock.Any())
			yttOutput := ytt.Output{}
			yttOutput.Output = &template.Output{}
			yttOutput.Files = []files.OutputFile{}
			yttClient.EXPECT().NamespaceTemplateFromApplication(gomock.Any(), gomock.Any()).Return(yttOutput, nil)
			yttClient.EXPECT().GenerateDynamicSecret(gomock.Any(), gomock.Any(), gomock.Any()).Return(yttOutput, nil)
			yttClient.EXPECT().ConfigMapTemplateFromApplication(gomock.Any(), gomock.Any()).Return(yttOutput, nil)
			yttClient.EXPECT().TemplateFromApplication(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(yttOutput, nil)

			runtimeClient := k8sMocks.NewMockControllerRuntimeInterface(ctrl)
			runtimeClient.EXPECT().CreateNameSpace(gomock.Any(), gomock.Any()).Return(nil)
			runtimeClient.EXPECT().CreateConfigMap(gomock.Any(), gomock.Any()).Return(nil)

			kappClient := kappMocks.NewMockInterface(ctrl)
			kappClient.EXPECT().GetDeployDiff(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", errors.New("error"))

			handler := &ApplicationHandler{
				applicationStore:            applicationStore,
				applicationStateChangeStore: applicationStateChangeStore,
				clusterStore:                clusterStore,
				taskStore:                   taskStore,
				yttClient:                   yttClient,
				kappClient:                  kappClient,
				runtimeClientFunc: func(data []byte, logger echo.Logger) (k8s.ControllerRuntimeInterface, error) {
					return runtimeClient, nil
				},
			}

			return handler, ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			handler, ctrl := tc(t)
			oldFunction := CreateReplicationController
			defer func() {
				CreateReplicationController = oldFunction
			}()
			CreateReplicationController = func(sourceConfig string, targetConfig string, c echo.Context) error {
				return nil
			}
			handler.captureApplicationDiff(context.Background(), 1, model.Task{}, echo.New().NewContext(nil, nil))
			ctrl.Finish()
		})
	}
}
