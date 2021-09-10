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
	"testing"

	"github.com/dell/csm-deployment/prechecks/mocks"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_K8sClientCertManagerInterface(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, CertManagerValidator, *gomock.Controller){
		"success": func(*testing.T) (bool, CertManagerValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			k8sclient := mocks.NewMockK8sClientCertManagerInterface(ctrl)

			k8sclient.EXPECT().GetAPIResource(gomock.Any(), gomock.Any()).Times(6).Return(&metav1.APIResource{}, " apiextensions.k8s.io/v1", nil)

			pod := corev1.Pod{}
			pod.Status.Phase = corev1.PodRunning
			k8sclient.EXPECT().GetCertManagerPods(gomock.Any(), gomock.Any(), gomock.Any()).Times(3).Return(&corev1.PodList{Items: []corev1.Pod{pod}}, nil)

			certManagerValidator := CertManagerValidator{
				K8sClient: k8sclient,
			}

			return true, certManagerValidator, ctrl
		},
		"error one pod is still in pending": func(*testing.T) (bool, CertManagerValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			k8sclient := mocks.NewMockK8sClientCertManagerInterface(ctrl)
			k8sclient.EXPECT().GetAPIResource(gomock.Any(), gomock.Any()).Times(6).Return(&metav1.APIResource{}, " apiextensions.k8s.io/v1", nil)

			pod := corev1.Pod{}
			pod.Status.Phase = corev1.PodRunning
			k8sclient.EXPECT().GetCertManagerPods(gomock.Any(), gomock.Any(), gomock.Any()).Times(2).Return(&corev1.PodList{Items: []corev1.Pod{pod}}, nil)
			pod.Status.Phase = corev1.PodPending
			k8sclient.EXPECT().GetCertManagerPods(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&corev1.PodList{Items: []corev1.Pod{pod}}, nil)

			certManagerValidator := CertManagerValidator{
				K8sClient: k8sclient,
			}

			return false, certManagerValidator, ctrl
		},
		"error component missing": func(*testing.T) (bool, CertManagerValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			k8sclient := mocks.NewMockK8sClientCertManagerInterface(ctrl)
			k8sclient.EXPECT().GetAPIResource(gomock.Any(), gomock.Any()).Times(6).Return(&metav1.APIResource{}, " apiextensions.k8s.io/v1", nil)
			k8sclient.EXPECT().GetCertManagerPods(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&corev1.PodList{}, errors.New("error"))

			certManagerValidator := CertManagerValidator{
				K8sClient: k8sclient,
			}

			return false, certManagerValidator, ctrl
		},
		"error empty resources in namespace": func(*testing.T) (bool, CertManagerValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			k8sclient := mocks.NewMockK8sClientCertManagerInterface(ctrl)
			k8sclient.EXPECT().GetAPIResource(gomock.Any(), gomock.Any()).Times(6).Return(&metav1.APIResource{}, " apiextensions.k8s.io/v1", nil)
			k8sclient.EXPECT().GetCertManagerPods(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&corev1.PodList{Items: []corev1.Pod{}}, nil)

			certManagerValidator := CertManagerValidator{
				K8sClient: k8sclient,
			}

			return false, certManagerValidator, ctrl
		},
		"error found v1alphav1 version of a crd": func(*testing.T) (bool, CertManagerValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			k8sclient := mocks.NewMockK8sClientCertManagerInterface(ctrl)
			k8sclient.EXPECT().GetAPIResource(gomock.Any(), gomock.Any()).Times(1).Return(&metav1.APIResource{}, " apiextensions.k8s.io/v1alpha1", nil)

			certManagerValidator := CertManagerValidator{
				K8sClient: k8sclient,
			}

			return false, certManagerValidator, ctrl
		},
		"error k8sclient returned error": func(*testing.T) (bool, CertManagerValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			k8sclient := mocks.NewMockK8sClientCertManagerInterface(ctrl)
			k8sclient.EXPECT().GetAPIResource(gomock.Any(), gomock.Any()).Times(1).Return(nil, "", errors.New("error"))

			certManagerValidator := CertManagerValidator{
				K8sClient: k8sclient,
			}

			return false, certManagerValidator, ctrl
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			expectSuccess, certManagerValidator, ctrl := tc(t)
			if expectSuccess {
				assert.NoError(t, certManagerValidator.Validate())
			} else {
				assert.Error(t, certManagerValidator.Validate())
			}
			ctrl.Finish()
		})
	}
}

