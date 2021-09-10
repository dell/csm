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

// ModuleTypeHandler is the handler for Module Type APIs
type ModuleTypeHandler struct {
	moduleTypeStore store.ModuleTypeStoreInterface
}

// NewModuleTypeHandler creates a new ModuleTypeHandler
func NewModuleTypeHandler(as store.ModuleTypeStoreInterface) *ModuleTypeHandler {
	return &ModuleTypeHandler{
		moduleTypeStore: as,
	}
}

// Register will register all Module Type APIs
func (h *ModuleTypeHandler) Register(api *echo.Group) {
	jwtMiddleware := middleware.JWT(utils.JWTSecret)

	moduleType := api.Group("/module-types", jwtMiddleware)
	moduleType.GET("/:id", h.GetModuleType)
	moduleType.GET("", h.ListModuleType)
}
