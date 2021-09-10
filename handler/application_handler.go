// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package handler

import (
	"github.com/dell/csm-deployment/k8s"
	"github.com/dell/csm-deployment/kapp"
	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/prechecks"
	"github.com/dell/csm-deployment/store"
	"github.com/dell/csm-deployment/utils"
	"github.com/dell/csm-deployment/ytt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// ApplicationHandler - Place holder for Api Interfaces
type ApplicationHandler struct {
	applicationStore            store.ApplicationStoreInterface
	arrayStore                  store.StorageArrayStoreInterface
	taskStore                   store.TaskStoreInterface
	clusterStore                store.ClusterStoreInterface
	applicationStateChangeStore store.ApplicationStateChangeStoreInterface
	ModuleTypeStore             store.ModuleTypeStoreInterface
	driverStore                 store.DriverTypeStoreInterface
	configFileStore             store.ConfigFileStoreInterface
	kappClient                  kapp.Interface
	yttClient                   ytt.Interface
	precheckHandler             PrecheckHandlerInterface
	runtimeClientFunc           func(data []byte, logger echo.Logger) (k8s.ControllerRuntimeInterface, error)
	SkipGoRoutine               bool
}

// PrecheckHandler holds references to perform driver and modules prechecks
type PrecheckHandler struct {
	driverStore     store.DriverTypeStoreInterface
	clusterStore    store.ClusterStoreInterface
	configFileStore store.ConfigFileStoreInterface
	precheckGetter  PrecheckGetterInterface
}

// PrecheckGetterInterface provides an interface to get prechecks for specific resources
//go:generate mockgen -destination=mocks/precheck_getter_interface.go -package=mocks github.com/dell/csm-deployment/handler PrecheckGetterInterface
type PrecheckGetterInterface interface {
	GetDriverPrechecks(driverType string, clusterData []byte, clusterNodeDetails string, modules []model.ModuleType, logger echo.Logger) []prechecks.Validator
	GetModuleTypePrechecks(moduleType, moduleConfig string, clusterData []byte, cfs []model.ConfigFile, availableModules map[string]string) []prechecks.Validator
}

// PrecheckHandlerInterface provides an interface to perform driver and module prechecks
//go:generate mockgen -destination=mocks/precheck_handler_interface.go -package=mocks github.com/dell/csm-deployment/handler PrecheckHandlerInterface
type PrecheckHandlerInterface interface {
	Precheck(c echo.Context, clusterID uint, driverID uint, modules []model.ModuleType, moduleConfig string) error
}

// NewApplicationHandler -  returns a new ApplicationHandler
func NewApplicationHandler(is store.ApplicationStoreInterface,
	ts store.TaskStoreInterface,
	cs store.ClusterStoreInterface,
	asc store.ApplicationStateChangeStoreInterface,
	as store.StorageArrayStoreInterface,
	ms store.ModuleTypeStoreInterface,
	ds store.DriverTypeStoreInterface,
	cf store.ConfigFileStoreInterface,
	kappClient kapp.Interface,
	yttClient ytt.Interface) *ApplicationHandler {
	return &ApplicationHandler{
		applicationStore:            is,
		taskStore:                   ts,
		clusterStore:                cs,
		applicationStateChangeStore: asc,
		arrayStore:                  as,
		ModuleTypeStore:             ms,
		driverStore:                 ds,
		configFileStore:             cf,
		kappClient:                  kappClient,
		precheckHandler:             &PrecheckHandler{driverStore: ds, clusterStore: cs, configFileStore: cf, precheckGetter: prechecks.PrecheckGetter{}},
		yttClient:                   yttClient,
		SkipGoRoutine:               false,
		runtimeClientFunc:           GetRuntimeClient,
	}
}

//Register new Application Handler
func (h *ApplicationHandler) Register(api *echo.Group) {
	jwtMiddleware := middleware.JWT(utils.JWTSecret)

	applications := api.Group("/applications", jwtMiddleware)
	applications.POST("", h.CreateApplication)
	applications.GET("/:id", h.GetApplication)
	applications.GET("", h.ListApplications)
	applications.DELETE("/:id", h.DeleteApplication)
}
