// Package api for API services
// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
package api

const (
	// UserLoginURI - User login URI
	UserLoginURI = "/api/v1/users/login"

	// ChangePasswordURI - Change Password URI
	ChangePasswordURI = "/api/v1/users/change-password?password=%s"

	// AddCLusterURI - Add CLuster URI
	AddCLusterURI = "/api/v1/clusters"

	// GetAllClustersURI - Get All Clusters URI
	GetAllClustersURI = "/api/v1/clusters"

	// GetClusterByNameURI - Get Cluster by name URI
	GetClusterByNameURI = GetAllClustersURI + "?cluster_name=%s"

	// PatchClusterURI - Patch Cluster URI
	PatchClusterURI = "/api/v1/clusters/%s"

	// DeleteClusterURI - Delete Cluster URI
	DeleteClusterURI = "/api/v1/clusters/%s"

	// AddStorageURI - Add Storage array URI
	AddStorageURI = "/api/v1/storage-arrays"

	// PatchStorageURI - Patch Storage URI
	PatchStorageURI = "/api/v1/storage-arrays/%s"

	// GetAllStorageURI - Get All Storage arrays URI
	GetAllStorageURI = "/api/v1/storage-arrays"

	// GetStorageByParamURI - Get Storage array by custom parameter URI
	GetStorageByParamURI = GetAllStorageURI + "?%s=%s"

	// GetStorageTypeURI - Get Storage array types URI
	GetStorageTypeURI = "/api/v1/storage-array-types"

	// DeleteStorageURI - Delete Storage array URI
	DeleteStorageURI = "/api/v1/storage-arrays/%s"

	// CreateApplicationURI - Create Application URI
	CreateApplicationURI = "/api/v1/applications"

	// UpdateApplicationURI - Update Application URI
	UpdateApplicationURI = "/api/v1/applications/%s"

	// GetAllApplicationURI - Get All Applications URI
	GetAllApplicationURI = "/api/v1/applications"

	// GetApplicationByNameURI - Get Application By Name URI
	GetApplicationByNameURI = GetAllApplicationURI + "?name=%s"

	// DeleteApplicationURI - Delete Application URI
	DeleteApplicationURI = "/api/v1/applications/%s"

	// GetAllTasksURI - Get All Tasks URI
	GetAllTasksURI = "/api/v1/tasks"

	// GetTaskByApplicationNameURI - Get Tasks filtered with application name URI
	GetTaskByApplicationNameURI = GetAllTasksURI + "?application_name=%s"

	// ApproveTaskURI - Approve the Task URI
	ApproveTaskURI = "/api/v1/tasks/%s/approve?updating=%t"

	// RejectTaskURI - Reject the Task URI
	RejectTaskURI = "/api/v1/tasks/%s/cancel?updating=%t"

	// GetModuleTypeURI - Get module types URI
	GetModuleTypeURI = "/api/v1/module-types"

	// GetDriverTypeURI - Get driver types URI
	GetDriverTypeURI = "/api/v1/driver-types"

	// ConfigurationURI - Add configuration file URI
	ConfigurationURI = "/api/v1/configuration-files"

	// GetConfigurationByNameURI - Get configuration file by name URI
	GetConfigurationByNameURI = ConfigurationURI + "?config_name=%s"

	// PatchConfigurationURI - Patch configuration file URI
	PatchConfigurationURI = ConfigurationURI + "/%s"

	// DeleteConfigurationURI - Delete configuration file URI
	DeleteConfigurationURI = ConfigurationURI + "/%s"
)
