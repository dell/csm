// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package handler

import (
	"net/http"
	"strconv"

	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
)

// GetModuleType godoc
// @Summary Get a module type
// @Description Get a module type
// @ID get-module-type
// @Tags module-type
// @Accept  json
// @Produce  json
// @Param id path string true "Module Type ID"
// @Success 200 {object} moduleResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /module-types/{id} [get]
func (h *ModuleTypeHandler) GetModuleType(c echo.Context) error {
	arrayID := c.Param("id")
	id, err := strconv.Atoi(arrayID)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.CriticalSeverity, "", err))
	}
	moduleType, err := h.moduleTypeStore.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}
	if moduleType == nil {
		return c.JSON(http.StatusNotFound, utils.NewErrorResponse(http.StatusNotFound, utils.ErrorSeverity, "", err))
	}
	return c.JSON(http.StatusOK, newModuleResponse(moduleType))
}

// ListModuleType godoc
// @Summary List all module types
// @Description List all module types
// @ID list-module-type
// @Tags module-type
// @Accept  json
// @Produce  json
// @Success 200 {array} moduleResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /module-types [get]
func (h *ModuleTypeHandler) ListModuleType(c echo.Context) error {
	moduleTypes, err := h.moduleTypeStore.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}
	resp := make([]*moduleResponse, 0)
	for _, arr := range moduleTypes {
		resp = append(resp, newModuleResponse(&arr))
	}
	return c.JSON(http.StatusOK, resp)
}
