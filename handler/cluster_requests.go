// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package handler

import (
	"fmt"

	"github.com/dell/csm-deployment/model"
)

type clusterResponse struct {
	ClusterID   string `json:"cluster_id"`
	ClusterName string `json:"cluster_name"`
	// The nodes
	Nodes string `json:"nodes"`
} //@name ClusterResponse

func newClusterResponse(u *model.Cluster) *clusterResponse {
	r := &clusterResponse{}
	r.ClusterID = fmt.Sprintf("%d", u.ID)
	r.ClusterName = u.ClusterName
	r.Nodes = u.ClusterDetails.Nodes
	return r
}
