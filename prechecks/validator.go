// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package prechecks

import (
	"fmt"

	"github.com/dell/csm-deployment/k8s"
	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
)

// Validator is the interface that all validation prechecks must implement
type Validator interface {
	Validate() error
}

const (
	// K8sMinimumSupportedVersion is the minimum supported version for k8s
	K8sMinimumSupportedVersion = "1.19"
	// K8sMaximumSupportedVersion is the maximum supported version for k8s
	K8sMaximumSupportedVersion = "1.22"
	// OpenshiftMinimumSupportedVersion is the minimum supported version for openshift
	OpenshiftMinimumSupportedVersion = "4.6"
	// OpenshiftMaximumSupportedVersion is the maximum supported version for openshift
	OpenshiftMaximumSupportedVersion = "4.7"
)

var (
	// AuthSupportedDrivers is the List of drivers currently supported by CSM Authorization
	AuthSupportedDrivers = []string{model.ArrayTypePowerMax, model.ArrayTypePowerFlex, model.ArrayTypePowerScale}
)

// PrecheckGetter will get a list of validators for different resources (e.g. drivers, modules)
type PrecheckGetter struct{}

// GetDriverPrechecks will return a list of prechecks for the specific driver and version
func (p PrecheckGetter) GetDriverPrechecks(driverType string, clusterData []byte, clusterNodeDetails string, modules []model.ModuleType, logger echo.Logger) []Validator {
	validators := make([]Validator, 0)

	// common prechecks for all drivers
	validators = append(validators, K8sVersionValidator{
		MinimumVersion: K8sMinimumSupportedVersion,
		MaximumVersion: K8sMaximumSupportedVersion,
		ClusterData:    clusterData,
		K8sClient:      k8s.Client{},
		Logger:         logger,
	})
	validators = append(validators, OpenshiftVersionValidator{
		MinimumVersion: OpenshiftMinimumSupportedVersion,
		MaximumVersion: OpenshiftMaximumSupportedVersion,
		ClusterData:    clusterData,
		K8sClient:      k8s.Client{},
		Logger:         logger,
	})
	validators = append(validators, VolumeSnapshotResourcesValidator{
		ClusterData: clusterData,
		K8sClient:   k8s.Client{},
	})

	switch driverType {
	case model.ArrayTypePowerFlex:
		validators = append(validators, SDCValidator{
			NodeInfo: clusterNodeDetails,
		}, SupportedModulesValidator{
			DriverType: driverType,
			Modules:    modules,
		})
	case model.ArrayTypePowerScale:
		validators = append(validators, NFSValidator{
			NodeInfo: clusterNodeDetails,
		}, SupportedModulesValidator{
			DriverType: driverType,
			Modules:    modules,
		})
	case model.ArrayTypePowerStore:
		validators = append(validators, ISCSIValidator{
			NodeInfo: clusterNodeDetails,
		}, NFSValidator{
			NodeInfo: clusterNodeDetails,
		}, SupportedModulesValidator{
			DriverType: driverType,
			Modules:    modules,
		})
	case model.ArrayTypeUnity:
		validators = append(validators, ISCSIValidator{
			NodeInfo: clusterNodeDetails,
		}, NFSValidator{
			NodeInfo: clusterNodeDetails,
		}, SupportedModulesValidator{
			DriverType: driverType,
			Modules:    modules,
		})
	case model.ArrayTypePowerMax:
		validators = append(validators, ISCSIValidator{
			NodeInfo: clusterNodeDetails,
		}, NFSValidator{
			NodeInfo: clusterNodeDetails,
		}, SupportedModulesValidator{
			DriverType: driverType,
			Modules:    modules,
		})
	}
	return validators
}

// GetModuleTypePrechecks will return a list of prechecks for the specific module
func (p PrecheckGetter) GetModuleTypePrechecks(moduleType, moduleConfig string, clusterData []byte, cfs []model.ConfigFile, availableModules map[string]string) []Validator {
	var filenames []string
	for _, cf := range cfs {
		filenames = append(filenames, cf.Name)
	}

	validators := make([]Validator, 0)

	switch moduleType {
	case model.ModuleTypeObservability:
		validators = append(validators, CertManagerValidator{
			ClusterData: clusterData,
			K8sClient:   k8s.Client{},
		})

		// standalone, check that the driver secret dependencies all exist
		if _, ok := availableModules["csidriver"]; !ok {
			validators = append(validators, DriverSecretConfigsValidator{
				ConfigFileNames: filenames,
				ModuleConfig:    moduleConfig,
			})
		}

	case model.ModuleTypeAuthorization:
		if _, ok := availableModules["csidriver"]; ok {
			if !utils.Find(AuthSupportedDrivers, availableModules["csidriver"]) {
				validators = append(validators, AuthorizationValidator{
					Skip: struct {
						Cond bool
						Msg  string
					}{Cond: true, Msg: fmt.Sprintf("csm authorization does not currently support %s", availableModules["csidriver"])},
				})
				break
			}
		} else {
			if _, ok := availableModules[model.ModuleTypeObservability]; !ok {
				validators = append(validators, AuthorizationValidator{
					Skip: struct {
						Cond bool
						Msg  string
					}{Cond: true, Msg: "authorization cannot be installed as standalone"},
				})
				break
			}
		}

		validators = append(validators, AuthorizationValidator{
			Skip: struct {
				Cond bool
				Msg  string
			}{Cond: false},
			ConfigFileNames: filenames,
			ModuleConfig:    moduleConfig,
		})
	case model.ModuleTypeVgSnapShotter:
		validators = append(validators, VolumeSnapshotResourcesValidator{
			ClusterData:       clusterData,
			K8sClient:         k8s.Client{},
			OnlyVgSnapshotter: true,
		})
	}

	return validators
}
