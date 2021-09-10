// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package handler

import (
	"github.com/dell/csm-deployment/store"
	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// StorageArrayTypeHandler is the handler for Storage Array Type APIs
type StorageArrayTypeHandler struct {
	storageArrayTypeStore store.StorageArrayTypeStoreInterface
}

// NewStorageArrayTypeHandler creates a new StorageArrayTypeHandler
func NewStorageArrayTypeHandler(sat store.StorageArrayTypeStoreInterface) *StorageArrayTypeHandler {
	return &StorageArrayTypeHandler{
		storageArrayTypeStore: sat,
	}
}

// Register will register all Storage Array Type APIs
func (h *StorageArrayTypeHandler) Register(api *echo.Group) {
	jwtMiddleware := middleware.JWT(utils.JWTSecret)

	storageArrayType := api.Group("/storage-array-types", jwtMiddleware)
	storageArrayType.GET("/:id", h.GetStorageArrayType)
	storageArrayType.GET("", h.ListStorageArrayTypes)
}
