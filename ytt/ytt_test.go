// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package ytt_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gorm.io/driver/sqlite"

	"github.com/davecgh/go-spew/spew"
	"github.com/dell/csm-deployment/db"
	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/store"
	"github.com/dell/csm-deployment/store/mocks"
	"github.com/dell/csm-deployment/utils"
	"github.com/dell/csm-deployment/ytt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/golang/mock/gomock"
)

type YttTestSuite struct {
	suite.Suite
	inputPaths []string
	outputPath string
	db         *gorm.DB
}

type testOverrides struct {
	decryptPassword      func(cipherText []byte) ([]byte, error)
	skipArrayCreation    bool
	errorGetByIDAppState bool
	nilGetByIDAppState   bool
	errorGetByIDCluster  bool
}

type testFilePath struct {
	configPath     string
	authCertPath   string
	authzTokenPath string
}

func RunTestSuite(t *testing.T, fn func(appID uint, suite *YttTestSuite, app *model.Application, stateChange *model.ApplicationStateChange) (ytt.Output, error),
	outputFiles []string, success, skipArrayCreation bool, array *model.StorageArray, app *model.Application, stateChange *model.ApplicationStateChange) {

	suite := new(YttTestSuite)
	suite.inputPaths = append(suite.inputPaths, "testdata/")
	suite.outputPath = "testdata/output/"
	suite.db = createTestDB("./../csm_test.db")
	err := db.AutoMigrate(suite.db)
	if err != nil {
		log.Fatal(err)
	}

	utils.CipherKey = []byte("thisisa32bytecharactercipherkey!")

	defer func() {
		// clean up
		_ = os.RemoveAll(suite.outputPath)
		assert.NoError(t, dropTestDB("./../csm_test.db"))
	}()

	// add modules to db
	for _, m := range stateChange.ModuleTypes {
		t := m
		suite.db.Where(&model.ModuleType{Name: t.Name}).FirstOrCreate(&t)
	}

	// add arrays
	if !skipArrayCreation {
		arrayType := &model.StorageArrayType{Name: array.StorageArrayType.Name}
		suite.db.Create(arrayType)
		suite.db.Create(array)

		// add app
		app.StorageArrays = []model.StorageArray{*array}
		stateChange.StorageArrays = []model.StorageArray{*array}
	}

	output, err := fn(1, suite, app, stateChange)

	if !success {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
		data := output.AsCombinedBytes()

		err = output.CreateAt(suite.outputPath)
		assert.NoError(t, err)

		spew.Dump(data)

		for _, f := range outputFiles {
			assert.FileExists(t, filepath.Join(suite.outputPath, f))
		}

	}
}

func DriverComponents(driverType string) (*model.StorageArray, *model.Application, *model.ApplicationStateChange) {
	arrayType := &model.StorageArrayType{Name: driverType}

	array := &model.StorageArray{
		StorageArrayType:   *arrayType,
		StorageArrayTypeID: arrayType.ID,
		ManagementEndpoint: "testing-endpoint",
		Username:           "testing-username",
		UniqueID:           "01234569",
	}

	app := &model.Application{
		Name:                driverType,
		Status:              model.DriverStatusOperational,
		ClusterID:           0,
		DriverTypeID:        1,
		ModuleTypes:         []model.ModuleType{},
		DriverConfiguration: fmt.Sprintf("namespace=%s", driverType),
	}

	stateChange := &model.ApplicationStateChange{
		ApplicationID:       app.ID,
		ClusterID:           0,
		DriverTypeID:        1,
		ModuleTypes:         []model.ModuleType{},
		ModuleConfiguration: "",
		DriverConfiguration: fmt.Sprintf("namespace=%s", driverType),
	}
	return array, app, stateChange
}

