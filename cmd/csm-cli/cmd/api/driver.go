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
	"fmt"
	"net/http"
	"strings"

	"github.com/dell/csm-deployment/cmd/csm-cli/cmd/api/types"
)

// GetDriverTypeID - returns ID of input driver type
func GetDriverTypeID(driverType string) (string, error) {
	getDriverResponse, err := GetDriverTypes()
	if err != nil {
		return "", fmt.Errorf("failed to get supported driver with error: %v", err)
	}
	// driverType is a combination of name and version
	// split driverType to get corresponding id
	driverTypeSplit := strings.Split(driverType, ":")
	name := ""
	if len(driverTypeSplit) > 0 {
		name = driverTypeSplit[0]
	}
	version := ""
	if len(driverTypeSplit) == 2 {
		version = strings.Trim(driverTypeSplit[1], "v")
	}
	if len(driverTypeSplit) <= 2 {
		for _, driver := range getDriverResponse {
			if driver.StorageType == name && driver.Version == version {
				return driver.ID, nil
			}
		}
	}
	return "", fmt.Errorf("invalid driver type: %s", driverType)
}

// GetDriverTypes - returns supported driver types by CSM API
func GetDriverTypes() ([]types.DriverTypeResponse, error) {
	getDriverTypeResponse := []types.DriverTypeResponse{}
	err := HTTPClient(http.MethodGet, GetDriverTypeURI, nil, &getDriverTypeResponse)
	if err != nil {
		return nil, err
	}
	for i, driver := range getDriverTypeResponse {
		storageTypesResponse, err := GetStorageTypes()
		if err != nil {
			return nil, err
		}
		for _, storage := range storageTypesResponse {
			if storage.ID == driver.StorageTypeID {
				getDriverTypeResponse[i].StorageType = storage.Name
			}
		}
	}
	return getDriverTypeResponse, nil
}
