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
	"io"
	"net/http"
	"strconv"

	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
)

// CreateConfigFile uploads a configuration file
// @Summary Create a new configuration file
// @Description Create a new configuration file
// @ID create-config-file
// @Tags configuration-file
// @Accept json
// @Produce json
// @Param name formData string true "Name of the configuration file"
// @Param file formData file true "Configuration file"
// @Success 201 {object} configFileResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /configuration-files [post]
func (h *ConfigFileHandler) CreateConfigFile(c echo.Context) error {
	cf := &model.ConfigFile{}

	// Read form fields
	name := c.FormValue("name")
	if len(name) == 0 {
		err := errors.New("name is required")
		if err != nil {
			return c.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, utils.CriticalSeverity, "", err))
		}
	}

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, utils.CriticalSeverity, "", err))
	}
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, utils.CriticalSeverity, "", err))
	}

	defer src.Close()
	data, err := io.ReadAll(src)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}

	cf.Name = name
	cf.ConfigFileData = data

	if err := h.configFileStore.Create(cf); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}

	return c.JSON(http.StatusCreated, newConfigFileResponse(cf))

}

// UpdateConfigFile updates a configuration file
// @Summary Update a configuration file
// @Description Update a configuration file
// @ID update-config-file
// @Tags configuration-file
// @Accept  json
// @Produce  json
// @Param id path string true "Configuration file ID"
// @Param name formData string true "Name of the configuration file"
// @Param file formData file true "Configuration file"
// @Success 204 "No Content"
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /configuration-files/{id} [patch]
func (h *ConfigFileHandler) UpdateConfigFile(c echo.Context) error {
	id := c.Param("id")
	cfID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.ErrorSeverity, "", err))
	}
	cf, err := h.configFileStore.GetByID(uint(cfID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}
	if cf == nil {
		return c.JSON(http.StatusNotFound, utils.NewErrorResponse(http.StatusNotFound, utils.ErrorSeverity, "", err))
	}

	// cluster name is optional during update, if not provided use existing cluster name
	name := c.FormValue("name")
	if len(name) == 0 {
		name = cf.Name
	}

	// configfile is optional during update, if not provided, use existing config file
	file, err := c.FormFile("file")
	var data []byte
	if err != nil {
		data = cf.ConfigFileData
	} else {
		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, utils.CriticalSeverity, "", err))
		}
		defer src.Close()
		data, err = io.ReadAll(src)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
		}
	}

	cf.Name = name
	cf.ConfigFileData = data

	if err := h.configFileStore.Update(cf); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}

	return c.JSON(http.StatusNoContent, nil)
}

// GetConfigFile gets a specific configuration file
// @Summary Get a configuration file
// @Description Get a configuration file
// @ID get-config-file
// @Tags configuration-file
// @Accept  json
// @Produce  json
// @Param id path string true "Configuration file ID"
// @Success 200 {object} configFileResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /configuration-files/{id} [get]
func (h *ConfigFileHandler) GetConfigFile(c echo.Context) error {
	id := c.Param("id")
	configFileID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.ErrorSeverity, "", err))
	}
	cf, err := h.configFileStore.GetByID(uint(configFileID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}
	if cf == nil {
		return c.JSON(http.StatusNotFound, utils.NewErrorResponse(http.StatusNotFound, utils.ErrorSeverity, "", err))
	}
	return c.JSON(http.StatusOK, newConfigFileResponse(cf))
}

// ListConfigFiles returns all the configuration files
// @Summary List all configuration files
// @Description List all configuration files
// @ID list-config-file
// @Tags configuration-file
// @Accept  json
// @Produce  json
// @Param config_name query string false "Name of the configuration file"
// @Success 200 {array} configFileResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /configuration-files [get]
func (h *ConfigFileHandler) ListConfigFiles(c echo.Context) error {

	name := c.QueryParam("config_name")
	var cfs []model.ConfigFile
	var err error
	if name != "" {
		cfs, err = h.configFileStore.GetAllByName(name)
	} else {
		cfs, err = h.configFileStore.GetAll()
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}

	resp := make([]*configFileResponse, 0)
	for _, cf := range cfs {
		resp = append(resp, newConfigFileResponse(&cf))
	}
	return c.JSON(http.StatusOK, resp)
}

// DeleteConfigFile deletes a specific configuration file
// @Summary Delete a configuration file
// @Description Delete a configuration file
// @ID delete-config-file
// @Tags configuration-file
// @Accept  json
// @Produce  json
// @Param id path string true "Configuration file ID"
// @Success 204
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /configuration-files/{id} [delete]
func (h *ConfigFileHandler) DeleteConfigFile(c echo.Context) error {
	id := c.Param("id")
	configFileID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.CriticalSeverity, "", err))
	}
	cf, err := h.configFileStore.GetByID(uint(configFileID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}
	if cf == nil {
		return c.JSON(http.StatusNotFound, utils.NewErrorResponse(http.StatusNotFound, utils.ErrorSeverity, "", err))
	}
	if err := h.configFileStore.Delete(cf); err != nil {
		c.Logger().Errorf("error deleting configuration file: %+v", err)
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}
	return c.JSON(http.StatusNoContent, nil)
}
