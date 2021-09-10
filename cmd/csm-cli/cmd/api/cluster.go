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

// AddCluster - Send call to API for add cluster
func AddCluster(clusterName, configFilePath string) (*types.ClusterResponse, error) {
	getClusterResp, err := GetClusterByName(clusterName)
	if err != nil {
		return nil, err
	}
	if len(getClusterResp) > 1 {
		return nil, errors.New("multiple clusters with same name exist")
	}
	if len(getClusterResp) == 1 {
		fmt.Println("cluster already exists with name: " + clusterName)
		return &getClusterResp[0], nil
	}

	reqFields := make(map[string]string)
	reqFields["name"] = clusterName

	addClusterResponse := &types.ClusterResponse{}
	err = HTTPFormDataClient(http.MethodPost, AddCLusterURI, configFilePath, reqFields, addClusterResponse)
	if err != nil {
		return nil, err
	}
	return addClusterResponse, nil
}

// GetClusterByName - Send call to API for get cluster by name
func GetClusterByName(clusterName string) ([]types.ClusterResponse, error) {
	getClusterResponse := []types.ClusterResponse{}
	err := HTTPClient(http.MethodGet, fmt.Sprintf(GetClusterByNameURI, clusterName), nil, &getClusterResponse)
	if err != nil {
		return nil, err
	}
	return getClusterResponse, nil
}

// GetAllClusters - Send call to API for get all clusters
func GetAllClusters() ([]types.ClusterResponse, error) {
	getClusterResponse := []types.ClusterResponse{}
	err := HTTPClient(http.MethodGet, GetClusterByNameURI, nil, &getClusterResponse)
	if err != nil {
		return nil, err
	}
	return getClusterResponse, nil
}

// PatchCluster - Send call to API for update cluster
func PatchCluster(clusterName, newClusterName, newConfigFilePath string) (*types.ClusterResponse, error) {
	getClusterResp, err := GetClusterByName(clusterName)
	if err != nil {
		return nil, err
	}
	if len(getClusterResp) == 0 {
		return nil, fmt.Errorf("cluster does not exist")
	}
	if len(getClusterResp) > 1 {
		return nil, errors.New("multiple clusters with same name exist")
	}

	reqFields := make(map[string]string)
	if newClusterName != "" {
		reqFields["name"] = newClusterName
	}

	patchClusterResponse := &types.ClusterResponse{}
	err = HTTPFormDataClient(http.MethodPatch, fmt.Sprintf(PatchClusterURI, getClusterResp[0].ClusterID), newConfigFilePath, reqFields, patchClusterResponse)
	if err != nil {
		return nil, err
	}
	return patchClusterResponse, nil
}

// DeleteCluster - Send call to API for delete cluster
func DeleteCluster(clusterName string) error {
	getClusterResp, err := GetClusterByName(clusterName)
	if err != nil {
		return err
	}
	if len(getClusterResp) == 0 {
		fmt.Println("cluster does not exist")
		return nil
	}
	if len(getClusterResp) > 1 {
		return errors.New("multiple clusters with same name exist")
	}

	err = HTTPClient(http.MethodDelete, fmt.Sprintf(DeleteClusterURI, getClusterResp[0].ClusterID), nil, nil)
	if err != nil {
		return err
	}
	return nil
}
