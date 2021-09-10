// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package model_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/dell/csm-deployment/db"
	"github.com/dell/csm-deployment/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"

	"gorm.io/gorm"
)

func Test_BeforeDelete(t *testing.T) {
	tests := map[string]func(t *testing.T) (*gorm.DB, model.Application, model.StorageArray, model.Cluster, *gomock.Controller){
		"success - test before delete": func(*testing.T) (*gorm.DB, model.Application, model.StorageArray, model.Cluster, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			database := createTestDB("./../csm_test1.db")
			assert.NoError(t, db.AutoMigrate(database))
			application := model.Application{Name: "app1"}
			cluster := model.Cluster{ClusterName: "fake-cluster"}
			storageArray := model.StorageArray{UniqueID: "unique"}

			defer func() {
				//clean up
				assert.NoError(t, dropTestDB("./../csm_test1.db"))
			}()
			return database, application, storageArray, cluster, ctrl
		},
		"failure - test before delete": func(*testing.T) (*gorm.DB, model.Application, model.StorageArray, model.Cluster, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			database := createTestDB("./../csm_test2.db")
			assert.NoError(t, db.AutoMigrate(database))
			application := model.Application{Name: "app"}
			cluster := model.Cluster{ClusterName: "fake-cluster", Applications: []model.Application{application}}
			storageArray := model.StorageArray{UniqueID: "unique", Applications: []model.Application{application}}

			defer func() {
				//clean up
				assert.NoError(t, dropTestDB("./../csm_test2.db"))
			}()
			return database, application, storageArray, cluster, ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			database, applicationModel, storageArray, cluster, ctrl := tc(t)
			database.Create(&applicationModel)
			database.Create(&cluster)
			database.Create(&storageArray)

			database.Delete(&applicationModel)
			database.Delete(&cluster)
			database.Delete(&storageArray)
			ctrl.Finish()
		})
	}
}

func Test_BeforeDeleteApplicationFail(t *testing.T) {
	tests := map[string]func(t *testing.T) (*gorm.DB, model.Application, *gomock.Controller){
		"Failure - test before delete application should fail with an error": func(*testing.T) (*gorm.DB, model.Application, *gomock.Controller) {
			ctrl := gomock.NewController(t)

			database := createTestDB("./../csm_test3.db")
			application := model.Application{Name: "a1"}

			return database, application, ctrl
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			database, applicationModel, ctrl := tc(t)
			// we will try to delete an application from database where there is no application table itself.
			database.Delete(&applicationModel)

			defer func() {
				//clean up
				assert.NoError(t, dropTestDB("./../csm_test3.db"))
			}()
			ctrl.Finish()
		})
	}
}

// createTestDB - Generates the .db file
func createTestDB(name string) *gorm.DB {
	database, err := gorm.Open(sqlite.Open(name), &gorm.Config{})
	if err != nil {
		fmt.Println("storage err: ", err)
		return nil
	}

	sqlDB, err := database.DB()
	if err != nil {
		fmt.Println("storage err: ", err)
		return database
	}
	sqlDB.SetMaxIdleConns(3)
	return database
}

// dropTestDB - Deletes sqlite DB created
func dropTestDB(name string) error {
	if err := os.Remove(name); err != nil {
		return err
	}
	return nil
}
