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

	"github.com/dell/csm-deployment/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DriverTypeStoreInterface is used to define the interface for persisting Driver type
//go:generate mockgen -destination=mocks/driver_type_store_interface.go -package=mocks github.com/dell/csm-deployment/store DriverTypeStoreInterface
type DriverTypeStoreInterface interface {
	GetByID(uint) (*model.DriverType, error)
	GetAll() ([]model.DriverType, error)
}

// DriverTypeStore - Placeholder for Driver type store
type DriverTypeStore struct {
	db *gorm.DB
}

// NewDriverTypeStore - returns an instance of DriverTypeStore in db
func NewDriverTypeStore(db *gorm.DB) *DriverTypeStore {
	return &DriverTypeStore{
		db: db,
	}
}

// GetByID - Method to get Driver Type by Id
func (sas *DriverTypeStore) GetByID(id uint) (*model.DriverType, error) {
	var sa model.DriverType
	if err := sas.db.Preload(clause.Associations).First(&sa, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sa, nil
}

// GetAll - Method to get all driver types
func (sas *DriverTypeStore) GetAll() ([]model.DriverType, error) {
	var sa []model.DriverType
	if err := sas.db.Preload(clause.Associations).Find(&sa).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return sa, nil
}
