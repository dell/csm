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

// StorageArrayTypeStoreInterface is used to define the interface for persisting Storage Array Types
//go:generate mockgen -destination=mocks/storage_array_type_store_interface.go -package=mocks github.com/dell/csm-deployment/store StorageArrayTypeStoreInterface
type StorageArrayTypeStoreInterface interface {
	GetByID(uint) (*model.StorageArrayType, error)
	GetAll() ([]model.StorageArrayType, error)
	GetAllByID(...uint) ([]model.StorageArrayType, error)
	GetByName(name string) (*model.StorageArrayType, error)
}

// StorageArrayTypeStore - Placeholder for Storage Array Type Store
type StorageArrayTypeStore struct {
	db *gorm.DB
}

// NewStorageArrayTypeStore returns an instance of StorageArrayTypeStore in db
func NewStorageArrayTypeStore(db *gorm.DB) *StorageArrayTypeStore {
	return &StorageArrayTypeStore{
		db: db,
	}
}

// GetByID returns an instance of StorageArrayTypeStore  that matches id in db
func (ms *StorageArrayTypeStore) GetByID(id uint) (*model.StorageArrayType, error) {
	var mt model.StorageArrayType
	if err := ms.db.Preload(clause.Associations).First(&mt, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &mt, nil
}

// GetAll returns all instances of StorageArrayTypeStore in db
func (ms *StorageArrayTypeStore) GetAll() ([]model.StorageArrayType, error) {
	var sa []model.StorageArrayType
	if err := ms.db.Preload(clause.Associations).Find(&sa).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return sa, nil
}

// GetAllByID returns instances of StorageArrayTypeStore in db that match all passed in IDs
func (ms *StorageArrayTypeStore) GetAllByID(v ...uint) ([]model.StorageArrayType, error) {
	var mt []model.StorageArrayType
	if len(v) > 0 {
		if err := ms.db.Preload(clause.Associations).Find(&mt, v).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, nil
			}
			return nil, err
		}
	}
	return mt, nil
}

// GetByName will return the storage type by name
func (ms *StorageArrayTypeStore) GetByName(name string) (*model.StorageArrayType, error) {
	var sat model.StorageArrayType
	if err := ms.db.
		Preload(clause.Associations).
		First(&sat, model.StorageArrayType{Name: name}).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sat, nil
}
