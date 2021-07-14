package kubectl

import (
	"os/exec"
)

// Kubectl is used to interact with the kubectl binary
type Kubectl struct{}

// Explain will return output from a `kubectl explain` command on the given resource
func (k Kubectl) Explain(configFile string, resource string) ([]byte, error) {
	cmd := exec.Command("/app/kubectl", "--kubeconfig", configFile, "explain", resource)
	return cmd.Output()
}