func Test_TemplateFromApplication(t *testing.T) {
	driverNoError := func(driverType string) (bool, testOverrides, []string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
		array, app, stateChange := DriverComponents(driverType)
		return true, testOverrides{}, []string{"node.yaml", "csidriver.yaml", "controller.yaml"}, []string{}, array, app, stateChange
	}
	tests := map[string]func(t *testing.T) (bool, testOverrides, []string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange){
		// successful cases
		"success observability standalone": func(*testing.T) (bool, testOverrides, []string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			observabilityModule := &model.ModuleType{
				Name:       "observability",
				Standalone: true,
			}

			array := &model.StorageArray{}

			app := &model.Application{
				Name:                "standalone-module-app",
				Status:              model.DriverStatusOperational,
				ClusterID:           0,
				ModuleTypes:         []model.ModuleType{*observabilityModule},
				ModuleConfiguration: "observability.namespace=test-observability-auth-namespace",
			}

			stateChange := &model.ApplicationStateChange{
				ApplicationID:       app.ID,
				ClusterID:           0,
				ModuleTypes:         []model.ModuleType{*observabilityModule},
				ModuleConfiguration: "observability.namespace=test-observability-auth-namespace",
			}

			return true, testOverrides{}, []string{"observability.yaml"}, []string{}, array, app, stateChange
		},
		"success driver-powermax with reverse-proxy": func(*testing.T) (bool, testOverrides, []string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			array, app, stateChange := DriverComponents(model.ArrayTypePowerMax)
			reverseproxy := &model.ModuleType{
				Name:       model.ModuleTypeReverseProxy,
				Standalone: true,
			}
			app.ModuleTypes = []model.ModuleType{*reverseproxy}
			stateChange.ModuleTypes = []model.ModuleType{*reverseproxy}
			array.MetaData = "portGroups=iscsi_csm_cicd"

			return true, testOverrides{}, []string{"node.yaml", "csidriver.yaml", "controller.yaml", "csireverseproxy.yaml"}, []string{}, array, app, stateChange
		},
		"success driver-powermax with reverse-proxy and authorization": func(*testing.T) (bool, testOverrides, []string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			array, app, stateChange := DriverComponents(model.ArrayTypePowerMax)
			reverseproxy := &model.ModuleType{
				Name:       model.ModuleTypeReverseProxy,
				Standalone: true,
			}
			authorizationModule := &model.ModuleType{
				Name:       model.ModuleTypeAuthorization,
				Version:    "0.2.1",
				Standalone: false,
			}

			app.ModuleTypes = []model.ModuleType{*reverseproxy, *authorizationModule}
			stateChange.ModuleTypes = []model.ModuleType{*reverseproxy, *authorizationModule}

			app.ModuleConfiguration = "csireverseproxy.deployAsSidecar=false"
			stateChange.ModuleConfiguration = "csireverseproxy.deployAsSidecar=false"
			array.MetaData = "portGroups=iscsi_csm_cicd"

			return true, testOverrides{}, []string{"node.yaml", "csidriver.yaml", "controller.yaml", "csireverseproxy.yaml"}, []string{}, array, app, stateChange
		},
		"success driver-powerflex only": func(*testing.T) (bool, testOverrides, []string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			return driverNoError(model.ArrayTypePowerFlex)
		},
		"success driver-powerscale only": func(*testing.T) (bool, testOverrides, []string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			return driverNoError(model.ArrayTypePowerScale)
		},
		"sucess driver-powerstore only": func(*testing.T) (bool, testOverrides, []string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			return driverNoError(model.ArrayTypePowerStore)
		},
		"success driver-powermax only": func(*testing.T) (bool, testOverrides, []string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			return driverNoError(model.ArrayTypePowerMax)
		},
		"success driver-unity only": func(*testing.T) (bool, testOverrides, []string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			return driverNoError(model.ArrayTypeUnity)
		},

		// failure
		"fail observability standalone": func(*testing.T) (bool, testOverrides, []string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			observabilityModule := &model.ModuleType{
				Name:       "observability",
				Standalone: true,
			}

			array := &model.StorageArray{}

			app := &model.Application{
				Name:                "standalone-module-app",
				Status:              model.DriverStatusOperational,
				ClusterID:           0,
				ModuleTypes:         []model.ModuleType{*observabilityModule},
				ModuleConfiguration: "observability.namespace=fail=test-observability-auth-namespace",
			}

			stateChange := &model.ApplicationStateChange{
				ApplicationID:       app.ID,
				ClusterID:           0,
				ModuleTypes:         []model.ModuleType{*observabilityModule},
				ModuleConfiguration: "observability.namespace=fail=test-observability-auth-namespace",
			}

			return false, testOverrides{}, []string{}, []string{}, array, app, stateChange
		},
		"fail podmon as standalone": func(*testing.T) (bool, testOverrides, []string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			podManModule := &model.ModuleType{
				Name: model.ModuleTypePodMon,
			}

			array := &model.StorageArray{}

			app := &model.Application{
				Name:        "standalone-module-app",
				Status:      model.DriverStatusOperational,
				ClusterID:   0,
				ModuleTypes: []model.ModuleType{*podManModule},
			}

			stateChange := &model.ApplicationStateChange{
				ApplicationID: app.ID,
				ClusterID:     0,
				ModuleTypes:   []model.ModuleType{*podManModule},
			}

			return false, testOverrides{}, []string{}, []string{}, array, app, stateChange
		},
		"fail couldn't find storage arrays for app state with": func(*testing.T) (bool, testOverrides, []string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			array, app, stateChange := DriverComponents(model.ArrayTypePowerFlex)
			patch := testOverrides{
				skipArrayCreation: true,
			}
			return false, patch, []string{}, []string{}, array, app, stateChange
		},
		"error querying app state db": func(*testing.T) (bool, testOverrides, []string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			array, app, stateChange := DriverComponents(model.ArrayTypePowerFlex)
			patch := testOverrides{
				errorGetByIDAppState: true,
			}
			return false, patch, []string{}, []string{}, array, app, stateChange
		},
		"error querying cluster db": func(*testing.T) (bool, testOverrides, []string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			array, app, stateChange := DriverComponents(model.ArrayTypePowerFlex)
			patch := testOverrides{
				errorGetByIDCluster: true,
			}
			return false, patch, []string{}, []string{}, array, app, stateChange
		},
		"nil result from db": func(*testing.T) (bool, testOverrides, []string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			array, app, stateChange := DriverComponents(model.ArrayTypePowerFlex)
			patch := testOverrides{
				nilGetByIDAppState: true,
			}
			return false, patch, []string{}, []string{}, array, app, stateChange
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			success, patch, outputFiles, configFilePaths, array, app, stateChange := tc(t)

			if patch.decryptPassword != nil {
				oldDecryptPassword := utils.DecryptPassword
				defer func() { utils.DecryptPassword = oldDecryptPassword }()
				utils.DecryptPassword = patch.decryptPassword
			}

			templateFromApplicationWrapper := func(appID uint, suite *YttTestSuite,
				app *model.Application, stateChange *model.ApplicationStateChange) (ytt.Output, error) {
				client := ytt.NewClient()
				client.SetOptions(ytt.WithTemplatePath("../"))

				cf := store.NewConfigFileStore(suite.db)
				clusters := store.NewClusterStore(suite.db)
				application := store.NewApplicationStateChangeStore(suite.db)

				// add cluster
				cluster := &model.Cluster{
					ClusterName:      "test-cluster",
					OrchestratorType: "k8s",
					Status:           model.ClusterStatusConnected,
					K8sVersion:       "1.22",
					ConfigFileData:   []byte(""),
				}
				suite.db.Create(cluster)

				// add configuration files
				for _, path := range configFilePaths {
					cfData, err := ioutil.ReadFile(path)
					if err != nil {
						panic(err)
					}
					cfInstance := &model.ConfigFile{Name: path, ConfigFileData: cfData}
					suite.db.Create(cfInstance)
				}

				app.ClusterID = cluster.ID
				stateChange.ClusterID = cluster.ID

				suite.db.Create(app)
				suite.db.Create(stateChange)

				if patch.errorGetByIDAppState {
					ctrl := gomock.NewController(t)

					as := mocks.NewMockApplicationStateChangeStoreInterface(ctrl)
					as.EXPECT().GetByID(appID).Times(1).Return(nil, errors.New("error"))
					defer func() { ctrl.Finish() }()

					return client.TemplateFromApplication(appID, as, clusters, cf)

				} else if patch.errorGetByIDCluster {
					ctrl := gomock.NewController(t)

					cl := mocks.NewMockClusterStoreInterface(ctrl)
					cl.EXPECT().GetByID(appID).Times(1).Return(nil, errors.New("error"))
					defer func() { ctrl.Finish() }()

					return client.TemplateFromApplication(appID, application, cl, cf)

				} else if patch.nilGetByIDAppState {
					ctrl := gomock.NewController(t)

					as := mocks.NewMockApplicationStateChangeStoreInterface(ctrl)
					as.EXPECT().GetByID(appID).Times(1).Return(nil, nil)
					defer func() { ctrl.Finish() }()

					return client.TemplateFromApplication(appID, as, clusters, cf)

				}

				return client.TemplateFromApplication(1, application, clusters, cf)

			}

			RunTestSuite(t, templateFromApplicationWrapper, outputFiles, success, patch.skipArrayCreation, array, app, stateChange)
		})
	}
}
func Test_NamespaceTemplateFromApplication(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, testOverrides, *model.StorageArray, *model.Application, *model.ApplicationStateChange){
		"success observability standalone": func(*testing.T) (bool, testOverrides, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {

			observabilityModule := &model.ModuleType{
				Name:       "observability",
				Standalone: true,
			}

			array := &model.StorageArray{}

			app := &model.Application{
				Name:        "standalone-module-app",
				Status:      model.DriverStatusOperational,
				ClusterID:   0,
				ModuleTypes: []model.ModuleType{*observabilityModule},
			}

			stateChange := &model.ApplicationStateChange{
				ApplicationID:       app.ID,
				ClusterID:           0,
				ModuleTypes:         []model.ModuleType{*observabilityModule},
				ModuleConfiguration: "",
			}

			return true, testOverrides{}, array, app, stateChange
		},
		"success driver": func(*testing.T) (bool, testOverrides, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			array, app, stateChange := DriverComponents(model.ArrayTypePowerFlex)
			return true, testOverrides{}, array, app, stateChange
		},
		"fail driver no storage": func(*testing.T) (bool, testOverrides, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			array, app, stateChange := DriverComponents(model.ArrayTypePowerFlex)
			patch := testOverrides{
				skipArrayCreation: true,
			}
			return false, patch, array, app, stateChange
		},
		"error querying app state db": func(*testing.T) (bool, testOverrides, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			array, app, stateChange := DriverComponents(model.ArrayTypePowerFlex)
			patch := testOverrides{
				errorGetByIDAppState: true,
			}
			return false, patch, array, app, stateChange
		},
		"nil result from db": func(*testing.T) (bool, testOverrides, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			array, app, stateChange := DriverComponents(model.ArrayTypePowerFlex)
			patch := testOverrides{
				nilGetByIDAppState: true,
			}
			return false, patch, array, app, stateChange
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			success, patch, array, app, stateChange := tc(t)

			namespaceTemplateFromApplicationWrapper := func(appID uint, suite *YttTestSuite, app *model.Application, stateChange *model.ApplicationStateChange) (ytt.Output, error) {
				client := ytt.NewClient()
				client.SetOptions(ytt.WithTemplatePath("../"))

				if patch.errorGetByIDAppState {
					ctrl := gomock.NewController(t)

					as := mocks.NewMockApplicationStateChangeStoreInterface(ctrl)
					as.EXPECT().GetByID(appID).Times(1).Return(nil, errors.New("error"))
					defer func() { ctrl.Finish() }()

					return client.NamespaceTemplateFromApplication(appID, as)

				} else if patch.nilGetByIDAppState {
					ctrl := gomock.NewController(t)

					as := mocks.NewMockApplicationStateChangeStoreInterface(ctrl)
					as.EXPECT().GetByID(appID).Times(1).Return(nil, nil)
					defer func() { ctrl.Finish() }()

					return client.NamespaceTemplateFromApplication(appID, as)

				}

				application := store.NewApplicationStateChangeStore(suite.db)

				suite.db.Create(app)
				suite.db.Create(stateChange)

				return client.NamespaceTemplateFromApplication(appID, application)

			}

			RunTestSuite(t, namespaceTemplateFromApplicationWrapper, []string{"namespace.yaml"}, success, patch.skipArrayCreation, array, app, stateChange)
		})
	}
}
func Test_GenerateDynamicSecret(t *testing.T) {
	driverNoError := func(driverType, outputFile string) (bool, testOverrides, string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
		array, app, stateChange := DriverComponents(driverType)
		return true, testOverrides{}, outputFile, []string{}, array, app, stateChange
	}

	driverError := func(driverType, outputFile string) (bool, testOverrides, string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
		patch := testOverrides{
			decryptPassword: func(cipherText []byte) ([]byte, error) {
				return nil, errors.New("error")
			},
		}
		array, app, stateChange := DriverComponents(driverType)

		return false, patch, outputFile, []string{}, array, app, stateChange
	}

	driverWithAuth := func(driverType, outputFile string) (bool, testOverrides, string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
		authorizationModule := &model.ModuleType{
			Name:       model.ModuleTypeAuthorization,
			Version:    "0.2.1",
			Standalone: false,
		}

		// file name contents will overwrite values in each data
		confilePaths := []string{"testdata/auth_token.yaml", "testdata/auth_certificate.crt"}
		moduleConfig := fmt.Sprintf(`
		   karaviAuthorizationProxy.proxyAuthzToken.filename=%s
		   karaviAuthorizationProxy.rootCertificate.filename=%s
		`, confilePaths[0], confilePaths[1])

		array, app, stateChange := DriverComponents(driverType)

		app.ModuleTypes = []model.ModuleType{*authorizationModule}
		app.ModuleConfiguration = moduleConfig
		stateChange.ModuleTypes = []model.ModuleType{*authorizationModule}
		stateChange.ModuleConfiguration = moduleConfig

		return true, testOverrides{}, outputFile, confilePaths, array, app, stateChange
	}

	tests := map[string]func(t *testing.T) (bool, testOverrides, string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange){
		// successful cases
		"success observability standalone": func(*testing.T) (bool, testOverrides, string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {

			observabilityModule := &model.ModuleType{
				Name:       "observability",
				Standalone: true,
			}

			array := &model.StorageArray{}

			app := &model.Application{
				Name:        "standalone-module-app",
				Status:      model.DriverStatusOperational,
				ClusterID:   0,
				ModuleTypes: []model.ModuleType{*observabilityModule},
			}

			// file name contents will overwrite values in each data
			confilePaths := []string{"testdata/vxflexos_secret.yaml"}

			moduleConfig := fmt.Sprintf(`
				   karaviMetricsPowerflex.enabled=true
				   karaviMetricsPowerflex.driverConfig.filename=%s
				`, "testdata/vxflexos_secret.yaml")

			stateChange := &model.ApplicationStateChange{
				ApplicationID:       app.ID,
				ClusterID:           0,
				ModuleTypes:         []model.ModuleType{*observabilityModule},
				ModuleConfiguration: moduleConfig,
			}

			return true, testOverrides{}, "observability-secret.yaml", confilePaths, array, app, stateChange
		},
		"success observability standalone + authorization": func(*testing.T) (bool, testOverrides, string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {

			observabilityModule := &model.ModuleType{
				Name:       "observability",
				Standalone: true,
			}
			authorizationModule := &model.ModuleType{
				Name:       model.ModuleTypeAuthorization,
				Version:    "0.2.1",
				Standalone: false,
			}

			array := &model.StorageArray{}

			app := &model.Application{
				Name:          "oberv-with-auth",
				Status:        model.DriverStatusOperational,
				ClusterID:     0,
				ModuleTypes:   []model.ModuleType{*observabilityModule, *authorizationModule},
				StorageArrays: []model.StorageArray{*array},
			}

			// file name contents will overwrite values in each data
			confilePaths := []string{"testdata/vxflexos_secret.yaml", "testdata/auth_token.yaml", "testdata/auth_certificate.crt"}
			moduleConfig := fmt.Sprintf(`
				   observability.namespace=test-observability-auth-namespace
				   karaviMetricsPowerflex.enabled=true
				   karaviMetricsPowerflex.driverConfig.filename=%s
				   karaviAuthorizationProxy.proxyAuthzToken.filename=%s
		           karaviAuthorizationProxy.rootCertificate.filename=%s
				`, confilePaths[0], confilePaths[1], confilePaths[2])

			stateChange := &model.ApplicationStateChange{
				ApplicationID:       app.ID,
				ClusterID:           0,
				ModuleTypes:         []model.ModuleType{*observabilityModule, *authorizationModule},
				StorageArrays:       []model.StorageArray{*array},
				ModuleConfiguration: moduleConfig,
			}

			return true, testOverrides{}, "observability-secret.yaml", confilePaths, array, app, stateChange
		},
		"success observability + driver": func(*testing.T) (bool, testOverrides, string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {

			observabilityModule := &model.ModuleType{
				Name:       "observability",
				Standalone: true,
			}

			arrayType := &model.StorageArrayType{Name: model.ArrayTypePowerFlex}

			array := &model.StorageArray{
				StorageArrayType:   *arrayType,
				StorageArrayTypeID: arrayType.ID,
			}

			app := &model.Application{
				Name:         "oberv-with-driver",
				Status:       model.DriverStatusOperational,
				ClusterID:    0,
				DriverTypeID: 1,
				ModuleTypes:  []model.ModuleType{*observabilityModule},
			}

			stateChange := &model.ApplicationStateChange{
				ApplicationID:       app.ID,
				ClusterID:           0,
				DriverTypeID:        1,
				ModuleTypes:         []model.ModuleType{*observabilityModule},
				ModuleConfiguration: "",
			}

			return true, testOverrides{}, "driver-secret.yaml", []string{}, array, app, stateChange
		},
		"success driver-powerflex + authorization": func(*testing.T) (bool, testOverrides, string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			return driverWithAuth(model.ArrayTypePowerFlex, "driver-secret.yaml")
		},
		"success driver-powermax + authorization": func(*testing.T) (bool, testOverrides, string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			return driverWithAuth(model.ArrayTypePowerMax, "driver-secret-powermax.yaml")
		},
		"success driver-powerscale + authorization": func(*testing.T) (bool, testOverrides, string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			return driverWithAuth(model.ArrayTypePowerScale, "driver-secret.yaml")
		},

		"success driver-powerflex only": func(*testing.T) (bool, testOverrides, string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			return driverNoError(model.ArrayTypePowerFlex, "driver-secret.yaml")
		},
		"success driver-powerscale only": func(*testing.T) (bool, testOverrides, string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			return driverNoError(model.ArrayTypePowerScale, "driver-secret.yaml")
		},
		"sucess driver-powerstore only": func(*testing.T) (bool, testOverrides, string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			return driverNoError(model.ArrayTypePowerStore, "driver-secret.yaml")
		},
		"success driver-powermax only": func(*testing.T) (bool, testOverrides, string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			return driverNoError(model.ArrayTypePowerMax, "driver-secret-powermax.yaml")
		},
		"success driver-unity only": func(*testing.T) (bool, testOverrides, string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			return driverNoError(model.ArrayTypeUnity, "driver-secret.yaml")
		},
		// failure cases
		"fail invalid ytt value format": func(*testing.T) (bool, testOverrides, string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {

			observabilityModule := &model.ModuleType{
				Name:       "observability",
				Standalone: true,
			}

			array := &model.StorageArray{}

			app := &model.Application{
				Name:        "standalone-module-app",
				Status:      model.DriverStatusOperational,
				ClusterID:   0,
				ModuleTypes: []model.ModuleType{*observabilityModule},
			}

			// file name contents will overwrite values in each data
			confilePaths := []string{"testdata/vxflexos_secret.yaml"}

			moduleConfig := fmt.Sprintf(`
			   karaviMetricsPowerflex.enabled=true=1
			   karaviMetricsPowerflex.driverConfig.filename=%s
			`, "testdata/vxflexos_secret.yaml")

			stateChange := &model.ApplicationStateChange{
				ApplicationID:       app.ID,
				ClusterID:           0,
				ModuleTypes:         []model.ModuleType{*observabilityModule},
				ModuleConfiguration: moduleConfig,
			}

			return false, testOverrides{}, "", confilePaths, array, app, stateChange
		},
		"fail couldn't find storage arrays for app state with": func(*testing.T) (bool, testOverrides, string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			array, app, stateChange := DriverComponents(model.ArrayTypeUnity)
			patch := testOverrides{
				skipArrayCreation: true,
			}
			return false, patch, "", []string{}, array, app, stateChange
		},
		"fail driver-powerflex only": func(*testing.T) (bool, testOverrides, string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			return driverError(model.ArrayTypePowerFlex, "")
		},
		"fail driver-powerscale only": func(*testing.T) (bool, testOverrides, string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			return driverError(model.ArrayTypePowerScale, "")
		},
		"fail driver-powerstore only": func(*testing.T) (bool, testOverrides, string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			return driverError(model.ArrayTypePowerStore, "")
		},
		"fail driver-powermax only": func(*testing.T) (bool, testOverrides, string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			return driverError(model.ArrayTypePowerMax, "")
		},
		"fail driver-unity only": func(*testing.T) (bool, testOverrides, string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			return driverError(model.ArrayTypeUnity, "")
		},
		"error querying app state db": func(*testing.T) (bool, testOverrides, string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			array, app, stateChange := DriverComponents(model.ArrayTypePowerFlex)
			patch := testOverrides{
				errorGetByIDAppState: true,
			}
			return false, patch, "", []string{}, array, app, stateChange
		},
		"nil result from db": func(*testing.T) (bool, testOverrides, string, []string, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			array, app, stateChange := DriverComponents(model.ArrayTypePowerFlex)
			patch := testOverrides{
				nilGetByIDAppState: true,
			}
			return false, patch, "", []string{}, array, app, stateChange
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			success, patch, outputFile, configFilePaths, array, app, stateChange := tc(t)

			if patch.decryptPassword != nil {
				oldDecryptPassword := utils.DecryptPassword
				defer func() { utils.DecryptPassword = oldDecryptPassword }()
				utils.DecryptPassword = patch.decryptPassword
			}

			generateDynamicSecretWrapper := func(appID uint, suite *YttTestSuite, app *model.Application, stateChange *model.ApplicationStateChange) (ytt.Output, error) {
				client := ytt.NewClient()
				client.SetOptions(ytt.WithTemplatePath("../"))

				application := store.NewApplicationStateChangeStore(suite.db)
				cf := store.NewConfigFileStore(suite.db)

				for _, path := range configFilePaths {
					cfData, err := ioutil.ReadFile(path)
					if err != nil {
						panic(err)
					}
					cfInstance := &model.ConfigFile{Name: path, ConfigFileData: cfData}
					suite.db.Create(cfInstance)
				}

				suite.db.Create(app)
				suite.db.Create(stateChange)

				if patch.errorGetByIDAppState {
					ctrl := gomock.NewController(t)

					as := mocks.NewMockApplicationStateChangeStoreInterface(ctrl)
					as.EXPECT().GetByID(appID).Times(1).Return(nil, errors.New("error"))
					defer func() { ctrl.Finish() }()

					return client.GenerateDynamicSecret(appID, as, cf)

				} else if patch.nilGetByIDAppState {
					ctrl := gomock.NewController(t)

					as := mocks.NewMockApplicationStateChangeStoreInterface(ctrl)
					as.EXPECT().GetByID(appID).Times(1).Return(nil, nil)
					defer func() { ctrl.Finish() }()

					return client.GenerateDynamicSecret(appID, as, cf)

				}

				return client.GenerateDynamicSecret(appID, application, cf)

			}

			RunTestSuite(t, generateDynamicSecretWrapper, []string{outputFile}, success, patch.skipArrayCreation, array, app, stateChange)
		})
	}
}

func Test_ConfigMapTemplateFromApplication(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, testOverrides, *model.StorageArray, *model.Application, *model.ApplicationStateChange){
		"success driver": func(*testing.T) (bool, testOverrides, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {

			array, app, stateChange := DriverComponents(model.ArrayTypePowerFlex)
			return true, testOverrides{}, array, app, stateChange
		},
		"error querying app state db": func(*testing.T) (bool, testOverrides, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			array, app, stateChange := DriverComponents(model.ArrayTypePowerFlex)
			patch := testOverrides{
				errorGetByIDAppState: true,
			}
			return false, patch, array, app, stateChange
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			success, patch, array, app, stateChange := tc(t)

			configMapTemplateFromApplicationWrapper := func(appID uint, suite *YttTestSuite, app *model.Application, stateChange *model.ApplicationStateChange) (ytt.Output, error) {
				application := store.NewApplicationStateChangeStore(suite.db)

				suite.db.Create(app)
				suite.db.Create(stateChange)

				client := ytt.NewClient()
				client.SetOptions(ytt.WithTemplatePath("../"))

				if patch.errorGetByIDAppState {
					ctrl := gomock.NewController(t)

					as := mocks.NewMockApplicationStateChangeStoreInterface(ctrl)
					as.EXPECT().GetByID(appID).Times(1).Return(nil, errors.New("error"))
					defer func() { ctrl.Finish() }()

					return client.ConfigMapTemplateFromApplication(appID, as)

				}

				return client.ConfigMapTemplateFromApplication(appID, application)

			}

			RunTestSuite(t, configMapTemplateFromApplicationWrapper, []string{"driver-config-params.yaml"}, success, patch.skipArrayCreation, array, app, stateChange)

		})
	}
}

func Test_GetEmptySecret(t *testing.T) {
	tests := map[string]func(t *testing.T) (bool, testOverrides, *model.StorageArray, *model.Application, *model.ApplicationStateChange){
		"success driver": func(*testing.T) (bool, testOverrides, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			array, app, stateChange := DriverComponents(model.ArrayTypePowerFlex)
			return true, testOverrides{}, array, app, stateChange
		},
		"fail driver no storage": func(*testing.T) (bool, testOverrides, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			array, app, stateChange := DriverComponents(model.ArrayTypePowerFlex)
			patch := testOverrides{
				skipArrayCreation: true,
			}
			return false, patch, array, app, stateChange
		},
		"error querying app state db": func(*testing.T) (bool, testOverrides, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			array, app, stateChange := DriverComponents(model.ArrayTypePowerFlex)
			patch := testOverrides{
				errorGetByIDAppState: true,
			}
			return false, patch, array, app, stateChange
		},
		"nil result from db": func(*testing.T) (bool, testOverrides, *model.StorageArray, *model.Application, *model.ApplicationStateChange) {
			array, app, stateChange := DriverComponents(model.ArrayTypePowerFlex)
			patch := testOverrides{
				nilGetByIDAppState: true,
			}
			return false, patch, array, app, stateChange
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			success, patch, array, app, stateChange := tc(t)

			getEmptySecretWrapper := func(appID uint, suite *YttTestSuite, app *model.Application, stateChange *model.ApplicationStateChange) (ytt.Output, error) {
				client := ytt.NewClient()
				client.SetOptions(ytt.WithTemplatePath("../"))

				if patch.errorGetByIDAppState {
					ctrl := gomock.NewController(t)

					as := mocks.NewMockApplicationStateChangeStoreInterface(ctrl)
					as.EXPECT().GetByID(appID).Times(1).Return(nil, errors.New("error"))
					defer func() { ctrl.Finish() }()

					return client.GetEmptySecret(appID, as)

				} else if patch.nilGetByIDAppState {
					ctrl := gomock.NewController(t)

					as := mocks.NewMockApplicationStateChangeStoreInterface(ctrl)
					as.EXPECT().GetByID(appID).Times(1).Return(nil, nil)
					defer func() { ctrl.Finish() }()

					return client.GetEmptySecret(appID, as)

				}

				application := store.NewApplicationStateChangeStore(suite.db)

				suite.db.Create(app)
				suite.db.Create(stateChange)

				return client.GetEmptySecret(appID, application)

			}

			RunTestSuite(t, getEmptySecretWrapper, []string{"empty-secret.yaml"}, success, patch.skipArrayCreation, array, app, stateChange)
		})
	}
}

func Test_EchoLoggerWrapper(t *testing.T) {
	r, w, _ := os.Pipe()

	logger := echo.New().Logger
	logger.SetOutput(w)
	logger.SetLevel(1)

	client := ytt.NewClient()
	client.SetOptions(ytt.WithLogger(logger, true))

	wLogger := ytt.EchoLoggerWrapper{
		EchoLogger: logger,
		Debug:      true,
	}

	wLogger.Printf("%s", "Printf- hello")
	wLogger.Debugf("%s", "Debugf- every")
	wLogger.Warnf("%s", "Warnf-one")

	err := w.Close()
	assert.NoError(t, err, "Error on OS pipe close")
	out, err := ioutil.ReadAll(r)
	assert.NoError(t, err)

	outSTr := fmt.Sprintf("%s", out)
	assert.True(t, strings.Contains(outSTr, "Printf- hello"))
	assert.True(t, strings.Contains(outSTr, "Debugf- every"))
	assert.True(t, strings.Contains(outSTr, "Warnf-one"))

	fmt.Print(outSTr)

	assert.EqualValues(t, os.Stderr, wLogger.DebugWriter())
	wLogger.Debug = false
	assert.EqualValues(t, ytt.NoopWriter{}, wLogger.DebugWriter())

	d, err := wLogger.DebugWriter().Write([]byte(""))
	assert.NoError(t, err)
	assert.EqualValues(t, 0, d)
}

func (suite *YttTestSuite) SetupSuite() {
	suite.inputPaths = append(suite.inputPaths, "testdata/")
	suite.outputPath = "testdata/output/"
	// we will create separate test DB files for each packages to avoid data duplication issue and constraint errors
	suite.db = createTestDB("./../csm_test.db")
	err := db.AutoMigrate(suite.db)
	assert.NoError(suite.T(), err, "Error in DB migration")

	path := testFilePath{
		configPath:     "testdata/kubeconfig.yaml",
		authzTokenPath: "testdata/auth_token.yaml",
		authCertPath:   "testdata/auth_certificate.crt",
	}
	populateTestDb(suite.db, path)
}

func (suite *YttTestSuite) TearDownSuite() {
	err := os.RemoveAll(suite.outputPath)
	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), dropTestDB("./../csm_test.db"))
}

