package model_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/dell/csm-deployment/db"
	"github.com/dell/csm-deployment/model"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ModelTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *ModelTestSuite) SetupSuite() {
	suite.db = db.TestDB()
}

func (suite *ModelTestSuite) TearDownSuite() {
	err := db.DropTestDB()
	if err != nil {
		suite.NoError(err)
	}
}

func (suite *ModelTestSuite) TestModel() {
	db.AutoMigrate(suite.db)

	file, err := os.Open("testdata/data.yaml")
	suite.NoError(err)
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	suite.NoError(err)

	cluster := &model.Cluster{
		ClusterName:      "some cluster name",
		OrchestratorType: model.OrchestratorTypeKubernetes,
		Status:           model.ClusterStatusConnected,
		K8sVersion:       "1.21",
		ConfigFileData:   data,
	}

	// Create another cluster for multi cluster use case
	cluster1 := &model.Cluster{
		ClusterName:      "another cluster name",
		OrchestratorType: model.OrchestratorTypeKubernetes,
		Status:           model.ClusterStatusConnected,
		K8sVersion:       "1.20",
		ConfigFileData:   data,
	}

	err = suite.db.Create(cluster).Error
	suite.NoError(err)

	// Create another cluster
	err = suite.db.Create(cluster1).Error
	suite.NoError(err)

	arrays := []model.StorageArray{
		{
			UniqueID:           "unique-1",
			Username:           "user",
			Password:           "password",
			ManagementEndpoint: "127.0.0.1",
			StorageArrayTypeID: 1,
		},
		{
			UniqueID:           "unique-2",
			Username:           "user",
			Password:           "password",
			ManagementEndpoint: "127.0.0.2",
			StorageArrayTypeID: 1,
		},
	}

	application := &model.Application{
		Name:         "app1",
		Status:       model.DriverStatusOperational,
		ClusterID:    cluster.ID,
		DriverTypeID: 1,
		ModuleTypes: []model.ModuleType{
			model.ModuleType{Name: "module-1"},
			model.ModuleType{Name: "module-2"},
		},
		StorageArrays: arrays,
	}

	err = suite.db.Create(application).Error
	suite.NoError(err)

	err = suite.db.Create(&model.Application{
		Name:      "app2",
		Status:    model.DriverStatusOperational,
		ClusterID: cluster.ID,
	}).Error
	suite.NoError(err)

	// Multi cluster
	err = suite.db.Create(&model.Application{
		Name:      "app3",
		Status:    model.DriverStatusOperational,
		ClusterID: cluster.ID,
	}).Error

	suite.NoError(err)
	err = suite.db.Create(&model.Application{
		Name:      "app3",
		Status:    model.DriverStatusOperational,
		ClusterID: cluster1.ID,
	}).Error
	suite.NoError(err)

	var cl []model.Cluster
	suite.db.Preload("Config").Preload("Applications").Find(&cl)
	suite.Equal(2, len(cl))
	suite.Equal(3, len(cl[0].Applications))

	// Try to launch tasks

	applicationInstallationTask := &model.Task{
		Status:        model.TaskStatusInProgress,
		TaskType:      model.TaskTypeInstall,
		ApplicationID: application.ID,
	}
	suite.db.Create(applicationInstallationTask)

	var tasks []model.Task
	suite.db.Preload("Application").Find(&tasks)
	suite.Equal(1, len(tasks))
	suite.Equal(application.ID, tasks[0].ApplicationID)
	spew.Dump(tasks)
}

func TestModelSuite(t *testing.T) {
	suite.Run(t, new(ModelTestSuite))
}
