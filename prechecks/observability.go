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
	"regexp"
	"strings"

	"github.com/dell/csm-deployment/utils"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	namespace             = "cert-manager"
	certManagerComponents = []string{"cert-manager", "cainjector", "webhook"}
	certManagerResources  = []string{"CertificateRequest", "Certificate", "Challenge", "ClusterIssuer", "Issuer", "Order"}
)

// K8sClientCertManagerInterface is the required interface for querying the k8s cluster for cert manager
//go:generate mockgen -destination=mocks/k8s_client_cert_manager_interface.go -package=mocks github.com/dell/csm-deployment/prechecks K8sClientCertManagerInterface
type K8sClientCertManagerInterface interface {
	GetCertManagerPods([]byte, string, string) (*corev1.PodList, error)
	GetAPIResource([]byte, string) (*metav1.APIResource, string, error)
}

// CertManagerValidator validates the required cert-manager components are installed
type CertManagerValidator struct {
	ClusterData []byte
	K8sClient   K8sClientCertManagerInterface
}

// Validate will check that all the expected cert-manager components exist and are running
func (k CertManagerValidator) Validate() error {
	for _, resource := range certManagerResources {
		_, groupVersion, err := k.K8sClient.GetAPIResource(k.ClusterData, resource)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("unable to find cert-manager CRD for %s", resource))
		}
		re := regexp.MustCompile(`v1alpha1`)
		hasAlphav1 := re.Match([]byte(groupVersion))
		if hasAlphav1 {
			return fmt.Errorf("has alphav1 of %s", resource)
		}
	}

	for _, component := range certManagerComponents {
		allPods, err := k.K8sClient.GetCertManagerPods(k.ClusterData, namespace, component)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("unable to get pods for %s", namespace))
		}
		if allPods == nil || len(allPods.Items) == 0 {
			return fmt.Errorf("no resources found for %s namepace. Please install cert-manager and the componets", namespace)
		}
		for _, pod := range allPods.Items {
			if pod.Status.Phase != corev1.PodRunning {
				return fmt.Errorf("the pod for the component, %s, has a status of %s", component, pod.Status.Phase)
			}
		}
	}

	return nil
}

// DriverSecretConfigsValidator validates the required driver secret configs in SecretConfigs are installed
type DriverSecretConfigsValidator struct {
	ModuleConfig    string
	ConfigFileNames []string
}

// Validate will check that secret config exist in namespace
func (k DriverSecretConfigsValidator) Validate() error {
	if strings.Contains(k.ModuleConfig, "karaviMetricsPowerflex.enabled=true") &&
		!(strings.Contains(k.ModuleConfig, "karaviMetricsPowerflex.driverConfig.filename") ||
			strings.Contains(k.ModuleConfig, "karaviMetricsPowerflex.driverConfig.data")) {
		return errors.New("missing PowerFlex secret config")
	}

	if strings.Contains(k.ModuleConfig, "karaviMetricsPowerstore.enabled=true") &&
		!(strings.Contains(k.ModuleConfig, "karaviMetricsPowerstore.driverConfig.filename") ||
			strings.Contains(k.ModuleConfig, "karaviMetricsPowerstore.driverConfig.data")) {
		return errors.New("missing PowerStore secret config")
	}

	if strings.Contains(k.ModuleConfig, "karaviMetricsPowerflex.driverConfig.filename") &&
		strings.Contains(k.ModuleConfig, "karaviMetricsPowerflex.driverConfig.data") {
		return errors.New("both filename and data are set for  PowerFlex secret config, set either but not both")
	}

	if strings.Contains(k.ModuleConfig, "karaviMetricsPowerstore.driverConfig.filename") &&
		strings.Contains(k.ModuleConfig, "karaviMetricsPowerstore.driverConfig.data") {
		return errors.New("both filename and data are set for  PowerStore secret config, set either but not both")
	}

	re := regexp.MustCompile(`[^\s]+`)
	for _, v := range re.FindAllString(k.ModuleConfig, -1) {
		if strings.Contains(v, ".driverConfig.filename") {
			configFile := strings.Split(v, "=")
			if len(configFile) != 2 {
				return errors.New("invalid ytt value format. It should be key1=value1")
			}
			if !utils.Find(k.ConfigFileNames, configFile[1]) {
				return fmt.Errorf("the filename %s for the module configuration value %s does not exist in configuration-files", configFile[1], configFile[0])
			}
		}
	}
	return nil
}
