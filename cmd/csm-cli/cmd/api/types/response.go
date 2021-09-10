// Package types to hold request and response
// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
package types

// JWTToken - captures JWT Token
var JWTToken string

// ClusterResponse - Struct to capture cluster response
type ClusterResponse struct {
	ClusterID   string `json:"cluster_id"`
	ClusterName string `json:"cluster_name"`
	Nodes       string `json:"nodes"`
}

// StorageResponse - Struct to capture storage array response
type StorageResponse struct {
	ID            string   `json:"id"`
	StorageTypeID string   `json:"storage_array_type_id"`
	UniqueID      string   `json:"unique_id"`
	Username      string   `json:"username"`
	Endpoint      string   `json:"management_endpoint"`
	MetaData      []string `json:"meta_data"`
}

// ApplicationResponse - Struct to capture application response
type ApplicationResponse struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	ClusterID           string   `json:"cluster_id"`
	DriverTypeID        string   `json:"driver_type_id"`
	ModuleTypes         []string `json:"module_types"`
	StorageArrays       []string `json:"storage_arrays"`
	DriverConfiguration []string `json:"driver_configuration"`
	ModuleConfiguration []string `json:"module_configuration"`
	ApplicationOutput   string   `json:"application_output"`
	Status              string
}

// TaskResponse - Struct to capture task response
type TaskResponse struct {
	ID              string                       `json:"id"`
	Status          string                       `json:"status"`
	ApplicationName string                       `json:"application_name"`
	Logs            string                       `json:"logs"`
	Links           map[string]map[string]string `json:"_links"`
}

// StorageTypeResponse - Struct to capture storage type response
type StorageTypeResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ModuleTypeResponse - struct to capture module type response
type ModuleTypeResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Version    string `json:"version"`
	Standalone bool   `json:"standalone"`
}

// DriverTypeResponse - struct to capture driver type response
type DriverTypeResponse struct {
	ID            string `json:"id"`
	StorageTypeID string `json:"storage_array_type_id"`
	Version       string `json:"version"`
	StorageType   string
}

// ConfigurationFileResponse - struct to capture configuration file response
type ConfigurationFileResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ErrorResponse - struct to capture error response
type ErrorResponse struct {
	StatusCode int             `json:"http_status_code"`
	Messages   []ErrorMessages `json:"messages"`
}

// ErrorMessages - struct to capture error messages
type ErrorMessages struct {
	Message     string `json:"message"`
	MessageDesc string `json:"message_l10n"`
}
