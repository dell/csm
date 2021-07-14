package prechecks

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_SDCValidator(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, SDCValidator, *gomock.Controller){
		"success": func(*testing.T) (bool, SDCValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			nodeInfo := "[{\"host_name\":\"host_1\",\"installed_software\":{\"sdc\":\"enabled\"}},{\"host_name\":\"host_2\",\"installed_software\":{\"sdc\":\"enabled\"}}]"
			sdcValidator := SDCValidator{
				NodeInfo: nodeInfo,
			}

			return true, sdcValidator, ctrl
		},
		"error - host_2 doesn't have sdc enabled": func(*testing.T) (bool, SDCValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			nodeInfo := "[{\"host_name\":\"host_1\",\"installed_software\":{\"sdc\":\"enabled\"}},{\"host_name\":\"host_2\",\"installed_software\":{}}]"
			sdcValidator := SDCValidator{
				NodeInfo: nodeInfo,
			}

			return false, sdcValidator, ctrl
		},
		"error - invalid json format": func(*testing.T) (bool, SDCValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			nodeInfo := "invalid-json"
			sdcValidator := SDCValidator{
				NodeInfo: nodeInfo,
			}

			return false, sdcValidator, ctrl
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
