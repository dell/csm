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

// GetStorageArrayType - Get Storage Array Types
// @Summary Get a storage array type
// @Description Get a storage array type
// @ID get-storage-array-type
// @Tags storage-array-type
// @Accept  json
// @Produce  json
// @Param id path string true "Storage Array Type ID"
// @Success 200 {object} storageArrayTypeResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /storage-array-types/{id} [get]
func (h *StorageArrayTypeHandler) GetStorageArrayType(c echo.Context) error {
	storageArrayTypeID := c.Param("id")
	id, err := strconv.Atoi(storageArrayTypeID)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.CriticalSeverity, "", err))
	}
	storageArrayType, err := h.storageArrayTypeStore.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}
	if storageArrayType == nil {
		return c.JSON(http.StatusNotFound, utils.NewErrorResponse(http.StatusNotFound, utils.ErrorSeverity, "", err))
	}
	return c.JSON(http.StatusOK, newStorageArrayTypeResponse(storageArrayType))
}

// ListStorageArrayTypes - List Storage Array Types
// @Summary List all storage array types
// @Description List all storage array types
// @ID list-storage-array-type
// @Tags storage-array-type
// @Accept  json
// @Produce  json
// @Success 200 {array} storageArrayTypeResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /storage-array-types [get]
func (h *StorageArrayTypeHandler) ListStorageArrayTypes(c echo.Context) error {
	storageArrayTypes, err := h.storageArrayTypeStore.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}
	resp := make([]*storageArrayTypeResponse, 0)
	for _, arr := range storageArrayTypes {
		resp = append(resp, newStorageArrayTypeResponse(&arr))
	}
	return c.JSON(http.StatusOK, resp)
}
