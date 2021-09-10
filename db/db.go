// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package db

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/dell/csm-deployment/store"
	"github.com/dell/csm-deployment/utils"
	"github.com/dell/csm-deployment/utils/constants"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/dell/csm-deployment/model"
)

// TestFilePath struct store paths to files for testing with sqlite unit tests
type TestFilePath struct {
	ConfigPath     string
	AuthCertPath   string
	AuthzTokenPath string
}

// SpaceChar - Placeholder for space character
const SpaceChar = " "

// New - Creates a new DB instance
func New(userName, password string) (*gorm.DB, error) {

	postgresDbStr := getDBConnectionDsn(userName, password)
	connection := postgres.Open(postgresDbStr)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
	db, err := gorm.Open(connection, &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		fmt.Println("storage err: ", err)
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(3)
	return db, nil
}

// AutoMigrate - Migrates the DB
func AutoMigrate(db InventoryDatabaseInterface) error {
	err := db.AutoMigrate(
		&model.User{},
		&model.Cluster{},
		&model.ClusterDetails{},
		&model.Application{},
		&model.Task{},
		&model.StorageArray{},
		&model.StorageArrayType{},
		&model.DriverType{},
		&model.ModuleType{},
		&model.ApplicationStateChange{},
		&model.ConfigFile{},
	)
	if err != nil {
		return err
	}
	return nil
}

func addStorageArrayType(db InventoryDatabaseInterface, storageArrayType *model.StorageArrayType) {
	db.FirstOrCreate(storageArrayType, model.StorageArrayType{Name: storageArrayType.Name})
}

func addModule(db InventoryDatabaseInterface, module *model.ModuleType) {
	db.FirstOrCreate(module, model.ModuleType{Name: module.Name, Version: module.Version})
}

func addDriver(db InventoryDatabaseInterface, driver *model.DriverType, storageArrayTypeName string, storageArrayTypeStore store.StorageArrayTypeStoreInterface) {
	storageArrayType, err := storageArrayTypeStore.GetByName(storageArrayTypeName)
	if err != nil {
		log.Printf("error looking up storage array type for %s", storageArrayTypeName)
		return
	}

	if storageArrayType == nil {
		log.Printf("unable to find storage array type for %s", storageArrayTypeName)
		return
	}

	driver.StorageArrayTypeID = storageArrayType.ID

	db.FirstOrCreate(driver, model.DriverType{Version: driver.Version, StorageArrayTypeID: driver.StorageArrayTypeID})
}

// InventoryDatabaseInterface provides an interface to perform db operations during inventory population
//go:generate mockgen -destination=mocks/inventory_database_inventory.go -package=mocks github.com/dell/csm-deployment/db InventoryDatabaseInterface
type InventoryDatabaseInterface interface {
	FirstOrCreate(interface{}, ...interface{}) (tx *gorm.DB)
	AutoMigrate(...interface{}) error
}

// PopulateInventory - Adds minimum data to the database
func PopulateInventory(db InventoryDatabaseInterface, storageArrayTypeStore store.StorageArrayTypeStoreInterface) {
	// add storage types
	addStorageArrayType(db, &model.StorageArrayType{Name: model.ArrayTypePowerFlex})
	addStorageArrayType(db, &model.StorageArrayType{Name: model.ArrayTypePowerMax})
	addStorageArrayType(db, &model.StorageArrayType{Name: model.ArrayTypePowerScale})
	addStorageArrayType(db, &model.StorageArrayType{Name: model.ArrayTypePowerStore})
	addStorageArrayType(db, &model.StorageArrayType{Name: model.ArrayTypeUnity})

	// add modules
	addModule(db, &model.ModuleType{
		Name:    model.ModuleTypePodMon,
		Version: "1.0.0",
	})
	addModule(db, &model.ModuleType{
		Name:       model.ModuleTypeObservability,
		Version:    "1.0.0",
		Standalone: true,
	})
	addModule(db, &model.ModuleType{
		Name:    model.ModuleTypeVgSnapShotter,
		Version: "1.0.0",
	})
	addModule(db, &model.ModuleType{
		Name:    model.ModuleTypeReplication,
		Version: "1.0.0",
	})
	addModule(db, &model.ModuleType{
		Name:       model.ModuleTypeAuthorization,
		Version:    "1.0.0",
		Standalone: false,
	})
	addModule(db, &model.ModuleType{
		Name:       model.ModuleTypeReverseProxy,
		Version:    "2.0.0",
		Standalone: false,
	})
	// add drivers
	addDriver(db, &model.DriverType{
		Version: "2.0.0",
	}, model.ArrayTypePowerFlex, storageArrayTypeStore)

	addDriver(db, &model.DriverType{
		Version: "2.0.0",
	}, model.ArrayTypePowerScale, storageArrayTypeStore)

	addDriver(db, &model.DriverType{
		Version: "2.0.0",
	}, model.ArrayTypeUnity, storageArrayTypeStore)

	addDriver(db, &model.DriverType{
		Version: "2.0.0",
	}, model.ArrayTypePowerStore, storageArrayTypeStore)

	addDriver(db, &model.DriverType{
		Version: "2.0.0",
	}, model.ArrayTypePowerMax, storageArrayTypeStore)
}

// getDBConnectionDsn - Populate DB Connection String from Environment variables
func getDBConnectionDsn(dbUserName, dbPassword string) string {

	var buffer bytes.Buffer

	buffer.WriteString("host=")
	buffer.WriteString(utils.GetEnv(constants.EnvDBHost, "cockroachdb-public"))
	buffer.WriteString(SpaceChar)

	buffer.WriteString("port=")
	buffer.WriteString(utils.GetEnv(constants.EnvDBPort, "26257"))
	buffer.WriteString(SpaceChar)

	buffer.WriteString("user=")
	buffer.WriteString(dbUserName)
	buffer.WriteString(SpaceChar)

	if dbPassword != "" {
		buffer.WriteString("password=")
		buffer.WriteString(dbPassword)
		buffer.WriteString(SpaceChar)
	}

	sslEnabled, err := strconv.ParseBool(utils.GetEnv(constants.EnvDBSSLEnabled, "true"))
	if err != nil {
		sslEnabled = false
	}

	if !sslEnabled {
		buffer.WriteString("sslmode=disable")
		buffer.WriteString(SpaceChar)
	} else {

		dbCertDir := "/app/dbclient-certificates"
		dbCertFileName := "tls.crt"
		dbKeyFileName := "tls.key"

		buffer.WriteString("sslmode=require")
		buffer.WriteString(SpaceChar)

		buffer.WriteString("sslcert=")
		buffer.WriteString(path.Join(dbCertDir, dbCertFileName))
		buffer.WriteString(SpaceChar)

		buffer.WriteString("sslkey=")
		buffer.WriteString(path.Join(dbCertDir, dbKeyFileName))
	}

	return buffer.String()
}
