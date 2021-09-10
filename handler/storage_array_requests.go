// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package handler

import (
	"fmt"
	"github.com/dell/csm-deployment/utils/constants"
	"strings"

	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
)

type storageArrayCreateRequest struct {
	StorageArrayType   string   `json:"storage_array_type" validate:"required"`
	UniqueID           string   `json:"unique_id" validate:"required"`
	Username           string   `json:"username" validate:"required"`
	Password           string   `json:"password" validate:"required"`
	ManagementEndpoint string   `json:"management_endpoint" validate:"required"`
	MetaData           []string `json:"meta_data"`
} //@name StorageArrayCreateRequest

type storageArrayUpdateRequest struct {
	StorageArrayType   string   `json:"storage_array_type"`
	UniqueID           string   `json:"unique_id"`
	Username           string   `json:"username"`
	Password           string   `json:"password"`
	ManagementEndpoint string   `json:"management_endpoint"`
	MetaData           []string `json:"meta_data"`
} //@name StorageArrayUpdateRequest

type storageArrayResponse struct {
	ID                 string   `json:"id"`
	StorageArrayTypeID string   `json:"storage_array_type_id"`
	UniqueID           string   `json:"unique_id"`
	Username           string   `json:"username"`
	ManagementEndpoint string   `json:"management_endpoint"`
	MetaData           []string `json:"meta_data"`
} //@name StorageArrayResponse

func newStorageArrayResponse(arr *model.StorageArray) *storageArrayResponse {
	r := new(storageArrayResponse)
	r.ID = fmt.Sprintf("%d", arr.ID)
	r.UniqueID = arr.UniqueID
	r.Username = arr.Username
	r.StorageArrayTypeID = fmt.Sprintf("%d", arr.StorageArrayTypeID)
	r.ManagementEndpoint = arr.ManagementEndpoint
	if len(arr.MetaData) == 0 {
		r.MetaData = []string{}
	} else {
		r.MetaData = strings.Split(arr.MetaData, constants.MetadataDelimeter)
	}
	return r
}

func (r *storageArrayCreateRequest) bind(c echo.Context, array *model.StorageArray) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	array.UniqueID = r.UniqueID
	array.Username = r.Username
	array.ManagementEndpoint = r.ManagementEndpoint
	if r.MetaData != nil {
		array.MetaData = strings.Join(r.MetaData, constants.MetadataDelimeter)
	}
	// TODO: it better to store hash, but we will need password for secret creation
	encrypted, err := utils.EncryptPassword([]byte(r.Password))
	if err != nil {
		return err
	}
	array.Password = encrypted
	return nil
}

func (r *storageArrayUpdateRequest) bind(c echo.Context, array *model.StorageArray) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	array.UniqueID = r.UniqueID
	array.Username = r.Username
	array.ManagementEndpoint = r.ManagementEndpoint
	if r.MetaData != nil {
		array.MetaData = strings.Join(r.MetaData, constants.MetadataDelimeter)
	}

	// TODO: it better to store hash, but we will need password for secret creation
	encrypted, err := utils.EncryptPassword([]byte(r.Password))
	if err != nil {
		return err
	}
	array.Password = encrypted
	return nil
}
