// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package store

import (
	"errors"

	"gorm.io/gorm/clause"

	"github.com/dell/csm-deployment/model"
	"gorm.io/gorm"
)

// ConfigFileStoreInterface is used to define the interface for persisting configuration files
//go:generate mockgen -destination=mocks/config_file_store_interface.go -package=mocks github.com/dell/csm-deployment/store ConfigFileStoreInterface
type ConfigFileStoreInterface interface {
	Create(*model.ConfigFile) error
	GetByID(clusterID uint) (*model.ConfigFile, error)
	Delete(u *model.ConfigFile) error
	Update(u *model.ConfigFile) error
	GetAll() ([]model.ConfigFile, error)
	GetAllByName(string) ([]model.ConfigFile, error)
}

// ConfigFileStore is used to operate on the configuration files persistent store
type ConfigFileStore struct {
	db *gorm.DB
}

// NewConfigFileStore creates a new configuration file
func NewConfigFileStore(db *gorm.DB) *ConfigFileStore {
	return &ConfigFileStore{
		db: db,
	}
}

// GetAllByName will return all configuration files with a matching name
func (us *ConfigFileStore) GetAllByName(name string) ([]model.ConfigFile, error) {
	var cfs []model.ConfigFile
	if err := us.db.
		Preload(clause.Associations).
		Where(&model.ConfigFile{Name: name}).
		Find(&cfs).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return cfs, nil
}

// GetByID returns a configuration file with the given ID
func (us *ConfigFileStore) GetByID(cfID uint) (*model.ConfigFile, error) {
	var m model.ConfigFile
	if err := us.db.Preload(clause.Associations).First(&m, cfID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

// GetAll returns all configuration files
func (us *ConfigFileStore) GetAll() ([]model.ConfigFile, error) {
	var cfs []model.ConfigFile
	if err := us.db.Preload(clause.Associations).Find(&cfs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return cfs, nil
}

// Create will persist a new configuration file in the database
func (us *ConfigFileStore) Create(u *model.ConfigFile) (err error) {
	return us.db.Create(u).Error
}

// Update will update an existing configuration file in the database
func (us *ConfigFileStore) Update(u *model.ConfigFile) error {
	return us.db.Save(u).Error
}

// Delete will delete an existing configuration file from the database
func (us *ConfigFileStore) Delete(u *model.ConfigFile) error {
	return us.db.Unscoped().Delete(u).Error
}
