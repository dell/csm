package ytt_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/dell/csm-deployment/db"
	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/store"
	"github.com/dell/csm-deployment/ytt"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type YttTestSuite struct {
	suite.Suite
	inputPaths []string
	outputPath string
	db         *gorm.DB
}

func (suite *YttTestSuite) SetupSuite() {
	suite.inputPaths = append(suite.inputPaths, "testdata/")
	suite.outputPath = "testdata/output/"
	suite.db = db.TestDB()
	db.AutoMigrate(suite.db)
	db.PopulateTestDb(suite.db, "testdata/kubeconfig.yaml")
}

func (suite *YttTestSuite) TearDownSuite() {
	_ = os.RemoveAll(suite.outputPath)

	err := db.DropTestDB()
	if err != nil {
		suite.NoError(err)
	}
}

func (suite *YttTestSuite) TestTemplate() {
	client := ytt.NewClient()
	output, err := client.Template(suite.inputPaths, []string{"port=9090", "text=\"hey, people\""})
	suite.NoError(err)

	err = output.CreateAt(suite.outputPath)
	suite.NoError(err)
	suite.FileExists(filepath.Join(suite.outputPath, "config.yaml"))

	data := output.AsBytes()
	suite.Equal(2, len(data))
	spew.Dump(data)
}

func (suite *YttTestSuite) TestTemplateFromApplication() {
	applications := store.NewApplicationStateChangeStore(suite.db)
	clusters := store.NewClusterStore(suite.db)

	client := ytt.NewClient(ytt.WithTemplatePath("../"))
	output, err := client.TemplateFromApplication(1, applications, clusters)
	suite.NoError(err)

	err = output.CreateAt(suite.outputPath)
	suite.NoError(err)
	suite.FileExists(filepath.Join(suite.outputPath, "node.yaml"))
	suite.FileExists(filepath.Join(suite.outputPath, "csidriver.yaml"))
	suite.FileExists(filepath.Join(suite.outputPath, "controller.yaml"))

	data := output.AsCombinedBytes()
	err = ioutil.WriteFile(filepath.Join(suite.outputPath, "output.yaml"), data, 0644)
	suite.NoError(err)
}

func (suite *YttTestSuite) TestTemplateFromApplication_StandaloneModule() {
	applications := store.NewApplicationStateChangeStore(suite.db)
	clusters := store.NewClusterStore(suite.db)

	observabilityModule := &model.ModuleType{
		Name:       "observability",
		Standalone: true,
	}
	suite.db.Where(&model.ModuleType{Name: observabilityModule.Name}).FirstOrCreate(observabilityModule)

	array := &model.StorageArray{}
	suite.db.First(array)

	app := &model.Application{
		Name:          "standalone-module-app",
		Status:        model.DriverStatusOperational,
		ClusterID:     1,
		ModuleTypes:   []model.ModuleType{*observabilityModule},
		StorageArrays: []model.StorageArray{*array},
	}
	suite.db.Create(app)

	stateChange := &model.ApplicationStateChange{
		ApplicationID:       app.ID,
		ClusterID:           1,
		ModuleTypes:         []model.ModuleType{*observabilityModule},
		StorageArrays:       []model.StorageArray{*array},
		ModuleConfiguration: "observability.namespace=test-csm-namespace",
	}
	suite.db.Create(stateChange)

	client := ytt.NewClient(ytt.WithTemplatePath("../"))
	output, err := client.TemplateFromApplication(stateChange.ID, applications, clusters)
	suite.NoError(err)

	err = output.CreateAt(suite.outputPath)
	suite.NoError(err)
	suite.FileExists(filepath.Join(suite.outputPath, "observability-0.3.0.yaml"))

	data := output.AsCombinedBytes()
	err = ioutil.WriteFile(filepath.Join(suite.outputPath, "output.yaml"), data, 0644)
	suite.NoError(err)
}

func TestYttSuite(t *testing.T) {
	suite.Run(t, new(YttTestSuite))
}
