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

	"github.com/dell/csm-deployment/model"
	"github.com/stretchr/testify/assert"
)

func Test_GetDriverPrechecks(t *testing.T) {
	tests := map[string]func(t *testing.T) string{
		"success getting powerflex validations": func(*testing.T) string {
			return model.ArrayTypePowerFlex
		},
		"success getting powerstore validations": func(*testing.T) string {
			return model.ArrayTypePowerStore
		},
		"success getting powermax validations": func(*testing.T) string {
			return model.ArrayTypePowerMax
		},
		"success getting powerscale validations": func(*testing.T) string {
			return model.ArrayTypePowerScale
		},
		"success getting unity validations": func(*testing.T) string {
			return model.ArrayTypeUnity
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			driverType := tc(t)
			validators := PrecheckGetter{}.GetDriverPrechecks(driverType, nil, "", nil, nil)
			assert.True(t, len(validators) > 0)
		})
	}
}

func Test_GetModuleTypePrechecks(t *testing.T) {
	tests := map[string]func(t *testing.T) (string, map[string]string){
		"success getting observability validations": func(*testing.T) (string, map[string]string) {
			return model.ModuleTypeObservability, nil
		},
		"success getting authorization validations": func(*testing.T) (string, map[string]string) {
			return model.ModuleTypeAuthorization, nil
		},
		"success getting authorization validations for specific driver": func(*testing.T) (string, map[string]string) {
			return model.ModuleTypeAuthorization, map[string]string{"csidriver": model.ArrayTypePowerScale}
		},
		"success getting authorization validations with observability": func(*testing.T) (string, map[string]string) {
			return model.ModuleTypeAuthorization, map[string]string{model.ModuleTypeObservability: "true"}
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			moduleType, availableModules := tc(t)
			validators := PrecheckGetter{}.GetModuleTypePrechecks(moduleType, "", nil, []model.ConfigFile{{Name: "custom-config"}}, availableModules)
			assert.True(t, len(validators) > 0)
		})
	}
}
