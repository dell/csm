package prechecks

import (
	"fmt"
	"log"
	"strconv"
)

// K8sClientVersionInterface is the required interface for querying the k8s cluster
//go:generate mockgen -destination=mocks/k8s_client_version_interface.go -package=mocks github.com/dell/csm-deployment/prechecks K8sClientVersionInterface
type K8sClientVersionInterface interface {
	IsOpenShift([]byte) (bool, error)
	GetVersion([]byte) (string, error)
}

// K8sVersionValidator will validate the k8s version of the cluster
type K8sVersionValidator struct {
	MinimumVersion string
	MaximumVersion string
	ClusterData    []byte
	K8sClient      K8sClientVersionInterface
}

// Validate will validate the version of the k8s cluster is between the min/max supported versions
func (k K8sVersionValidator) Validate() error {
	isOpenshift, err := k.K8sClient.IsOpenShift(k.ClusterData)
	if err != nil {
		return err
	}
	if isOpenshift {
		log.Printf("cluster is openshift, skipping k8s version validator")
		return nil
	}
	version, err := k.K8sClient.GetVersion(k.ClusterData)
	if err != nil {
		return err
	}
	minVersion, err := strconv.ParseFloat(k.MinimumVersion, 64)
	if err != nil {
		return err
	}
	maxVersion, err := strconv.ParseFloat(k.MaximumVersion, 64)
	if err != nil {
		return err
	}
	currentVersion, err := strconv.ParseFloat(version, 64)
	if err != nil {
		return err
	}
	if currentVersion < minVersion {
		return fmt.Errorf("version %s is less than minimum supported version of %s", version, k.MinimumVersion)
	}
	if currentVersion > maxVersion {
		return fmt.Errorf("version %s is greater than maximum supported version of %s", version, k.MaximumVersion)
	}

	return nil
}