func Test_DriverSecretConfigsValidator(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, DriverSecretConfigsValidator, *gomock.Controller){
		"success": func(*testing.T) (bool, DriverSecretConfigsValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			configFileNames := []string{"vxflexos-config"}
			moduleConfig := `
			karaviMetricsPowerflex.enabled=true
			karaviMetricsPowerflex.driverConfig.filename=vxflexos-config 
			karaviMetricsPowerstore.enabled=true
			karaviMetricsPowerstore.driverConfig.data=fakePowerStoreConfig
			`

			driverSecretConfigsValidator := DriverSecretConfigsValidator{
				ConfigFileNames: configFileNames,
				ModuleConfig:    moduleConfig,
			}

			return true, driverSecretConfigsValidator, ctrl
		},
		"fail due to missing powerflex secret config": func(*testing.T) (bool, DriverSecretConfigsValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			configFileNames := []string{"vxflexos-config"}
			moduleConfig := `
			karaviMetricsPowerflex.enabled=true
			`

			driverSecretConfigsValidator := DriverSecretConfigsValidator{
				ConfigFileNames: configFileNames,
				ModuleConfig:    moduleConfig,
			}

			return false, driverSecretConfigsValidator, ctrl
		},
		"fail due to missing powerstore secret config": func(*testing.T) (bool, DriverSecretConfigsValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			configFileNames := []string{"vxflexos-config"}
			moduleConfig := `
			karaviMetricsPowerstore.enabled=true
			`

			driverSecretConfigsValidator := DriverSecretConfigsValidator{
				ConfigFileNames: configFileNames,
				ModuleConfig:    moduleConfig,
			}

			return false, driverSecretConfigsValidator, ctrl
		},
		"fail due to both powerflex filename and data specified": func(*testing.T) (bool, DriverSecretConfigsValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			configFileNames := []string{"vxflexos-config"}
			moduleConfig := `
			karaviMetricsPowerflex.enabled=true
			karaviMetricsPowerflex.driverConfig.filename=file
			karaviMetricsPowerflex.driverConfig.data=data
			`

			driverSecretConfigsValidator := DriverSecretConfigsValidator{
				ConfigFileNames: configFileNames,
				ModuleConfig:    moduleConfig,
			}

			return false, driverSecretConfigsValidator, ctrl
		},
		"fail due to both powerstore filename and data specified": func(*testing.T) (bool, DriverSecretConfigsValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			configFileNames := []string{"vxflexos-config"}
			moduleConfig := `
			karaviMetricsPowerstore.enabled=true
			karaviMetricsPowerstore.driverConfig.filename=file
			karaviMetricsPowerstore.driverConfig.data=data
			`

			driverSecretConfigsValidator := DriverSecretConfigsValidator{
				ConfigFileNames: configFileNames,
				ModuleConfig:    moduleConfig,
			}

			return false, driverSecretConfigsValidator, ctrl
		},
		"fail due to missing the set value filename in configuration-files db": func(*testing.T) (bool, DriverSecretConfigsValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			configFileNames := []string{""}
			moduleConfig := `
			karaviMetricsPowerflex.enabled=true
			karaviMetricsPowerflex.driverConfig.filename=vxflexos-config 
			karaviMetricsPowerstore.enabled=true
			karaviMetricsPowerstore.driverConfig.data=fakePowerStoreConfig
			`

			driverSecretConfigsValidator := DriverSecretConfigsValidator{
				ConfigFileNames: configFileNames,
				ModuleConfig:    moduleConfig,
			}

			return false, driverSecretConfigsValidator, ctrl
		},
		"fail due to setting both filename and data values instead of one": func(*testing.T) (bool, DriverSecretConfigsValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			configFileNames := []string{""}
			moduleConfig := `
			karaviMetricsPowerflex.enabled=true
			karaviMetricsPowerflex.driverConfig.filename=vxflexos-config
			karaviMetricsPowerflex.driverConfig.data=fakePowerflexConfig
			karaviMetricsPowerstore.enabled=true
			karaviMetricsPowerstore.driverConfig.data=fakePowerStoreConfig
			`

			driverSecretConfigsValidator := DriverSecretConfigsValidator{
				ConfigFileNames: configFileNames,
				ModuleConfig:    moduleConfig,
			}

			return false, driverSecretConfigsValidator, ctrl
		},
		"fail due to  wrong key=value format": func(*testing.T) (bool, DriverSecretConfigsValidator, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			configFileNames := []string{""}
			moduleConfig := `
			karaviMetricsPowerflex.enabled=true
			karaviMetricsPowerflex.driverConfig.filename:vxflexos-config
			karaviMetricsPowerstore.enabled=true
			karaviMetricsPowerstore.driverConfig.data=fakePowerStoreConfig
			`

			driverSecretConfigsValidator := DriverSecretConfigsValidator{
				ConfigFileNames: configFileNames,
				ModuleConfig:    moduleConfig,
			}

			return false, driverSecretConfigsValidator, ctrl
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			expectSuccess, authorizationValidator, ctrl := tc(t)
			if expectSuccess {
				assert.NoError(t, authorizationValidator.Validate())
			} else {
				assert.Error(t, authorizationValidator.Validate())
			}
			ctrl.Finish()
		})
	}
}
