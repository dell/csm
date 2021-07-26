package model

import (
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

	// ArrayTypeIsilon - placeholder for constant isilon
	ArrayTypeIsilon string = "isilon"

	// ArrayTypeUnity - placeholder for constant unity
	ArrayTypeUnity string = "unity"

	// ArrayTypePowerStore - placeholder for constant powerstore
	ArrayTypePowerStore string = "powerstore"
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

// ClusterDetails - Placeholder for Cluster and Nodes information
type ClusterDetails struct {
	gorm.Model // This already contains ID field
	ClusterID  string
	Nodes      string
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
	Password           string
	ManagementEndpoint string
	StorageArrayTypeID uint
	StorageArrayType   StorageArrayType
	Applications       []Application `gorm:"many2many:storage_array_applications;"`
}
