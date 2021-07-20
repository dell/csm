package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dell/csm-deployment/cmd/csm-cli/cmd/api/types"
)

//@TODO implement idempotency in all api methods (Decide whether to implement at api level or cli level)
func AddCluster(clusterName, configFilePath string) (*types.ClusterResponse, error) {
	reqFields := make(map[string]string)
	reqFields["name"] = clusterName

	addClusterResponse := &types.ClusterResponse{}
	err := HttpClusterClient(http.MethodPost, AddCLusterURI, configFilePath, reqFields, addClusterResponse)
	if err != nil {
		return nil, err
	}
	return addClusterResponse, nil
}

func GetClusterByName(clusterName string) ([]types.ClusterResponse, error) {
	getClusterResponse := []types.ClusterResponse{}
	err := HttpClient(http.MethodGet, fmt.Sprintf(GetClusterByNameURI, clusterName), nil, &getClusterResponse)
	if err != nil {
		return nil, err
	}
	return getClusterResponse, nil
}

func GetAllClusters() ([]types.ClusterResponse, error) {
	getClusterResponse := []types.ClusterResponse{}
	err := HttpClient(http.MethodGet, GetClusterByNameURI, nil, &getClusterResponse)
	if err != nil {
		return nil, err
	}
	return getClusterResponse, nil
}

func PatchCluster(clusterName, newClusterName, newConfigFilePath string) (*types.ClusterResponse, error) {
	getClusterResp, err := GetClusterByName(clusterName)
	if err != nil {
		return nil, errors.New("cluster does not exist")
	}
	if len(getClusterResp) > 1 {
		return nil, errors.New("multiple clusters with same name exist")
	}

	reqFields := make(map[string]string)
	if newClusterName != "" {
		reqFields["name"] = newClusterName
	}

	patchClusterResponse := &types.ClusterResponse{}
	err = HttpClusterClient(http.MethodPatch, fmt.Sprintf(PatchClusterURI, getClusterResp[0].ClusterId), newConfigFilePath, reqFields, patchClusterResponse)
	if err != nil {
		return nil, err
	}
	return patchClusterResponse, nil
}

func DeleteCluster(clusterName string) error {
	getClusterResp, err := GetClusterByName(clusterName)
	if err != nil {
		return errors.New("cluster does not exist")
	}
	if len(getClusterResp) > 1 {
		return errors.New("multiple clusters with same name exist")
	}

	err = HttpClient(http.MethodDelete, fmt.Sprintf(DeleteClusterURI, getClusterResp[0].ClusterId), nil, nil)
	if err != nil {
		return err
	}
	return nil
}
