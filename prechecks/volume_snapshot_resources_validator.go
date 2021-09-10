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

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	snapshotResources = []string{"VolumeSnapshotClass", "VolumeSnapshotContent", "VolumeSnapshot"}
)

// K8sClientAPIResourceInterface is the required interface for querying the k8s cluster
//go:generate mockgen -destination=mocks/k8s_client_api_resource_interface.go -package=mocks github.com/dell/csm-deployment/prechecks K8sClientAPIResourceInterface
type K8sClientAPIResourceInterface interface {
	GetAPIResource([]byte, string) (*metav1.APIResource, string, error)
}

// VolumeSnapshotResourcesValidator validates the required VolumeSnapshot CRDs and versions on the k8s cluster
type VolumeSnapshotResourcesValidator struct {
	ClusterData       []byte
	K8sClient         K8sClientAPIResourceInterface
	OnlyVgSnapshotter bool
}

// Validate will check that the expected CRD resources exist and that they are not of the version 'v1alphav1'
func (k VolumeSnapshotResourcesValidator) Validate() error {
	snapshotResourcesList := snapshotResources
	if k.OnlyVgSnapshotter {
		snapshotResourcesList = []string{"DellCsiVolumeGroupSnapshot"}
	}
	for _, resource := range snapshotResourcesList {
		_, groupVersion, err := k.K8sClient.GetAPIResource(k.ClusterData, resource)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("unable to find CRD for %s", resource))
		}
		re := regexp.MustCompile(`v1alpha1`)
		hasAlphav1 := re.Match([]byte(groupVersion))
		if hasAlphav1 {
			if k.OnlyVgSnapshotter { // For now, it's okay for DellCsiVolumeGroupSnapshot to have storedVersions of v1alpha1
				continue
			}
			return fmt.Errorf("has alphav1 of %s", resource)
		}
	}

	return nil
}
