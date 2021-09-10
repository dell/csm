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

// ISCSIValidator will validate that iscsi is installed on the cluster nodes
type ISCSIValidator struct {
	NodeInfo string
}

// Validate will validate that all hosts in the cluster have iscsi installed
func (k ISCSIValidator) Validate() error {

	var nodes []k8s.Node
	err := json.Unmarshal([]byte(k.NodeInfo), &nodes)
	if err != nil {
		return err
	}
	for _, node := range nodes {
		if _, ok := node.InstalledSoftware["iscsi"]; !ok {
			return fmt.Errorf("iscsi is not installed on host %s", node.HostName)
		}
	}
	return nil
}
