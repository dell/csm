package prechecks

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/pkg/errors"
)

var (
	snapshotResources = []string{"VolumeSnapshotClasses", "VolumeSnapshotContents", "VolumeSnapshots"}
)

// KubectlExplainInterface is the interface required to explain resources using kubectl
//go:generate mockgen -destination=mocks/kubectl_explain_interface.go -package=mocks github.com/dell/csm-deployment/prechecks KubectlExplainInterface
type KubectlExplainInterface interface {
	Explain(string, string) ([]byte, error)
}

// VolumeSnapshotResourcesValidator validates the required VolumeSnapshot CRDs and versions on the k8s cluster
type VolumeSnapshotResourcesValidator struct {
	ClusterData   []byte
	KubectlClient KubectlExplainInterface
}

// Validate will check that the expected CRD resources exist and that they are not of the version 'v1alphav1'
func (k VolumeSnapshotResourcesValidator) Validate() error {
	tmpFile, err := ioutil.TempFile("", "config")
	if err != nil {
		return err
	}

	_, err = tmpFile.Write(k.ClusterData)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	for _, resource := range snapshotResources {

		out, err := k.KubectlClient.Explain(tmpFile.Name(), resource)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("unable to find CRD for %s", resource))
		}

		re := regexp.MustCompile(`VERSION.*v1alpha1$`)
		hasAlphav1 := re.Match([]byte(out))
		if hasAlphav1 {
			return fmt.Errorf("has alphav1 of %s", resource)
		}
	}

	return nil
}
