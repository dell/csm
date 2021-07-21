package prechecks

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_ISCSIValidator(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, ISCSIValidator, *gomock.Controller){
		"success": func(*testing.T) (bool, ISCSIValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			nodeInfo := "[{\"host_name\":\"host_1\",\"installed_software\":{\"iscsi\":\"enabled\"}},{\"host_name\":\"host_2\",\"installed_software\":{\"iscsi\":\"enabled\"}}]"
			iscsiValidator := ISCSIValidator{
				NodeInfo: nodeInfo,
			}

			return true, iscsiValidator, ctrl
		},
		"error - host_2 doesn't have iscsi enabled": func(*testing.T) (bool, ISCSIValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			nodeInfo := "[{\"host_name\":\"host_1\",\"installed_software\":{\"iscsi\":\"enabled\"}},{\"host_name\":\"host_2\",\"installed_software\":{}}]"
			iscsiValidator := ISCSIValidator{
				NodeInfo: nodeInfo,
			}

			return false, iscsiValidator, ctrl
		},
		"error - invalid json format": func(*testing.T) (bool, ISCSIValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			nodeInfo := "invalid-json"
			iscsiValidator := ISCSIValidator{
				NodeInfo: nodeInfo,
			}

			return false, iscsiValidator, ctrl
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
