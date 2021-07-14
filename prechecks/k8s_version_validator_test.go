package prechecks

import (
	"errors"
	"testing"

	"github.com/dell/csm-deployment/prechecks/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_K8sVersionValidator(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, K8sVersionValidator, *gomock.Controller){
		"success": func(*testing.T) (bool, K8sVersionValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			versionInterface := mocks.NewMockK8sClientVersionInterface(ctrl)
			versionInterface.EXPECT().IsOpenShift(gomock.Any()).Times(1).Return(false, nil)
			versionInterface.EXPECT().GetVersion(gomock.Any()).Times(1).Return("1.20", nil)

			versionValidator := K8sVersionValidator{
				MinimumVersion: "1.19",
				MaximumVersion: "1.21",
				K8sClient:      versionInterface,
			}

			return true, versionValidator, ctrl
		},
		"success - at minimum version": func(*testing.T) (bool, K8sVersionValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			versionInterface := mocks.NewMockK8sClientVersionInterface(ctrl)
			versionInterface.EXPECT().IsOpenShift(gomock.Any()).Times(1).Return(false, nil)
			versionInterface.EXPECT().GetVersion(gomock.Any()).Times(1).Return("1.19", nil)

			versionValidator := K8sVersionValidator{
				MinimumVersion: "1.19",
				MaximumVersion: "1.21",
				K8sClient:      versionInterface,
			}

			return true, versionValidator, ctrl
		},
		"success - at maximum version": func(*testing.T) (bool, K8sVersionValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			versionInterface := mocks.NewMockK8sClientVersionInterface(ctrl)
			versionInterface.EXPECT().IsOpenShift(gomock.Any()).Times(1).Return(false, nil)
			versionInterface.EXPECT().GetVersion(gomock.Any()).Times(1).Return("1.21", nil)

			versionValidator := K8sVersionValidator{
				MinimumVersion: "1.19",
				MaximumVersion: "1.21",
				K8sClient:      versionInterface,
			}

			return true, versionValidator, ctrl
		},
		"success - skip openshift": func(*testing.T) (bool, K8sVersionValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			versionInterface := mocks.NewMockK8sClientVersionInterface(ctrl)
			versionInterface.EXPECT().IsOpenShift(gomock.Any()).Times(1).Return(true, nil)

			versionValidator := K8sVersionValidator{
				K8sClient: versionInterface,
			}

			return true, versionValidator, ctrl
		},
		"error - below minimum version": func(*testing.T) (bool, K8sVersionValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			versionInterface := mocks.NewMockK8sClientVersionInterface(ctrl)
			versionInterface.EXPECT().IsOpenShift(gomock.Any()).Times(1).Return(false, nil)
			versionInterface.EXPECT().GetVersion(gomock.Any()).Times(1).Return("1.18", nil)

			versionValidator := K8sVersionValidator{
				MinimumVersion: "1.19",
				MaximumVersion: "1.21",
				K8sClient:      versionInterface,
			}

			return false, versionValidator, ctrl
		},
		"error - above maximum minimum version": func(*testing.T) (bool, K8sVersionValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			versionInterface := mocks.NewMockK8sClientVersionInterface(ctrl)
			versionInterface.EXPECT().IsOpenShift(gomock.Any()).Times(1).Return(false, nil)
			versionInterface.EXPECT().GetVersion(gomock.Any()).Times(1).Return("1.22", nil)

			versionValidator := K8sVersionValidator{
				MinimumVersion: "1.19",
				MaximumVersion: "1.21",
				K8sClient:      versionInterface,
			}

			return false, versionValidator, ctrl
		},
		"error - checking openshift": func(*testing.T) (bool, K8sVersionValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			versionInterface := mocks.NewMockK8sClientVersionInterface(ctrl)
			versionInterface.EXPECT().IsOpenShift(gomock.Any()).Times(1).Return(false, errors.New("error"))

			versionValidator := K8sVersionValidator{
				K8sClient: versionInterface,
			}

			return false, versionValidator, ctrl
		},
		"error - getting version": func(*testing.T) (bool, K8sVersionValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			versionInterface := mocks.NewMockK8sClientVersionInterface(ctrl)
			versionInterface.EXPECT().IsOpenShift(gomock.Any()).Times(1).Return(false, nil)
			versionInterface.EXPECT().GetVersion(gomock.Any()).Times(1).Return("", errors.New("error"))

			versionValidator := K8sVersionValidator{
				K8sClient: versionInterface,
			}

			return false, versionValidator, ctrl
		},
		"error - invalid min version": func(*testing.T) (bool, K8sVersionValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			versionInterface := mocks.NewMockK8sClientVersionInterface(ctrl)
			versionInterface.EXPECT().IsOpenShift(gomock.Any()).Times(1).Return(false, nil)
			versionInterface.EXPECT().GetVersion(gomock.Any()).Times(1).Return("1.22", nil)

			versionValidator := K8sVersionValidator{
				MinimumVersion: "invalid-version",
				MaximumVersion: "1.21",
				K8sClient:      versionInterface,
			}

			return false, versionValidator, ctrl
		},
		"error - invalid max version": func(*testing.T) (bool, K8sVersionValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			versionInterface := mocks.NewMockK8sClientVersionInterface(ctrl)
			versionInterface.EXPECT().IsOpenShift(gomock.Any()).Times(1).Return(false, nil)
			versionInterface.EXPECT().GetVersion(gomock.Any()).Times(1).Return("1.22", nil)

			versionValidator := K8sVersionValidator{
				MinimumVersion: "1.19",
				MaximumVersion: "invalid-version",
				K8sClient:      versionInterface,
			}

			return false, versionValidator, ctrl
		},
		"error - invalid version": func(*testing.T) (bool, K8sVersionValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			versionInterface := mocks.NewMockK8sClientVersionInterface(ctrl)
			versionInterface.EXPECT().IsOpenShift(gomock.Any()).Times(1).Return(false, nil)
			versionInterface.EXPECT().GetVersion(gomock.Any()).Times(1).Return("invalid-version", nil)

			versionValidator := K8sVersionValidator{
				MinimumVersion: "1.19",
				MaximumVersion: "1.21",
				K8sClient:      versionInterface,
			}

			return false, versionValidator, ctrl
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			expectSuccess, versionValidator, ctrl := tc(t)
			if expectSuccess {
				assert.NoError(t, versionValidator.Validate())
			} else {
				assert.Error(t, versionValidator.Validate())
			}
			ctrl.Finish()
		})
	}
}
