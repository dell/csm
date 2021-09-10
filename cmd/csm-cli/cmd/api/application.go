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
	"strings"

	"github.com/dell/csm-deployment/cmd/csm-cli/cmd/api/types"
)

const (
	// ReplicationTargetClusterParam - Parameter for target cluster in replication
	ReplicationTargetClusterParam = "target_cluster"
)

// CreateApplication - Send call to API for create application
func CreateApplication(name, clusterName, driverType string, driverConfiguration, storageArrays, moduleTypes, moduleConfiguration []string) error {
	clusterResp, err := GetClusterByName(clusterName)
	if err != nil {
		return err
	}
	if len(clusterResp) == 0 {
		return fmt.Errorf("clusters not found")
	}
	if len(clusterResp) > 1 {
		return fmt.Errorf("multiple clusters found with same name")
	}

	driverTypeID, err := GetDriverTypeID(driverType)
	if err != nil {
		return err
	}

	storageArraysID := make([]string, 0)
	for _, uniqueID := range storageArrays {
		storageResp, err := GetStorageByParam(StorageUniqueIDResponseField, uniqueID)
		if err != nil {
			return err
		} else if len(storageResp) == 0 {
			return fmt.Errorf("storage array with unique id: %s not found", uniqueID)
		} else if len(storageResp) > 1 {
			return fmt.Errorf("multiple storage arrays found with same unique Id: %s", uniqueID)
		} else {
			storageArraysID = append(storageArraysID, storageResp[0].ID)
		}
	}
	moduleTypeIDList := make([]string, 0)
	if len(moduleTypes) > 0 {
		for _, moduleType := range moduleTypes {
			moduleTypeID, err := GetModuleTypeID(moduleType)
			if err != nil {
				return err
			}
			moduleTypeIDList = append(moduleTypeIDList, moduleTypeID)
		}
	}

	for i, moduleParam := range moduleConfiguration {
		param := strings.Split(moduleParam, "=")
		if len(param) != 2 {
			return fmt.Errorf("invalid target cluster for replication")
		}
		paramKey := param[0]
		paramValue := param[1]
		if paramKey == ReplicationTargetClusterParam {
			getTargetCluster, err := GetClusterByName(paramValue)
			if err != nil {
				return fmt.Errorf("failed to get target cluster")
			} else if len(getTargetCluster) == 0 {
				return fmt.Errorf("target cluster not found")
			} else if len(getTargetCluster) > 1 {
				return fmt.Errorf("multiple target cluster with same name found")
			}
			moduleConfiguration[i] = fmt.Sprintf("%s=%s", paramKey, getTargetCluster[0].ClusterID)
		}
	}

	createApplicationReq := &types.Application{
		Name:                name,
		ClusterID:           clusterResp[0].ClusterID,
		DriverTypeID:        driverTypeID,
		StorageArrays:       storageArraysID,
		DriverConfiguration: driverConfiguration,
		ModuleTypes:         moduleTypeIDList,
		ModuleConfiguration: moduleConfiguration,
	}

	return HTTPClient(http.MethodPost, CreateApplicationURI, createApplicationReq, nil)
}

// GetAllApplications - Send call to CSM API for get all applications
func GetAllApplications() ([]types.ApplicationResponse, error) {
	getApplicationResponse := []types.ApplicationResponse{}
	err := HTTPClient(http.MethodGet, GetAllApplicationURI, nil, &getApplicationResponse)
	if err != nil {
		return nil, err
	}
	for i, application := range getApplicationResponse {
		getApplicationResponse[i].Status = GetApplicationTaskStatus(application.Name)
	}
	return getApplicationResponse, nil
}

// GetApplicationByName - call to CSM API for get application by name
func GetApplicationByName(name string) ([]types.ApplicationResponse, error) {
	getApplicationResponse := []types.ApplicationResponse{}
	err := HTTPClient(http.MethodGet, fmt.Sprintf(GetApplicationByNameURI, name), nil, &getApplicationResponse)
	if err != nil {
		return nil, err
	}
	for i, application := range getApplicationResponse {
		getApplicationResponse[i].Status = GetApplicationTaskStatus(application.Name)
	}
	return getApplicationResponse, nil
}

// GetApplicationTaskStatus - call to CSM API for get status for task corresponding to application
func GetApplicationTaskStatus(name string) string {
	taskResp, err := GetTaskByApplicationName(name)
	if err != nil {
		return ""
	}
	if len(taskResp) == 0 {
		return "NoTaskForApplication"
	}
	return taskResp[0].Status
}

// DeleteApplication - call to CSM API for Delete Application
func DeleteApplication(name string) error {
	getApplicationResp, err := GetApplicationByName(name)
	if err != nil {
		return fmt.Errorf("find application failed with error: %v", err)
	}
	if len(getApplicationResp) == 0 {
		fmt.Println("application does not exist")
		return nil
	}
	if len(getApplicationResp) > 1 {
		return errors.New("multiple applications with same unique id exist")
	}

	return HTTPClient(http.MethodDelete, fmt.Sprintf(DeleteApplicationURI, getApplicationResp[0].ID), nil, nil)
}
