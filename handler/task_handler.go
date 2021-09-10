// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package handler

import (
	"github.com/dell/csm-deployment/kapp"
	"github.com/dell/csm-deployment/store"
	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// TaskHandler - Place holder for taska nd application Interfaces
type TaskHandler struct {
	taskStore                   store.TaskStoreInterface
	applicationStore            store.ApplicationStoreInterface
	applicationStateChangeStore store.ApplicationStateChangeStoreInterface
	clusterStore                store.ClusterStoreInterface
	kappClient                  kapp.Interface
}

//NewTaskHandler - returns a new TaskHandler
func NewTaskHandler(ts store.TaskStoreInterface, as store.ApplicationStoreInterface, asc store.ApplicationStateChangeStoreInterface, cs store.ClusterStoreInterface, kapp kapp.Interface) *TaskHandler {
	return &TaskHandler{
		taskStore:                   ts,
		applicationStore:            as,
		applicationStateChangeStore: asc,
		clusterStore:                cs,
		kappClient:                  kapp,
	}
}

// Register -registers a new TaskHandler
func (h *TaskHandler) Register(api *echo.Group) {
	jwtMiddleware := middleware.JWT(utils.JWTSecret)

	tasks := api.Group("/tasks", jwtMiddleware)
	tasks.GET("/:id", h.GetTask)
	tasks.GET("", h.ListTasks)
	tasks.POST("/:id/approve", h.ApproveStateChange)
	tasks.POST("/:id/cancel", h.CancelStateChange)
}
