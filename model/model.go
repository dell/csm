// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package model

import (
	"fmt"

	"gorm.io/gorm"
)

const (
	// OrchestratorTypeKubernetes - placeholder for constant Kubernetes
	OrchestratorTypeKubernetes string = "k8s"

	// OrchestratorTypeOpenshift - placeholder for constant openshift
	OrchestratorTypeOpenshift string = "openshift"

	// ClusterStatusConnected - placeholder for constant connected
	ClusterStatusConnected string = "connected"

	// ClusterStatusDisconnected - placeholder for constant disconnected
	ClusterStatusDisconnected string = "disconnected"

	// DriverStatusOperational - placeholder for constant operational
	DriverStatusOperational string = "operational"

	// DriverStatusFailing - placeholder for constant failing
	DriverStatusFailing string = "failing"

	// ModuleTypeStandalone - placeholder for constant standalone
	ModuleTypeStandalone string = "standalone"

	// ModuleTypeSidecar - placeholder for constant sidecar
	ModuleTypeSidecar string = "sidecar"

	// ModuleTypeReplication - placeholder for replication constant
	ModuleTypeReplication string = "replication"

	// ModuleTypeObservability - placeholder for constant observability
	ModuleTypeObservability string = "observability"

	// ModuleTypePodMon - placeholder for constant podmon
	ModuleTypePodMon string = "podmon"

	// ModuleTypeVgSnapShotter - placeholder for constant vgsnapshotter
	ModuleTypeVgSnapShotter string = "vgsnapshotter"

	// ModuleTypeAuthorization - placeholder for constant authorization
	ModuleTypeAuthorization string = "authorization"

	// ModuleTypeReverseProxy - placeholder for constant csireverseproxy
	ModuleTypeReverseProxy string = "csireverseproxy"

	// TaskStatusInProgress - placeholder for constant task in progress
	TaskStatusInProgress string = "in progress"

	// TaskStatusCompleted - placeholder for constant task completed
	TaskStatusCompleted string = "completed"

	// TaskStatusPrompting - placeholder for constant task prompting
	TaskStatusPrompting string = "prompting"

	// TaskStatusFailed - placeholder for constant task failed
	TaskStatusFailed string = "failed"

	// TaskTypeInstall - placeholder for constant task install
	TaskTypeInstall string = "install"

	// TaskTypeUpdate - placeholder for constant task update
	TaskTypeUpdate string = "update"

	// TaskTypeDelete - placeholder for constant task delete
	TaskTypeDelete string = "delete"

	// ArrayTypePowerFlex - placeholder for constant powerflex
	ArrayTypePowerFlex string = "powerflex"

	// ArrayTypePowerMax - placeholder for constant powermax
	ArrayTypePowerMax string = "powermax"

	// ArrayTypePowerScale - placeholder for constant PowerScale
	ArrayTypePowerScale string = "isilon"

	// ArrayTypeUnity - placeholder for constant unity
	ArrayTypeUnity string = "unity"

	// ArrayTypePowerStore - placeholder for constant powerstore
	ArrayTypePowerStore string = "powerstore"

	// ReplicationNamespace - replication sidecar installation namespace
	ReplicationNamespace string = "dell-replication-controller"
	//ObservabilityNamespace - placeholder for constant observability namespace
	ObservabilityNamespace string = "csm-observability"
)

// User - Placeholder for User information
type User struct {
	gorm.Model
	Username string `gorm:"unique_index;not null"`
	Password string `gorm:"not null"`
	Admin    bool
}

// Cluster - Placeholder for Cluster information
type Cluster struct {
	gorm.Model              // This already contains ID field
	ClusterName      string `gorm:"unique;not null"`
	ConfigFileData   []byte `gorm:"not null"`
	OrchestratorType string `gorm:"not null"`
	Status           string `gorm:"not null"`
	K8sVersion       string
	Applications     []Application
	ClusterDetails   ClusterDetails
}

