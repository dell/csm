package handler

import (
	"github.com/dell/csm-deployment/model"
)

// ClusterListResponse responds body for list of cluster instances
type ClusterListResponse struct {
	// List of ClusterResponse
	Clusters []*ClusterResponse `json:"clusters"`
} //@name ClusterListResponse

// ClusterResponse responds body for cluster instance
type ClusterResponse struct {
	// Unique identifier of the cluster
	ClusterID uint `json:"cluster_id"`
	// Unique name of the cluster
	// This value must contain 128 or fewer printable Unicode characters.
	ClusterName string `json:"cluster_name"`
	// The nodes
	Nodes string `json:"nodes"`
} //@name ClusterResponse

func newClusterResponse(u *model.Cluster) *ClusterResponse {
	r := &ClusterResponse{}
	r.ClusterID = u.ID
	r.ClusterName = u.ClusterName
	r.Nodes = u.ClusterDetails.Nodes
	return r
}
