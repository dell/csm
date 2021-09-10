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

// ApplicationStoreInterface is used to define the interface for persisting Application
//go:generate mockgen -destination=mocks/application_store_interface.go -package=mocks github.com/dell/csm-deployment/store ApplicationStoreInterface
type ApplicationStoreInterface interface {
	Create(*model.Application) error
	GetByID(string) (*model.Application, error)
	GetByName(string) (*model.Application, error)
	Update(*model.Application) error
	Delete(a *model.Application) error
	GetAll() ([]model.Application, error)
}

// ApplicationStore - Placeholder for Application
type ApplicationStore struct {
	db *gorm.DB
}

// NewApplicationStore - returns an instance of ApplicationStore in db
func NewApplicationStore(db *gorm.DB) *ApplicationStore {
	return &ApplicationStore{
		db: db,
	}
}

// GetAll returns all applications
func (as *ApplicationStore) GetAll() ([]model.Application, error) {
	var applications []model.Application
	if err := as.db.Preload(clause.Associations).Find(&applications).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return applications, nil
}

// GetByID - returns Application by Id
func (as *ApplicationStore) GetByID(id string) (*model.Application, error) {
	var m model.Application
	if err := as.db.Preload(clause.Associations).Preload("StorageArrays.StorageArrayType").First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

// GetByName - returns Application by Name
func (as *ApplicationStore) GetByName(name string) (*model.Application, error) {
	var application model.Application
	if err := as.db.
		Preload(clause.Associations).
		First(&application, &model.Application{Name: name}).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &application, nil
}

// Create new Application
func (as *ApplicationStore) Create(a *model.Application) (err error) {
	return as.db.Save(a).Error
}

// Update Application info
func (as *ApplicationStore) Update(a *model.Application) (err error) {
	return as.db.Save(a).Error
}

// Delete Application
func (as *ApplicationStore) Delete(a *model.Application) error {
	return as.db.Unscoped().Delete(a).Error
}
