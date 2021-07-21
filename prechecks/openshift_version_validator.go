package prechecks

import (
	"fmt"
	"log"
	"strconv"
)

// OpenshiftVersionValidator will validate the openshift version of the cluster
type OpenshiftVersionValidator struct {
	MinimumVersion string
	MaximumVersion string
	ClusterData    []byte
	K8sClient      K8sClientVersionInterface
}

// Validate will validate the version of the openshift cluster is between the min/max supported versions
func (k OpenshiftVersionValidator) Validate() error {
	isOpenshift, err := k.K8sClient.IsOpenShift(k.ClusterData)
	if err != nil {
		return err
	}
	if !isOpenshift {
		log.Printf("cluster is k8s, skipping openshift version validator")
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
