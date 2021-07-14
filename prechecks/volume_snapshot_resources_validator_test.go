package prechecks

import (
	"errors"
	"testing"

	"github.com/dell/csm-deployment/prechecks/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_VolumeSnapshotResourcesValidator(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, VolumeSnapshotResourcesValidator, *gomock.Controller){
		"success": func(*testing.T) (bool, VolumeSnapshotResourcesValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			kubectl := mocks.NewMockKubectlExplainInterface(ctrl)
			kubectl.EXPECT().Explain(gomock.Any(), gomock.Any()).Times(3).Return([]byte(""), nil)

			snapshotValidator := VolumeSnapshotResourcesValidator{
				KubectlClient: kubectl,
			}

			return true, snapshotValidator, ctrl
		},
		"error - found v1alphav1 version of a crd": func(*testing.T) (bool, VolumeSnapshotResourcesValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			kubectl := mocks.NewMockKubectlExplainInterface(ctrl)
			kubectl.EXPECT().Explain(gomock.Any(), gomock.Any()).Times(1).Return([]byte("VERSION: snapshot.storage.k8s.io/v1alpha1"), nil)

			snapshotValidator := VolumeSnapshotResourcesValidator{
				KubectlClient: kubectl,
			}

			return false, snapshotValidator, ctrl
		},
		"error - kubectl returned error": func(*testing.T) (bool, VolumeSnapshotResourcesValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			kubectl := mocks.NewMockKubectlExplainInterface(ctrl)
			kubectl.EXPECT().Explain(gomock.Any(), gomock.Any()).Times(1).Return([]byte(""), errors.New("error"))

			snapshotValidator := VolumeSnapshotResourcesValidator{
				KubectlClient: kubectl,
			}

			return false, snapshotValidator, ctrl
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			expectSuccess, snapshotValidator, ctrl := tc(t)
			if expectSuccess {
				assert.NoError(t, snapshotValidator.Validate())
			} else {
				assert.Error(t, snapshotValidator.Validate())
			}
			ctrl.Finish()
		})
	}
}
