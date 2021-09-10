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

// StorageArrayHandler constains the store interface for StorageArray
type StorageArrayHandler struct {
	arrayStore store.StorageArrayStoreInterface
}

// NewStorageArrayHandler creates an handler for a StorageArray
func NewStorageArrayHandler(as store.StorageArrayStoreInterface) *StorageArrayHandler {
	return &StorageArrayHandler{
		arrayStore: as,
	}
}

// Register rosters all the API endpoints for StorageArray
func (h *StorageArrayHandler) Register(api *echo.Group) {
	jwtMiddleware := middleware.JWT(utils.JWTSecret)

	storageArrays := api.Group("/storage-arrays", jwtMiddleware)
	storageArrays.POST("", h.CreateStorageArray)
	storageArrays.GET("", h.ListStorageArrays)
	storageArrays.GET("/:id", h.GetStorageArray)
	storageArrays.DELETE("/:id", h.DeleteStorageArray)
	storageArrays.PATCH("/:id", h.UpdateStorageArray)
}
