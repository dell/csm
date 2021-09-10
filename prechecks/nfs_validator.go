// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package prechecks

import (
	"encoding/json"
	"fmt"

	"github.com/dell/csm-deployment/k8s"
)

// NFSValidator will validate that nfs is installed on at least one cluster node
type NFSValidator struct {
	NodeInfo string
}

// Validate will validate that at least one host in the cluster have nfs installed
func (k NFSValidator) Validate() error {

	var nodes []k8s.Node
	err := json.Unmarshal([]byte(k.NodeInfo), &nodes)
	if err != nil {
		return err
	}
	for _, node := range nodes {
		if _, ok := node.InstalledSoftware["nfs"]; ok {
			return nil
		}
	}
	return fmt.Errorf("nfs is not installed on any host")
}
