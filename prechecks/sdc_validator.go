package prechecks

import (
	"encoding/json"
	"fmt"

	"github.com/dell/csm-deployment/k8s"
)

// SDCValidator will validate that the SDC is installed on the cluster nodes
type SDCValidator struct {
	NodeInfo string
}

// Validate will validate that all hosts in the cluster have the SDC installed
func (k SDCValidator) Validate() error {

	var nodes []k8s.Node
	err := json.Unmarshal([]byte(k.NodeInfo), &nodes)
	if err != nil {
		return err
	}
	for _, node := range nodes {
		if _, ok := node.InstalledSoftware["sdc"]; !ok {
			return fmt.Errorf("sdc is not installed on host %s", node.HostName)
		}
	}
	return nil
}
