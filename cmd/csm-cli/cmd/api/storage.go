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

	// StorageTypeIDResponseField - Place holder for field "storage_array_type_id"
	StorageTypeIDResponseField = "storage_array_type_id"

	// EndpointResponseField - Place holder for field "management_endpoint"
	EndpointResponseField = "management_endpoint"
)

// AddStorage - Create new storage array
func AddStorage(endpoint, username, password, uniqueID, storageType string) (*types.StorageResponse, error) {
	addStorageReq := &types.Storage{
		Endpoint:    endpoint,
		Username:    username,
		Password:    password,
		UniqueID:    uniqueID,
		StorageType: storageType,
	}

	addStorageResponse := &types.StorageResponse{}
	err := HTTPClient(http.MethodPost, AddStorageURI, addStorageReq, addStorageResponse)
	if err != nil {
		return nil, err
	}
	return addStorageResponse, nil
}

// GetStorageByParam - return storage array based on parameter and value
func GetStorageByParam(param, value string) ([]types.StorageResponse, error) {
	getStorageResponse := []types.StorageResponse{}
	err := HTTPClient(http.MethodGet, fmt.Sprintf(GetStorageByParamURI, param, value), nil, &getStorageResponse)
	if err != nil {
		return nil, err
	}
	return getStorageResponse, nil
}

// GetAllStorage - returns all storage arrays
func GetAllStorage() ([]types.StorageResponse, error) {
	getStorageResponse := []types.StorageResponse{}
	err := HTTPClient(http.MethodGet, GetStorageByParamURI, nil, &getStorageResponse)
	if err != nil {
		return nil, err
	}
	return getStorageResponse, nil
}

// DeleteStorage - Delete storage array based on ID
func DeleteStorage(uniqueID string) error {
	getStorageResp, err := GetStorageByParam(StorageUniqueIDResponseField, uniqueID)
	if err != nil {
		return errors.New("storage array does not exist")
	}
	if len(getStorageResp) > 1 {
		return errors.New("multiple storage array with same unique id exist")
	}

	err = HTTPClient(http.MethodDelete, fmt.Sprintf(DeleteClusterURI, getStorageResp[0].ID), nil, nil)
	if err != nil {
		return err
	}
	return nil
}

// GetStorageTypeID - returns storageType based on ID
func GetStorageTypeID(storageType string) string {
	mapID := make(map[string]string)
	mapID["unity"] = "0"
	mapID["powermax"] = "1"
	mapID["vxflexos"] = "2"
	mapID["powerstore"] = "3"
	mapID["powerscale"] = "4"
	return mapID[storageType]
}
