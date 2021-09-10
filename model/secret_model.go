// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package model

// UnityStorageArrayList - Placeholder for Unity Storage array list
type UnityStorageArrayList struct {
	StorageArrayList []UnityStorageArray `yaml:"storageArrayList"`
}

// UnityStorageArray - Placeholder for Unity storage array list
type UnityStorageArray struct {
	ArrayID                   string `yaml:"arrayId"`
	Username                  string `yaml:"username"`
	Password                  string `yaml:"password"`
	Endpoint                  string `yaml:"endpoint"`
	IsDefault                 *bool  `yaml:"isDefault"`
	SkipCertificateValidation *bool  `yaml:"skipCertificateValidation,omitempty"`
}

// IsilonClusters - Placeholder for Isilon cluster list
type IsilonClusters struct {
	IsilonCluster []IsilonCluster `yaml:"isilonClusters"`
}

// IsilonCluster - Placeholder for Isilon cluster information
type IsilonCluster struct {
	ClusterName               string `yaml:"clusterName"`
	Username                  string `yaml:"username"`
	Password                  string `yaml:"password"`
	Endpoint                  string `yaml:"endpoint"`
	EndpointPort              uint   `yaml:"endpointPort"`
	IsDefault                 *bool  `yaml:"isDefault"`
	SkipCertificateValidation *bool  `yaml:"skipCertificateValidation,omitempty"`
}

// PowerflexArray - Placeholder for Powerflex array Information
type PowerflexArray struct {
	SystemID                  string `yaml:"systemID"`
	Username                  string `yaml:"username"`
	Password                  string `yaml:"password"`
	Endpoint                  string `yaml:"endpoint"`
	IsDefault                 *bool  `yaml:"isDefault"`
	SkipCertificateValidation *bool  `yaml:"skipCertificateValidation,omitempty"`
	Mdm                       string `yaml:"mdm,omitempty"`
}

// PowerstoreSecret - Placeholder for Powerstore secret
type PowerstoreSecret struct {
	Arrays []PowerstoreArray `yaml:"arrays"`
}

// PowerstoreArray - Placeholder for Powerstore array Information
type PowerstoreArray struct {
	GlobalID                  string `yaml:"globalID"`
	Username                  string `yaml:"username"`
	Password                  string `yaml:"password"`
	Endpoint                  string `yaml:"endpoint"`
	IsDefault                 *bool  `yaml:"isDefault"`
	SkipCertificateValidation *bool  `yaml:"skipCertificateValidation,omitempty"`
	NasName                   string `yaml:"nasName"`
	BlockProtocol             string `yaml:"blockProtocol"`
}
