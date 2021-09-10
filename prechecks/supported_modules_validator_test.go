// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

// Package  prechecks for application creation
package prechecks

import (
	"github.com/dell/csm-deployment/model"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_SupportedModulesValidator(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, SupportedModulesValidator, *gomock.Controller){
		"success": func(*testing.T) (bool, SupportedModulesValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			moduleTypes := []model.ModuleType{{
				Name: model.ModuleTypePodMon,
			}}

			SupportedModulesValidator := SupportedModulesValidator{
				Modules:    moduleTypes,
				DriverType: model.ArrayTypeUnity,
			}

			return true, SupportedModulesValidator, ctrl
		},
		"failed due to wrong module type": func(*testing.T) (bool, SupportedModulesValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			moduleTypes := []model.ModuleType{{
				Name: model.ModuleTypeVgSnapShotter,
			}}

			SupportedModulesValidator := SupportedModulesValidator{
				Modules:    moduleTypes,
				DriverType: model.ArrayTypeUnity,
			}

			return false, SupportedModulesValidator, ctrl
		},
		"Success with replication and reverse proxy module for pmax": func(*testing.T) (bool, SupportedModulesValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			moduleTypes := []model.ModuleType{{
				Name: model.ModuleTypeReverseProxy,
			}, {
				Name: model.ModuleTypeReplication,
			}, {
				Name: model.ModuleTypeAuthorization,
			}}

			SupportedModulesValidator := SupportedModulesValidator{
				Modules:    moduleTypes,
				DriverType: model.ArrayTypePowerMax,
			}

			return true, SupportedModulesValidator, ctrl
		},
		"Failure with invalid modules for pmax": func(*testing.T) (bool, SupportedModulesValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			moduleTypes := []model.ModuleType{{
				Name: model.ModuleTypeVgSnapShotter,
			}, {
				Name: model.ModuleTypeReplication,
			}}

			SupportedModulesValidator := SupportedModulesValidator{
				Modules:    moduleTypes,
				DriverType: model.ArrayTypePowerMax,
			}

			return false, SupportedModulesValidator, ctrl
		},
		"Success with valid modules for pstore": func(*testing.T) (bool, SupportedModulesValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			moduleTypes := []model.ModuleType{{
				Name: model.ModuleTypeReplication,
			}, {
				Name: model.ModuleTypeAuthorization,
			}, {
				Name: model.ModuleTypeObservability,
			}}

			SupportedModulesValidator := SupportedModulesValidator{
				Modules:    moduleTypes,
				DriverType: model.ArrayTypePowerStore,
			}

			return true, SupportedModulesValidator, ctrl
		},
		"Failure with invalid modules for pstore": func(*testing.T) (bool, SupportedModulesValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			moduleTypes := []model.ModuleType{{
				Name: model.ModuleTypePodMon,
			}}

			SupportedModulesValidator := SupportedModulesValidator{
				Modules:    moduleTypes,
				DriverType: model.ArrayTypePowerStore,
			}

			return false, SupportedModulesValidator, ctrl
		},
		"Failure with invalid modules for isilon": func(*testing.T) (bool, SupportedModulesValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			moduleTypes := []model.ModuleType{{
				Name: model.ModuleTypeVgSnapShotter,
			}}

			SupportedModulesValidator := SupportedModulesValidator{
				Modules:    moduleTypes,
				DriverType: model.ArrayTypePowerScale,
			}

			return false, SupportedModulesValidator, ctrl
		},
		"Success with valid modules for isilon": func(*testing.T) (bool, SupportedModulesValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			moduleTypes := []model.ModuleType{{
				Name: model.ModuleTypeAuthorization,
			}}

			SupportedModulesValidator := SupportedModulesValidator{
				Modules:    moduleTypes,
				DriverType: model.ArrayTypePowerScale,
			}

			return true, SupportedModulesValidator, ctrl
		},
		"Success with valid modules for powerflex": func(*testing.T) (bool, SupportedModulesValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			moduleTypes := []model.ModuleType{{
				Name: model.ModuleTypeVgSnapShotter,
			}, {
				Name: model.ModuleTypePodMon,
			}, {
				Name: model.ModuleTypeObservability,
			}, {
				Name: model.ModuleTypeAuthorization,
			}}

			SupportedModulesValidator := SupportedModulesValidator{
				Modules:    moduleTypes,
				DriverType: model.ArrayTypePowerFlex,
			}

			return true, SupportedModulesValidator, ctrl
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			expectSuccess, supportedModulesValidator, ctrl := tc(t)
			if expectSuccess {
				assert.NoError(t, supportedModulesValidator.Validate())
			} else {
				assert.Error(t, supportedModulesValidator.Validate())
			}
			ctrl.Finish()
		})
	}
}
