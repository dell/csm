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
