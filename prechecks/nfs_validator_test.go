// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package prechecks

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_NFSValidator(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, NFSValidator, *gomock.Controller){
		"success": func(*testing.T) (bool, NFSValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			nodeInfo := "[{\"host_name\":\"host_1\",\"installed_software\":{\"nfs\":\"enabled\"}},{\"host_name\":\"host_2\",\"installed_software\":{}}]"
			nfsValidator := NFSValidator{
				NodeInfo: nodeInfo,
			}

			return true, nfsValidator, ctrl
		},
		"error - hosts doesn't have nfs enabled": func(*testing.T) (bool, NFSValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			nodeInfo := "[{\"host_name\":\"host_1\",\"installed_software\":{}},{\"host_name\":\"host_2\",\"installed_software\":{}}]"
			nfsValidator := NFSValidator{
				NodeInfo: nodeInfo,
			}

			return false, nfsValidator, ctrl
		},
		"error - invalid json format": func(*testing.T) (bool, NFSValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			nodeInfo := "invalid-json"
			nfsValidator := NFSValidator{
				NodeInfo: nodeInfo,
			}

			return false, nfsValidator, ctrl
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			expectSuccess, nfsValidator, ctrl := tc(t)
			if expectSuccess {
				assert.NoError(t, nfsValidator.Validate())
			} else {
				assert.Error(t, nfsValidator.Validate())
			}
			ctrl.Finish()
		})
	}
}
