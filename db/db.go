package db

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"gorm.io/gorm"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm/logger"

	"path"

	"github.com/dell/csm-deployment/model"
)

func New(address string) (*gorm.DB, error) {
	dbFile := ""
	if address == "" {
		dbFile = "./csm.db"
	} else {
		dbFile = path.Join(address, "csm.db")
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{
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

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
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
	)
}

func PopulateInventory(db *gorm.DB) {
	powerflex := &model.StorageArrayType{Name: model.ArrayTypePowerFlex}
	db.Create(powerflex)

	isilon := &model.StorageArrayType{Name: model.ArrayTypeIsilon}
	db.Create(isilon)

	podmon := &model.ModuleType{
		Name:    "podmon",
		Version: "0.0.1",
	}
	db.Create(podmon)

	observabilityModule := &model.ModuleType{
		Name:       "observability",
		Standalone: true,
	}
	db.Create(observabilityModule)

	vgsnapshotter := &model.ModuleType{
		Name:    "vgsnapshotter",
		Version: "0.0.1",
	}
	db.Create(vgsnapshotter)

	powerflexdriver14 := &model.DriverType{
		Version:            "1.4.0",
		StorageArrayTypeID: powerflex.ID,
	}
	db.Create(powerflexdriver14)

	isilondriver15 := &model.DriverType{
		Version:            "1.5.0",
		StorageArrayTypeID: isilon.ID,
	}
	db.Create(isilondriver15)
}

func PopulateTestDb(db *gorm.DB, configPath string) {
	user := &model.User{
		Username: "admin",
		Password: "admin",
		Admin:    true,
	}
	h, err := user.HashPassword(user.Password)
	if err != nil {
		panic(err)
	}
	user.Password = h
	db.Create(user)

	arrayType := &model.StorageArrayType{Name: model.ArrayTypePowerFlex}
	db.Create(arrayType)

	array := &model.StorageArray{
		UniqueID:           "id-1",
		Username:           "user",
		Password:           "password",
		ManagementEndpoint: "127.0.01",
		StorageArrayTypeID: arrayType.ID,
	}
	db.Create(array)

	podmonModule := &model.ModuleType{
		Name:       "podmon",
		Standalone: false,
	}
	db.Create(podmonModule)

	observabilityModule := &model.ModuleType{
		Name:       "observability",
		Standalone: true,
	}
	db.Create(observabilityModule)

	kubeConfigData, err := ioutil.ReadFile(configPath)

	cluster := &model.Cluster{
		ClusterName:      "test-cluster",
		ConfigFileData:   kubeConfigData,
		OrchestratorType: "k8s",
		Status:           model.ClusterStatusConnected,
		K8sVersion:       "1.20",
	}
	if err := db.Create(cluster).Error; err != nil {
		panic(err)
	}

	app := &model.Application{
		Name:          "test-app",
		Status:        model.DriverStatusOperational,
		ClusterID:     cluster.ID,
		DriverTypeID:  1,
		ModuleTypes:   []model.ModuleType{*podmonModule},
		StorageArrays: []model.StorageArray{*array},
	}
	db.Create(app)

	stateChange := &model.ApplicationStateChange{
		ApplicationID: app.ID,
		ClusterID:     cluster.ID,
		DriverTypeID:  1,
		ModuleTypes:   []model.ModuleType{*podmonModule},
		StorageArrays: []model.StorageArray{*array},
	}
	db.Create(stateChange)
}

func TestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("./../csm_test.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("storage err: ", err)
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("storage err: ", err)
		return db
	}
	sqlDB.SetMaxIdleConns(3)
	return db
}

func DropTestDB() error {
	if err := os.Remove("./../csm_test.db"); err != nil {
		return err
	}
	return nil
}
