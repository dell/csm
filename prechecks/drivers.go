package prechecks

import (
	"github.com/dell/csm-deployment/k8s"
	"github.com/dell/csm-deployment/kubectl"
	"github.com/dell/csm-deployment/model"
)

// Validator is the interface that all validation prechecks must implement
type Validator interface {
	Validate() error
}

const (
	// K8sMinimumSupportedVersion is the minimum supported version for k8s
	K8sMinimumSupportedVersion = "1.19"
	// K8sMaximumSupportedVersion is the maximum supported version for k8s
	K8sMaximumSupportedVersion = "1.21"
	// OpenshiftMinimumSupportedVersion is the minimum supported version for openshift
	OpenshiftMinimumSupportedVersion = "4.6"
	// OpenshiftMaximumSupportedVersion is the maximum supported version for openshift
	OpenshiftMaximumSupportedVersion = "4.7"
)

// GetDriverPrechecks will return a list of prechecks for the specific driver and version
func GetDriverPrechecks(driverType string, clusterData []byte, clusterNodeDetails string) []Validator {
	validators := make([]Validator, 0)

	// common prechecks for all drivers
	validators = append(validators, K8sVersionValidator{
		MinimumVersion: K8sMinimumSupportedVersion,
		MaximumVersion: K8sMaximumSupportedVersion,
		ClusterData:    clusterData,
		K8sClient:      k8s.K8sClient{},
	})
	validators = append(validators, OpenshiftVersionValidator{
		MinimumVersion: OpenshiftMinimumSupportedVersion,
		MaximumVersion: OpenshiftMaximumSupportedVersion,
		ClusterData:    clusterData,
		K8sClient:      k8s.K8sClient{},
	})
	validators = append(validators, VolumeSnapshotResourcesValidator{
		ClusterData:   clusterData,
		KubectlClient: kubectl.Kubectl{},
	})

	switch driverType {
	case model.ArrayTypePowerFlex:
		validators = append(validators, SDCValidator{
			NodeInfo: clusterNodeDetails,
		})
	case model.ArrayTypeIsilon:
		// no specific prechecks
	case model.ArrayTypePowerStore:
		validators = append(validators, ISCSIValidator{
			NodeInfo: clusterNodeDetails,
		})
	case model.ArrayTypeUnity:
		validators = append(validators, ISCSIValidator{
			NodeInfo: clusterNodeDetails,
		})
	case model.ArrayTypePowerMax:
		validators = append(validators, ISCSIValidator{
			NodeInfo: clusterNodeDetails,
		})
	}
	return validators
}
