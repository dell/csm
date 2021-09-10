// Package types to hold request and response
// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
package types

// User - Struct for user
type User struct {
	Token string `json:"jwtToken,omitempty"`
}

// Cluster - Struct for cluster
type Cluster struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

// Storage - Struct for storage array
type Storage struct {
	Endpoint    string   `json:"management_endpoint,omitempty"`
	Username    string   `json:"username,omitempty"`
	Password    string   `json:"password,omitempty"`
	UniqueID    string   `json:"unique_id,omitempty"`
	StorageType string   `json:"storage_array_type,omitempty"`
	MetaData    []string `json:"meta_data"`
}

// GetStorage - Struct for get storage
type GetStorage struct {
	UniqueID    string `json:"unique_id,omitempty"`
	StorageType string `json:"storage_type,omitempty"`
}

// Application - Struct for application
type Application struct {
	Name                string   `json:"name" validate:"required"`
	ClusterID           string   `json:"cluster_id"`
	DriverTypeID        string   `json:"driver_type_id"`
	ModuleTypes         []string `json:"module_types"`
	StorageArrays       []string `json:"storage_arrays"`
	DriverConfiguration []string `json:"driver_configuration"`
	ModuleConfiguration []string `json:"module_configuration"`
}
