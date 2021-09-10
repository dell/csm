// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package prechecks

import (
	"github.com/dell/csm-deployment/model"
	"github.com/pkg/errors"
)

// SupportedModulesValidator validates whether a driver supports the given module or not
type SupportedModulesValidator struct {
	Modules    []model.ModuleType
	DriverType string
}

// Validate will check that all the expected cert-manager components exist and are running
func (k SupportedModulesValidator) Validate() error {
	switch k.DriverType {
	case model.ArrayTypePowerMax:
		for _, moduleType := range k.Modules {
			if !(moduleType.Name == model.ModuleTypeReplication ||
				moduleType.Name == model.ModuleTypeReverseProxy ||
				moduleType.Name == model.ModuleTypeAuthorization) {
				return errors.New("invalid module type given for the PowerMax driver")
			}
		}
	case model.ArrayTypePowerStore:
		for _, moduleType := range k.Modules {
			if !(moduleType.Name == model.ModuleTypeReplication ||
				moduleType.Name == model.ModuleTypeAuthorization ||
				moduleType.Name == model.ModuleTypeObservability) {
				return errors.New("invalid module type given for the PowerStore driver")
			}
		}
	case model.ArrayTypePowerFlex:
		for _, moduleType := range k.Modules {
			if !(moduleType.Name == model.ModuleTypePodMon ||
				moduleType.Name == model.ModuleTypeVgSnapShotter ||
				moduleType.Name == model.ModuleTypeObservability ||
				moduleType.Name == model.ModuleTypeAuthorization) {
				return errors.New("invalid module type given for the PowerFlex driver")
			}
		}
	case model.ArrayTypeUnity:
		for _, moduleType := range k.Modules {
			if moduleType.Name != model.ModuleTypePodMon {
				return errors.New("invalid module type given for the Unity driver")
			}
		}
	case model.ArrayTypePowerScale:
		for _, moduleType := range k.Modules {
			if moduleType.Name != model.ModuleTypeAuthorization {
				return errors.New("invalid module type given for the PowerScale driver")
			}
		}
	}

	return nil
}
