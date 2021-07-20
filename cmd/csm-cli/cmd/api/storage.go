package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dell/csm-deployment/cmd/csm-cli/cmd/api/types"
)

const (
	StorageUniqueIdResponseField = "unique_id"
	StorageTypeIdResponseField   = "storage_array_type_id"
	EndpointResponseField        = "management_endpoint"
)

func AddStorage(endpoint, username, password, uniqueId, storageType string) (*types.StorageResponse, error) {
	addStorageReq := &types.Storage{
		Endpoint:    endpoint,
		Username:    username,
		Password:    password,
		UniqueId:    uniqueId,
		StorageType: storageType,
	}

	addStorageResponse := &types.StorageResponse{}
	err := HttpClient(http.MethodPost, AddStorageURI, addStorageReq, addStorageResponse)
	if err != nil {
		return nil, err
	}
	return addStorageResponse, nil
}

func GetStorageByParam(param, value string) ([]types.StorageResponse, error) {
	getStorageResponse := []types.StorageResponse{}
	err := HttpClient(http.MethodGet, fmt.Sprintf(GetStorageByParamURI, param, value), nil, &getStorageResponse)
	if err != nil {
		return nil, err
	}
	return getStorageResponse, nil
}

func GetAllStorage() ([]types.StorageResponse, error) {
	getStorageResponse := []types.StorageResponse{}
	err := HttpClient(http.MethodGet, GetStorageByParamURI, nil, &getStorageResponse)
	if err != nil {
		return nil, err
	}
	return getStorageResponse, nil
}

func DeleteStorage(uniqueId string) error {
	getStorageResp, err := GetStorageByParam(StorageUniqueIdResponseField, uniqueId)
	if err != nil {
		return errors.New("storage array does not exist")
	}
	if len(getStorageResp) > 1 {
		return errors.New("multiple storage array with same unique id exist")
	}

	err = HttpClient(http.MethodDelete, fmt.Sprintf(DeleteClusterURI, getStorageResp[0].Id), nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func GetStorageTypeId(storageType string) string {
	mapId := make(map[string]string)
	mapId["unity"] = "0"
	mapId["powermax"] = "1"
	mapId["vxflexos"] = "2"
	mapId["powerstore"] = "3"
	mapId["powerscale"] = "4"
	return mapId[storageType]
}