func (suite *YttTestSuite) TestTemplate() {
	client := ytt.NewClient()
	output, err := client.Template(suite.inputPaths, []string{"port=9090", "text=\"hey, people\""})
	suite.NoError(err)

	err = output.CreateAt(suite.outputPath)
	suite.NoError(err)
	suite.FileExists(filepath.Join(suite.outputPath, "config.yaml"))

	data := output.AsBytes()
	suite.Equal(4, len(data))
	spew.Dump(data)
}

func populateTestDb(db *gorm.DB, path testFilePath) {
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

	encrypted, err := utils.EncryptPassword([]byte("password"))
	if err != nil {
		panic(err)
	}
	array := &model.StorageArray{
		UniqueID:           "id-1",
		Username:           "user",
		Password:           encrypted,
		ManagementEndpoint: "127.0.01",
		StorageArrayTypeID: arrayType.ID,
	}
	db.Create(array)

	powerflexdriver14 := &model.DriverType{
		Version:            "1.4.0",
		StorageArrayTypeID: arrayType.ID,
	}
	db.Create(powerflexdriver14)

	podmonModule := &model.ModuleType{
		Name:       model.ModuleTypePodMon,
		Standalone: false,
	}
	db.Create(podmonModule)

	observabilityModule := &model.ModuleType{
		Name:       model.ModuleTypeObservability,
		Standalone: true,
	}
	db.Create(observabilityModule)

	authorizationModule := &model.ModuleType{
		Name:       model.ModuleTypeAuthorization,
		Version:    "0.2.1",
		Standalone: false,
	}
	db.Create(authorizationModule)

	reverseproxyModule := &model.ModuleType{
		Name:       model.ModuleTypeReverseProxy,
		Version:    "1.3.0",
		Standalone: false,
	}
	db.Create(reverseproxyModule)

	kubeConfigData, err := ioutil.ReadFile(path.configPath)
	if err != nil {
		panic(err)
	}
	cluster := &model.Cluster{
		ClusterName:      "test-cluster",
		ConfigFileData:   kubeConfigData,
		OrchestratorType: "k8s",
		Status:           model.ClusterStatusConnected,
		K8sVersion:       "1.22",
	}
	if err := db.Create(cluster).Error; err != nil {
		panic(err)
	}

	// file name contents will overwrite values in each data
	moduleConfig := fmt.Sprintf(`
	   karaviAuthorizationProxy.proxyAuthzToken.filename=%s
	   karaviAuthorizationProxy.rootCertificate.filename=%s
	`, path.authzTokenPath, path.authCertPath)

	authzTokenConfigData, err := ioutil.ReadFile(path.authzTokenPath)
	if err != nil {
		panic(err)
	}
	cf := &model.ConfigFile{Name: path.authzTokenPath, ConfigFileData: authzTokenConfigData}
	db.Create(cf)

	authCertConfigData, err := ioutil.ReadFile(path.authCertPath)
	if err != nil {
		panic(err)
	}
	cf = &model.ConfigFile{Name: path.authCertPath, ConfigFileData: authCertConfigData}
	db.Create(cf)

	app := &model.Application{
		Name:                "test-app",
		Status:              model.DriverStatusOperational,
		ClusterID:           cluster.ID,
		DriverTypeID:        1,
		ModuleTypes:         []model.ModuleType{*podmonModule, *authorizationModule, *observabilityModule},
		StorageArrays:       []model.StorageArray{*array},
		ModuleConfiguration: moduleConfig,
	}
	db.Create(app)

	stateChange := &model.ApplicationStateChange{
		ApplicationID:       app.ID,
		ClusterID:           cluster.ID,
		DriverTypeID:        1,
		ModuleTypes:         []model.ModuleType{*podmonModule, *authorizationModule, *observabilityModule},
		StorageArrays:       []model.StorageArray{*array},
		ModuleConfiguration: moduleConfig,
	}
	db.Create(stateChange)
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

func TestYttSuite(t *testing.T) {
	suite.Run(t, new(YttTestSuite))
}
