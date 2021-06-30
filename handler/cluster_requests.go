package handler

import (
	"github.com/dell/csm-deployment/model"
)

type clusterListResponse struct {
	Clusters []*clusterResponse `json:"clusters"`
}

type clusterResponse struct {
	ClusterID   uint   `json:"cluster_id"`
	ClusterName string `json:"cluster_name"`
	Nodes       string `json:"nodes"`
}

func newClusterResponse(u *model.Cluster) *clusterResponse {
	r := &clusterResponse{}
	r.ClusterID = u.ID
	r.ClusterName = u.ClusterName
	r.Nodes = u.ClusterDetails.Nodes
	return r
}
