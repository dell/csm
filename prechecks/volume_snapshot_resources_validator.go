package prechecks

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	snapshotResources = []string{"VolumeSnapshotClass", "VolumeSnapshotContent", "VolumeSnapshot"}
)

// K8sClientExplainInterface is the required interface for querying the k8s cluster
//go:generate mockgen -destination=mocks/k8s_client_explain_interface.go -package=mocks github.com/dell/csm-deployment/prechecks K8sClientExplainInterface
type K8sClientExplainInterface interface {
	Explain([]byte, string) (*metav1.APIResource, string, error)
}

// VolumeSnapshotResourcesValidator validates the required VolumeSnapshot CRDs and versions on the k8s cluster
type VolumeSnapshotResourcesValidator struct {
	ClusterData []byte
	K8sClient   K8sClientExplainInterface
}

// Validate will check that the expected CRD resources exist and that they are not of the version 'v1alphav1'
func (k VolumeSnapshotResourcesValidator) Validate() error {
	for _, resource := range snapshotResources {
		_, groupVersion, err := k.K8sClient.Explain(k.ClusterData, resource)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("unable to find CRD for %s", resource))
		}
		re := regexp.MustCompile(`v1alpha1`)
		hasAlphav1 := re.Match([]byte(groupVersion))
		if hasAlphav1 {
			return fmt.Errorf("has alphav1 of %s", resource)
		}
	}

	return nil
}
