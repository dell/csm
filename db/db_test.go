// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package db

import (
	"errors"
	"github.com/dell/csm-deployment/utils/constants"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dell/csm-deployment/db/mocks"
	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/store"
	storeMocks "github.com/dell/csm-deployment/store/mocks"
	"github.com/golang/mock/gomock"
)

func Test_PopulateInventory(t *testing.T) {
	tests := map[string]func(t *testing.T) (InventoryDatabaseInterface, store.StorageArrayTypeStoreInterface, *gomock.Controller){
		"success": func(*testing.T) (InventoryDatabaseInterface, store.StorageArrayTypeStoreInterface, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			db := mocks.NewMockInventoryDatabaseInterface(ctrl)
			db.EXPECT().FirstOrCreate(gomock.Any(), gomock.Any()).MinTimes(1)
			storageArrayTypeStore := storeMocks.NewMockStorageArrayTypeStoreInterface(ctrl)
			fakeStorageArrayType := model.StorageArrayType{Name: "type"}
			fakeStorageArrayType.ID = 1
			storageArrayTypeStore.EXPECT().GetByName(gomock.Any()).MinTimes(1).Return(&fakeStorageArrayType, nil)
			return db, storageArrayTypeStore, ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			db, storageArrayTypeStore, ctrl := tc(t)
			PopulateInventory(db, storageArrayTypeStore)
			ctrl.Finish()
		})
	}
}

func Test_AutoMigrate(t *testing.T) {
	tests := map[string]func(t *testing.T) (InventoryDatabaseInterface, *gomock.Controller){
		"success": func(*testing.T) (InventoryDatabaseInterface, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			db := mocks.NewMockInventoryDatabaseInterface(ctrl)
			db.EXPECT().AutoMigrate(gomock.Any()).MinTimes(1)
			return db, ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			db, ctrl := tc(t)
			assert.NoError(t, AutoMigrate(db))
			ctrl.Finish()
		})
	}
}

func Test_AutoMigrateFailure(t *testing.T) {
	tests := map[string]func(t *testing.T) (InventoryDatabaseInterface, *gomock.Controller){
		"success": func(*testing.T) (InventoryDatabaseInterface, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			db := mocks.NewMockInventoryDatabaseInterface(ctrl)
			db.EXPECT().AutoMigrate(gomock.Any()).MinTimes(1).Return(errors.New("error while auto-migrating"))
			return db, ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			db, ctrl := tc(t)
			assert.Error(t, AutoMigrate(db))
			ctrl.Finish()
		})
	}
}

func Test_addDriver(t *testing.T) {
	tests := map[string]func(t *testing.T) (InventoryDatabaseInterface, *model.DriverType, string, store.StorageArrayTypeStoreInterface, *gomock.Controller){
		"Failure - should fail with invalid driver type": func(*testing.T) (InventoryDatabaseInterface, *model.DriverType, string, store.StorageArrayTypeStoreInterface, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			db := mocks.NewMockInventoryDatabaseInterface(ctrl)
			storageArrayTypeStore := storeMocks.NewMockStorageArrayTypeStoreInterface(ctrl)
			storageArrayTypeStore.EXPECT().GetByName(gomock.Any()).MinTimes(1).Return(nil, errors.New("error looking up storage array type"))
			return db, &model.DriverType{}, "", storageArrayTypeStore, ctrl
		},
		"Failure - should fail with invalid storageArrayType": func(*testing.T) (InventoryDatabaseInterface, *model.DriverType, string, store.StorageArrayTypeStoreInterface, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			db := mocks.NewMockInventoryDatabaseInterface(ctrl)
			storageArrayTypeStore := storeMocks.NewMockStorageArrayTypeStoreInterface(ctrl)
			storageArrayTypeStore.EXPECT().GetByName(gomock.Any()).MinTimes(1).Return(nil, nil)
			return db, nil, "type", storageArrayTypeStore, ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			db, driverType, arrayName, storageArrayTypeStore, ctrl := tc(t)
			addDriver(db, driverType, arrayName, storageArrayTypeStore)
			ctrl.Finish()
		})
	}
}

func Test_New(t *testing.T) {
	tests := map[string]func(t *testing.T) (string, string, *gomock.Controller){
		"Should fail to connect to cockroach-db": func(*testing.T) (string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			return "admin", "", ctrl
		},
		"Should fail to connect to cockroach-db with bad password": func(*testing.T) (string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			return "admin", "password123", ctrl
		},
		"Should fail to connect to cockroach-db with SSL disabled password": func(*testing.T) (string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			assert.NoError(t, os.Setenv(constants.EnvDBSSLEnabled, "false"))
			return "admin", "password123", ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			db, storageArrayTypeStore, ctrl := tc(t)
			_, err := New(db, storageArrayTypeStore)
			assert.Error(t, err)
			ctrl.Finish()
		})
	}
}
