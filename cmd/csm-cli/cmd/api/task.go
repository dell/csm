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

	"github.com/dell/csm-deployment/cmd/csm-cli/cmd/api/types"
)

// GetTaskByApplicationName - Send call to CSM API for get task by application name
func GetTaskByApplicationName(applicationName string) ([]types.TaskResponse, error) {
	getTaskResponse := []types.TaskResponse{}
	err := HTTPClient(http.MethodGet, fmt.Sprintf(GetTaskByApplicationNameURI, applicationName), nil, &getTaskResponse)
	if err != nil {
		return nil, err
	}
	return getTaskResponse, nil
}

// GetAllTasks - Send call to CSM API for get all tasks
func GetAllTasks() ([]types.TaskResponse, error) {
	getTaskResponse := []types.TaskResponse{}
	err := HTTPClient(http.MethodGet, GetAllTasksURI, nil, &getTaskResponse)
	if err != nil {
		return nil, err
	}
	return getTaskResponse, nil
}

// ApproveTask - Send call to CSM API for approve task for an application
func ApproveTask(applicationName string, update bool) error {
	getTaskResp, err := GetTaskByApplicationName(applicationName)
	if err != nil {
		return err
	}
	if len(getTaskResp) == 0 {
		return fmt.Errorf("no tasks for application %s", applicationName)
	}

	err = HTTPClient(http.MethodPost, fmt.Sprintf(ApproveTaskURI, getTaskResp[0].ID, update), nil, nil)
	if err != nil {
		return err
	}
	return nil
}

// RejectTask - Send call to CSM API for reject task for an application
func RejectTask(applicationName string, update bool) error {
	getTaskResp, err := GetTaskByApplicationName(applicationName)
	if err != nil {
		return err
	}
	if len(getTaskResp) == 0 {
		return fmt.Errorf("no tasks for application %s", applicationName)
	}

	err = HTTPClient(http.MethodPost, fmt.Sprintf(RejectTaskURI, getTaskResp[0].ID, update), nil, nil)
	if err != nil {
		return err
	}
	return nil
}
