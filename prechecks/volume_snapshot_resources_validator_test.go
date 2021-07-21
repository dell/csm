package prechecks

import (
	"errors"
	"testing"

	"github.com/dell/csm-deployment/prechecks/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_VolumeSnapshotResourcesValidator(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, VolumeSnapshotResourcesValidator, *gomock.Controller){
		"success": func(*testing.T) (bool, VolumeSnapshotResourcesValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			k8sclient := mocks.NewMockK8sClientExplainInterface(ctrl)
			k8sclient.EXPECT().Explain(gomock.Any(), gomock.Any()).Times(3).Return(&metav1.APIResource{}, "snapshot.storage.k8s.io/v1", nil)

			snapshotValidator := VolumeSnapshotResourcesValidator{
				K8sClient: k8sclient,
			}

			return true, snapshotValidator, ctrl
		},
		"error - found v1alphav1 version of a crd": func(*testing.T) (bool, VolumeSnapshotResourcesValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			k8sclient := mocks.NewMockK8sClientExplainInterface(ctrl)
			k8sclient.EXPECT().Explain(gomock.Any(), gomock.Any()).Times(1).Return(&metav1.APIResource{}, "snapshot.storage.k8s.io/v1alpha1", nil)

			snapshotValidator := VolumeSnapshotResourcesValidator{
				K8sClient: k8sclient,
			}

			return false, snapshotValidator, ctrl
		},
		"error - k8sclient returned error": func(*testing.T) (bool, VolumeSnapshotResourcesValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			k8sclient := mocks.NewMockK8sClientExplainInterface(ctrl)
			k8sclient.EXPECT().Explain(gomock.Any(), gomock.Any()).Times(1).Return(nil, "", errors.New("error"))

			snapshotValidator := VolumeSnapshotResourcesValidator{
				K8sClient: k8sclient,
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
