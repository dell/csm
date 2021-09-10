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

// AddConfiguration - Send call to API for add configuration file
func AddConfiguration(fileName, filePath string) (*types.ConfigurationFileResponse, error) {
	getConfigResp, err := GetConfigurationByName(fileName)
	if err != nil {
		return nil, err
	}
	if len(getConfigResp) > 1 {
		return nil, errors.New("multiple configuration files with same name exist")
	}
	if len(getConfigResp) == 1 {
		fmt.Println("configuration file already exists with name: " + fileName)
		return &getConfigResp[0], nil
	}

	reqFields := make(map[string]string)
	reqFields["name"] = fileName

	addConfigFileResponse := &types.ConfigurationFileResponse{}
	err = HTTPFormDataClient(http.MethodPost, ConfigurationURI, filePath, reqFields, addConfigFileResponse)
	if err != nil {
		return nil, err
	}
	return addConfigFileResponse, nil
}

// GetConfigurationByName - Send call to API for get configuration file by name
func GetConfigurationByName(fileName string) ([]types.ConfigurationFileResponse, error) {
	getConfigResponse := []types.ConfigurationFileResponse{}
	err := HTTPClient(http.MethodGet, fmt.Sprintf(GetConfigurationByNameURI, fileName), nil, &getConfigResponse)
	if err != nil {
		return nil, err
	}
	return getConfigResponse, nil
}

// GetAllConfigurations - Send call to API for get all configuration files
func GetAllConfigurations() ([]types.ConfigurationFileResponse, error) {
	getConfigResponse := []types.ConfigurationFileResponse{}
	err := HTTPClient(http.MethodGet, ConfigurationURI, nil, &getConfigResponse)
	if err != nil {
		return nil, err
	}
	return getConfigResponse, nil
}

// PatchConfiguration - Send call to API for update configuration file
func PatchConfiguration(fileName, newFileName, newFilePath string) (*types.ConfigurationFileResponse, error) {
	getConfigResp, err := GetConfigurationByName(fileName)
	if err != nil {
		return nil, err
	}
	if len(getConfigResp) == 0 {
		return nil, fmt.Errorf("configuration file does not exist")
	}
	if len(getConfigResp) > 1 {
		return nil, errors.New("multiple configuration files with same name exist")
	}

	reqFields := make(map[string]string)
	if newFileName != "" {
		reqFields["name"] = newFileName
	} else {
		reqFields["name"] = fileName
	}

	patchConfigResponse := &types.ConfigurationFileResponse{}
	err = HTTPFormDataClient(http.MethodPatch, fmt.Sprintf(PatchConfigurationURI, getConfigResp[0].ID), newFilePath, reqFields, patchConfigResponse)
	if err != nil {
		return nil, err
	}
	return patchConfigResponse, nil
}

// DeleteConfiguration - Send call to API for delete configuration file
func DeleteConfiguration(fileName string) error {
	getConfigResp, err := GetConfigurationByName(fileName)
	if err != nil {
		return err
	}
	if len(getConfigResp) == 0 {
		return fmt.Errorf("configuration file does not exist")
	}
	if len(getConfigResp) > 1 {
		return errors.New("multiple configuration files with same name exist")
	}

	err = HTTPClient(http.MethodDelete, fmt.Sprintf(DeleteConfigurationURI, getConfigResp[0].ID), nil, nil)
	if err != nil {
		return err
	}
	return nil
}
