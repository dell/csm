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

// ConfigFileHandler is the handler for configuration-file APIs
type ConfigFileHandler struct {
	configFileStore store.ConfigFileStoreInterface
}

// NewConfigFileHandler creates a new ConfigFileHandler
func NewConfigFileHandler(cfs store.ConfigFileStoreInterface) *ConfigFileHandler {
	return &ConfigFileHandler{
		configFileStore: cfs,
	}
}

// Register will register all  configuration-file APIs
func (h *ConfigFileHandler) Register(api *echo.Group) {
	jwtMiddleware := middleware.JWT(utils.JWTSecret)

	configFile := api.Group("/configuration-files", jwtMiddleware)
	configFile.GET("/:id", h.GetConfigFile)
	configFile.POST("", h.CreateConfigFile)
	configFile.GET("", h.ListConfigFiles)
	configFile.DELETE("/:id", h.DeleteConfigFile)
	configFile.PATCH("/:id", h.UpdateConfigFile)
}