// BeforeDelete hook defined for cascade delete
func (c *Cluster) BeforeDelete(tx *gorm.DB) error {
	if len(c.Applications) > 0 {
		return fmt.Errorf("cluster is in use by applications")
	}
	return tx.Model(&ClusterDetails{}).Where("cluster_id = ?", c.ID).Unscoped().Delete(&ClusterDetails{}).Error
}

// ClusterDetails - Placeholder for Cluster and Nodes information
type ClusterDetails struct {
	gorm.Model // This already contains ID field
	ClusterID  string
	Nodes      string
}

// ConfigFile - Placeholder forconfig file type
type ConfigFile struct {
	gorm.Model            // This already contains ID field
	Name           string `gorm:"unique;not null"`
	ConfigFileData []byte `gorm:"not null"`
}

// Application - Placeholder for Application
type Application struct {
	gorm.Model
	Name                string `gorm:"uniqueIndex:name_id"`
	Status              string `gorm:"not null"`
	ClusterID           uint   `gorm:"uniqueIndex:name_id"`
	DriverTypeID        uint
	ModuleTypes         []ModuleType   `gorm:"many2many:application_module_types;"`
	StorageArrays       []StorageArray `gorm:"many2many:application_storage_arrays;"`
	Tasks               []Task
	DriverConfiguration string
	ModuleConfiguration string
	// TODO: These can be deleted.
	ApplicationOutput string
}

// BeforeDelete hook defined for cascade delete
func (a *Application) BeforeDelete(tx *gorm.DB) error {
	err := tx.Table("application_module_types").Where("application_id = ?", a.ID).Unscoped().Delete(&ModuleType{}).Error
	if err != nil {
		return err
	}
	err = tx.Table("application_storage_arrays").Where("application_id = ?", a.ID).Unscoped().Delete(&StorageArray{}).Error
	if err != nil {
		return err
	}
	return tx.Model(&Task{}).Where("application_id = ?", a.ID).Unscoped().Delete(&Task{}).Error
}

// ApplicationStateChange - Placeholder for Application State Change
type ApplicationStateChange struct {
	gorm.Model
	ApplicationID       uint
	ClusterID           uint
	DriverTypeID        uint
	ModuleTypes         []ModuleType   `gorm:"many2many:application_state_change_module_types;"`
	StorageArrays       []StorageArray `gorm:"many2many:application_state_change_storage_arrays;"`
	Template            []byte
	DriverConfiguration string
	ModuleConfiguration string
}

// StorageArrayType - Placeholder for Storage array type
type StorageArrayType struct {
	gorm.Model
	Name string
}

// DriverType - Placeholder for Driver type
type DriverType struct {
	gorm.Model         // This already contains ID field
	Version            string
	StorageArrayTypeID uint
	StorageArrayType   StorageArrayType
}

// ModuleType holds details about module used in application
type ModuleType struct {
	gorm.Model // This already contains ID field
	Name       string
	Version    string
	Standalone bool
}

// Task - Placeholder for Task
type Task struct {
	gorm.Model
	Status        string `gorm:"not null"`
	TaskType      string `gorm:"not null"`
	Logs          []byte
	ApplicationID uint
	Application   Application
}

// ApplicationArray - Placeholder for Storage Array and Application
type ApplicationArray struct {
	StorageArray
	Application
}

// StorageArray - Placeholder for storage array
type StorageArray struct {
	gorm.Model
	UniqueID           string `gorm:"unique"`
	Username           string
	Password           []byte
	ManagementEndpoint string
	StorageArrayTypeID uint
	StorageArrayType   StorageArrayType
	Applications       []Application `gorm:"many2many:storage_array_applications;"`
	MetaData           string
}

// BeforeDelete hook defined for cascade delete
func (s *StorageArray) BeforeDelete(tx *gorm.DB) error {
	if len(s.Applications) > 0 {
		return fmt.Errorf("storage array is in use by applications")
	}
	return tx.Table("application_state_change_storage_arrays").Where("storage_array_id = ?", s.ID).Unscoped().Delete(&ModuleType{}).Error
}
