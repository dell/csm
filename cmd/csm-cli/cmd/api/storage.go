// Package api for API services
// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dell/csm-deployment/cmd/csm-cli/cmd/api/types"
)

const (
	// StorageUniqueIDResponseField - Place holder for field "unique_id"
	StorageUniqueIDResponseField = "unique_id"

	// StorageTypeResponseField - Place holder for field "storage_type"
	StorageTypeResponseField = "storage_type"
)

// AddStorage - Create new storage array
func AddStorage(endpoint, username, password, uniqueID, storageType string, metadata []string) (*types.StorageResponse, error) {
	getStorageResp, err := GetStorageByParam(StorageUniqueIDResponseField, uniqueID)
	if err != nil {
		return nil, err
	}
	if len(getStorageResp) > 1 {
		return nil, errors.New("multiple storage array with same unique id exist")
	}
	if len(getStorageResp) == 1 {
		fmt.Println("storage array already exists with unique ID: " + uniqueID)
		return &getStorageResp[0], nil
	}

	if password == "" {
		return nil, fmt.Errorf("empty password")
	}
	addStorageReq := &types.Storage{
		Endpoint:    endpoint,
		Username:    username,
		Password:    password,
		UniqueID:    uniqueID,
		StorageType: storageType,
		MetaData:    metadata,
	}

	addStorageResponse := &types.StorageResponse{}
	err = HTTPClient(http.MethodPost, AddStorageURI, addStorageReq, addStorageResponse)
	if err != nil {
		return nil, err
	}
	return addStorageResponse, nil
}

// PatchStorage - Send call to CSM API for update storage array
func PatchStorage(endpoint, username, password, oldUniqueID, storageType, newUniqueID string, metadata []string) (*types.StorageResponse, error) {
	getStorageResp, err := GetStorageByParam(StorageUniqueIDResponseField, oldUniqueID)
	if err != nil {
		return nil, err
	}
	if len(getStorageResp) == 0 {
		return nil, errors.New("storage array does not exist")
	}
	if len(getStorageResp) > 1 {
		return nil, errors.New("multiple storage array with same unique id exist")
	}

	updateStorageReq := &types.Storage{
		Endpoint:    endpoint,
		Username:    username,
		Password:    password,
		UniqueID:    newUniqueID,
		StorageType: storageType,
		MetaData:    metadata,
	}

	updateStorageResponse := &types.StorageResponse{}
	err = HTTPClient(http.MethodPatch, fmt.Sprintf(PatchStorageURI, getStorageResp[0].ID), updateStorageReq, updateStorageResponse)
	if err != nil {
		return nil, err
	}
	return updateStorageResponse, nil
}

// GetStorageByParam - Send call to CSM API for get storage by input param
func GetStorageByParam(param, value string) ([]types.StorageResponse, error) {
	getStorageResponse := []types.StorageResponse{}
	err := HTTPClient(http.MethodGet, fmt.Sprintf(GetStorageByParamURI, param, value), nil, &getStorageResponse)
	if err != nil {
		return nil, err
	}
	return getStorageResponse, nil
}

// GetAllStorage - Send call to CSM API for get all storage arrays
func GetAllStorage() ([]types.StorageResponse, error) {
	getStorageResponse := []types.StorageResponse{}
	err := HTTPClient(http.MethodGet, GetAllStorageURI, nil, &getStorageResponse)
	if err != nil {
		return nil, err
	}
	return getStorageResponse, nil
}

// DeleteStorage - Send call to CSM API for delete storage array based on unique ID
func DeleteStorage(uniqueID string) error {
	getStorageResp, err := GetStorageByParam(StorageUniqueIDResponseField, uniqueID)
	if err != nil {
		return err
	}
	if len(getStorageResp) == 0 {
		fmt.Println("storage array doesn't exist with name " + uniqueID)
		return nil
	}
	if len(getStorageResp) > 1 {
		return errors.New("multiple storage array with same unique id exist")
	}

	err = HTTPClient(http.MethodDelete, fmt.Sprintf(DeleteStorageURI, getStorageResp[0].ID), nil, nil)
	if err != nil {
		return err
	}
	return nil
}

// GetStorageTypes - returns supported storage types by CSM API
func GetStorageTypes() ([]types.StorageTypeResponse, error) {
	getStorageTypeResponse := []types.StorageTypeResponse{}
	err := HTTPClient(http.MethodGet, GetStorageTypeURI, nil, &getStorageTypeResponse)
	if err != nil {
		return nil, err
	}
	return getStorageTypeResponse, nil
}
