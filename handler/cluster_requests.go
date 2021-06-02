package handler

import (
	"github.com/dell/csm-deployment/model"
)

type clusterResponse struct {
	Cluster struct {
		ClusterID   uint   `json:"cluster_id"`
		ClusterName string `json:"cluster_name"`
		Nodes       string `json:"nodes"`
		ConfigFile  string `json:"config_file"`
	} `json:"cluster"`
}

func newClusterResponse(u *model.Cluster) *clusterResponse {
	r := new(clusterResponse)
	r.Cluster.ClusterID = u.ID
	r.Cluster.ClusterName = u.ClusterName
	r.Cluster.ConfigFile = string(u.ConfigFileData)
	r.Cluster.Nodes = u.ClusterDetails.Nodes
	return r
}
