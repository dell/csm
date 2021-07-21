package model

import (
	"gorm.io/gorm"
)

const (
	OrchestratorTypeKubernetes string = "k8s"
	OrchestratorTypeOpenshift  string = "openshift"
	ClusterStatusConnected     string = "connected"
	ClusterStatusDisconnected  string = "disconnected"
	DriverStatusOperational    string = "operational"
	DriverStatusFailing        string = "failing"
	ModuleTypeStandalone       string = "standalone"
	ModuleTypeSidecar          string = "sidecar"
	TaskStatusInProgress       string = "in progress"
	TaskStatusCompleted        string = "completed"
	TaskStatusPrompting        string = "prompting"
	TaskStatusFailed           string = "failed"
	TaskTypeInstall            string = "install"
	TaskTypeUpdate             string = "update"
	TaskTypeDelete             string = "delete"
	ArrayTypePowerFlex         string = "powerflex"
	ArrayTypePowerMax          string = "powermax"
	ArrayTypeIsilon            string = "isilon"
	ArrayTypeUnity             string = "unity"
	ArrayTypePowerStore        string = "powerstore"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique_index;not null"`
	Password string `gorm:"not null"`
	Admin    bool
}

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

type ClusterDetails struct {
	gorm.Model // This already contains ID field
	ClusterID  string
	Nodes      string
}

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

type StorageArrayType struct {
	gorm.Model
	Name string
}

type DriverType struct {
	gorm.Model         // This already contains ID field
	Version            string
	StorageArrayTypeID uint
	StorageArrayType   StorageArrayType
}

//  ModuleType holds details about module used in application
type ModuleType struct {
	gorm.Model // This already contains ID field
	Name       string
	Version    string
	Standalone bool
}

type Task struct {
	gorm.Model
	Status        string `gorm:"not null"`
	TaskType      string `gorm:"not null"`
	Logs          []byte
	ApplicationID uint
}

type ApplicationArray struct {
	StorageArray
	Application
}

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
