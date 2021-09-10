// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package ytt

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/dell/csm-deployment/utils"
	"github.com/dell/csm-deployment/utils/constants"

	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/store"
	"github.com/k14s/ytt/pkg/cmd/template"
	"github.com/k14s/ytt/pkg/cmd/ui"
	"github.com/k14s/ytt/pkg/files"
	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v3"
)

const (
	// NamespacePowerflex - Placeholder for powerflex namespace
	NamespacePowerflex = "vxflexos"
)

//EchoLoggerWrapper - Define log wrapper object
type EchoLoggerWrapper struct {
	EchoLogger echo.Logger
	Debug      bool
}

// Printf - Printf wrapper
func (e *EchoLoggerWrapper) Printf(s string, i ...interface{}) {
	e.EchoLogger.Printf(s, i)
}

// Debugf - Debug log wrapper
func (e *EchoLoggerWrapper) Debugf(s string, i ...interface{}) {
	if e.Debug {
		e.EchoLogger.Debugf(s, i)
	}
}

// Warnf - Warn log wrapper
func (e *EchoLoggerWrapper) Warnf(s string, i ...interface{}) {
	e.EchoLogger.Warnf(s, i)
}

// DebugWriter - Print err if debug enabled
func (e *EchoLoggerWrapper) DebugWriter() io.Writer {
	if e.Debug {
		return os.Stderr
	}
	return NoopWriter{}
}

// NoopWriter -
type NoopWriter struct{}

var _ io.Writer = NoopWriter{}

func (w NoopWriter) Write(data []byte) (int, error) { return len(data), nil }

// Output -
type Output struct {
	*template.Output
}

// CreateAt - Create file at outputPath
func (out *Output) CreateAt(outputPath string) error {
	for _, file := range out.Files {
		err := file.Create(outputPath)
		if err != nil {
			return err
		}
	}
	return nil
}

//AsBytes - Append file bytes
func (out *Output) AsBytes() [][]byte {
	var res [][]byte
	for _, file := range out.Files {
		res = append(res, file.Bytes())
	}
	return res
}

// AsCombinedBytes - Append file bytes
func (out *Output) AsCombinedBytes() []byte {
	var res []byte
	for _, file := range out.Files {
		res = append(res, file.Bytes()...)
		res = append(res, []byte("---\n")...)
	}
	return res
}

// Interface -
//go:generate mockgen -destination=mocks/ytt_interface.go -package=mocks github.com/dell/csm-deployment/ytt Interface
type Interface interface {
	Template([]string, []string) (Output, error)
	TemplateFromApplication(appID uint,
		as store.ApplicationStateChangeStoreInterface, cs store.ClusterStoreInterface, cf store.ConfigFileStoreInterface) (Output, error)
	NamespaceTemplateFromApplication(appID uint,
		as store.ApplicationStateChangeStoreInterface) (Output, error)
	GetEmptySecret(appID uint, as store.ApplicationStateChangeStoreInterface) (Output, error)
	GenerateDynamicSecret(appID uint, as store.ApplicationStateChangeStoreInterface, cf store.ConfigFileStoreInterface) (Output, error)
	ConfigMapTemplateFromApplication(appID uint, as store.ApplicationStateChangeStoreInterface) (Output, error)
	SetOptions(opts ...Option)
}

type client struct {
	logger       ui.UI
	templatePath string
}

// Option -
type Option func(*client)

// GlobalConfigStorageArrays -
type GlobalConfigStorageArrays []struct {
	ID             string `yaml:"storageArrayId"`
	Endpoint       string `yaml:"endpoint"`
	BackupEndpoint string `yaml:"backupEndpoint,omitempty"`
}

// GlobalConfigManagementServers -
type GlobalConfigManagementServers []struct {
	Endpoint                  string `yaml:"endpoint"`
	CredentialsSecret         string `yaml:"credentialsSecret,omitempty"`
	CertSecret                string `yaml:"certSecret,omitempty"`
	SkipCertificateValidation bool   `yaml:"skipCertificateValidation,omitempty"`
	Limits                    struct {
		MaxActiveRead       int `yaml:"maxActiveRead,omitempty"`
		MaxActiveWrite      int `yaml:"maxActiveWrite,omitempty"`
		MaxOutStandingRead  int `yaml:"maxOutStandingRead,omitempty"`
		MaxOutStandingWrite int `yaml:"maxOutStandingWrite,omitempty"`
	} `yaml:"limits,omitempty"`
}

// WithLogger - Set logger
func WithLogger(logger echo.Logger, debug bool) Option {
	return func(c *client) {
		c.logger = &EchoLoggerWrapper{EchoLogger: logger, Debug: debug}
	}
}

// WithTemplatePath - Set template path to client
func WithTemplatePath(path string) Option {
	return func(c *client) {
		c.templatePath = path
	}
}

var defaultOptions Option = func(c *client) {
	c.logger = ui.NewTTY(false)
	c.templatePath = "./"
}

// NewClient will return a new ytt client
func NewClient() Interface {
	c := &client{}
	defaultOptions(c)
	return c
}

// SetOptions - Sets options on a client
func (c *client) SetOptions(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
}

func (c *client) Template(paths []string, values []string) (Output, error) {
	yttFiles, err := files.NewSortedFilesFromPaths(paths, files.SymlinkAllowOpts{})
	if err != nil {
		return Output{}, err
	}
	in := template.Input{Files: yttFiles}

	opts := template.NewOptions()
	opts.DataValuesFlags = template.DataValuesFlags{KVsFromYAML: values}

	opts.Debug = true

	output := opts.RunWithFiles(in, c.logger)

	return Output{&output}, output.Err
}

// ProcessModuleConfig overwires specif value-key pair based on other value key-pair
var ProcessModuleConfig = func(mcf []string, enabledModules map[string]string, cf store.ConfigFileStoreInterface) ([]string, error) {
	var newModuleConfig []string
	for _, v := range mcf {
		configFile := strings.Split(v, "=")
		if len(configFile) != 2 {
			return nil, errors.New("invalid ytt value format. It should be key1=value1")
		}

		switch configFile[0] {
		case "karaviAuthorizationProxy.proxyAuthzToken.filename":
			type Token struct {
				Data struct {
					Access  string `yaml:"access"`
					Refresh string `yaml:"refresh"`
				} `yaml:"data"`
			}

			token := Token{}

			file, err := cf.GetAllByName(configFile[1])
			if err != nil {
				return nil, err
			}

			if err := yaml.Unmarshal(file[0].ConfigFileData, &token); err != nil {
				return nil, err
			}

			newModuleConfig = append(newModuleConfig, fmt.Sprintf("karaviAuthorizationProxy.proxyAuthzToken.data.access=%s", token.Data.Access))
			newModuleConfig = append(newModuleConfig, fmt.Sprintf("karaviAuthorizationProxy.proxyAuthzToken.data.refresh=%s", token.Data.Refresh))

		case "karaviAuthorizationProxy.rootCertificate.filename":
			file, err := cf.GetAllByName(configFile[1])
			if err != nil {
				return nil, err
			}
			newModuleConfig = append(newModuleConfig, fmt.Sprintf("karaviAuthorizationProxy.rootCertificate.data=%s", file[0].ConfigFileData))
		case "karaviMetricsPowerflex.driverConfig.filename":
			file, err := cf.GetAllByName(configFile[1])
			if err != nil {
				return nil, err
			}
			newModuleConfig = append(newModuleConfig, fmt.Sprintf("karaviMetricsPowerflex.driverConfig.data=%s", file[0].ConfigFileData))
			// set the secret for auth injection if Auth is enabled
			if _, ok := enabledModules[model.ModuleTypeAuthorization]; ok {
				type SecretConfig struct {
					Data struct {
						Config string `yaml:"config"`
					} `yaml:"data"`
				}

				secretConfig := SecretConfig{}

				if err := yaml.Unmarshal(file[0].ConfigFileData, &secretConfig); err != nil {
					return nil, err
				}

				newModuleConfig = append(newModuleConfig, fmt.Sprintf("karaviMetricsPowerflex.driverConfig.onlyConfig=%s", secretConfig.Data.Config))
			}

		default:
			newModuleConfig = append(newModuleConfig, v)
		}
	}

	return newModuleConfig, nil
}

func (c *client) TemplateFromApplication(appID uint,
	as store.ApplicationStateChangeStoreInterface, cs store.ClusterStoreInterface, cf store.ConfigFileStoreInterface) (Output, error) {
	c.logger.Printf("Generating template from app state with id %d \n", appID)
	appState, err := as.GetByID(appID)
	if err != nil {
		return Output{}, err
	}
	if appState == nil {
		return Output{}, fmt.Errorf("app state with id %d not found", appID)
	}

	if len(appState.StorageArrays) == 0 {
		return Output{}, fmt.Errorf("couldn't find storage arrays for app state with %d", appID)
	}
	storageArray := appState.StorageArrays[0]
	arrayType := storageArray.StorageArrayType
	c.logger.Printf("Array type is %q \n", arrayType.Name)

	cluster, err := cs.GetByID(appState.ClusterID)
	if err != nil {
		return Output{}, err
	}

	var values []string

	// only standalone modules can be installed without a driver
	if appState.DriverTypeID == 0 {
		for _, module := range appState.ModuleTypes {
			if !(module.Standalone || (module.Name == model.ModuleTypeAuthorization)) { //exception is authorization, it has been handled in precheck
				return Output{}, fmt.Errorf("unable to install module %s without specifying a driver", module.Name)
			}
		}
		values = append(values, "observability.standalone=true")
	}

	paths := make([]string, 0)

	// if driver is going to be installed, load appropriate template files
	if appState.DriverTypeID > 0 {
		paths = append(paths, []string{
			c.templatePath + "templates/controller.yaml",
			c.templatePath + "templates/node.yaml",
			c.templatePath + "templates/csidriver.yaml",
			c.templatePath + fmt.Sprintf("templates/configs/values-%s.yaml", arrayType.Name),
			c.templatePath + "templates/common/values.yaml",
			c.templatePath + fmt.Sprintf("templates/common/k8s-%s-values.yaml", cluster.K8sVersion),
			c.templatePath + "templates/configs/values-authorization.yaml",
			c.templatePath + "templates/authorization.yaml",
			c.templatePath + "templates/csireverseproxy.yaml",
		}...)

		if arrayType.Name == model.ArrayTypePowerMax {
			globalConfigStorageArrays := GlobalConfigStorageArrays{
				{ID: storageArray.UniqueID, Endpoint: storageArray.ManagementEndpoint},
			}
			globalConfigStorageArraysBytes, err := yaml.Marshal(globalConfigStorageArrays)
			if err != nil {
				return Output{}, err
			}

			globalConfigManagementServers := GlobalConfigManagementServers{
				{Endpoint: storageArray.ManagementEndpoint},
			}
			globalConfigManagementServersBytes, err := yaml.Marshal(globalConfigManagementServers)
			if err != nil {
				return Output{}, err
			}

			values = append(values, fmt.Sprintf("globalConfig.storageArrays=%s", globalConfigStorageArraysBytes))
			values = append(values, fmt.Sprintf("globalConfig.managementServers=%s", globalConfigManagementServersBytes))
			values = append(values, fmt.Sprintf(`portGroups="%s"`, utils.GetValueFromMetadataKey(storageArray.MetaData, constants.KeyPortGroups)))
		}
	}

	// load standalone modules (ex: observability) and other shared ytt module files
	paths = append(paths, []string{
		c.templatePath + "templates/modules/",
		c.templatePath + "templates/configs/values-observability.yaml",
		c.templatePath + "templates/observability.yaml",
	}...)

	enabledModules := make(map[string]string)
	for _, module := range appState.ModuleTypes {
		values = append(values, fmt.Sprintf("%s.enabled=true", module.Name))
		enabledModules[module.Name] = module.Name
	}

	// add configuration values for the driver and modules that were passed from the API to create the application
	// these are space-delimited strings of the format "parent1.key1=value1 parent2.key2=value2"
	if len(appState.DriverConfiguration) > 0 {
		driverValues := strings.Split(appState.DriverConfiguration, " ")
		values = append(values, driverValues...)
	}

	// enable CSM Metrics based on driver type
	switch arrayType.Name {
	case model.ArrayTypePowerFlex:
		values = append(values, "karaviMetricsPowerflex.enabled=true")
	case model.ArrayTypePowerStore:
		values = append(values, "karaviMetricsPowerstore.enabled=true")
	}

	if len(appState.ModuleTypes) != 0 {
		if appState.ModuleTypes[0].Name != model.ModuleTypeReplication { // we don't want any additional module configs for the replication.
			if len(appState.ModuleConfiguration) > 0 {
				re := regexp.MustCompile(`[^\s]+`)
				moduleValues, err := ProcessModuleConfig(re.FindAllString(appState.ModuleConfiguration, -1), enabledModules, cf)
				if err != nil {
					return Output{}, err
				}
				values = append(values, moduleValues...)
			}
		}
	}

	output, err := c.Template(paths, values)
	if err != nil {
		return Output{}, err
	}
	return output, output.Err
}

func (c *client) GetEmptySecret(appID uint, as store.ApplicationStateChangeStoreInterface) (Output, error) {
	c.logger.Printf("Generating secret template from app state with id %d \n", appID)
	appState, err := as.GetByID(appID)
	if err != nil {
		return Output{}, err
	}
	if appState == nil {
		return Output{}, fmt.Errorf("app state with id %d not found", appID)
	}

	if len(appState.StorageArrays) == 0 {
		return Output{}, fmt.Errorf("couldn't find storage arrays for app state with %d", appID)
	}
	arrayType := appState.StorageArrays[0].StorageArrayType
	c.logger.Printf("Array type is %q \n", arrayType.Name)
	paths := []string{
		c.templatePath + fmt.Sprintf("templates/configs/values-%s.yaml", arrayType.Name),
		c.templatePath + "templates/empty-secret.yaml",
	}

	output, err := c.Template(paths, nil)
	if err != nil {
		return Output{}, err
	}
	return output, output.Err
}

func (c *client) NamespaceTemplateFromApplication(appID uint, as store.ApplicationStateChangeStoreInterface) (Output, error) {
	c.logger.Printf("Generating namespace template from app state with id %d \n", appID)
	appState, err := as.GetByID(appID)
	if err != nil {
		return Output{}, err
	}
	if appState == nil {
		return Output{}, fmt.Errorf("app state with id %d not found", appID)
	}

	var values []string
	var paths []string

	if appState.DriverTypeID > 0 {
		if len(appState.StorageArrays) == 0 {
			return Output{}, fmt.Errorf("couldn't find storage arrays for app state with %d", appID)
		}
		arrayType := appState.StorageArrays[0].StorageArrayType
		paths = append(paths, c.templatePath+fmt.Sprintf("templates/configs/values-%s.yaml", arrayType.Name))
	} else {
		// TODO(Michael): we are assuming the only standalone module is observability, may be modified when other standalone are supported
		if !strings.Contains(appState.ModuleConfiguration, "namespace=") {
			values = append(values, fmt.Sprintf("namespace=%s", model.ObservabilityNamespace))
		}
		values = append(values, "observability.standalone=true")
	}
	paths = append(paths, c.templatePath+"templates/namespace.yaml")

	output, err := c.Template(paths, values)
	if err != nil {
		return Output{}, err
	}
	return output, output.Err
}

func (c *client) DynamicTemplate(paths []string, values []string) (Output, error) {
	yttFiles, err := files.NewSortedFilesFromPaths(paths, files.SymlinkAllowOpts{})
	if err != nil {
		return Output{}, err
	}
	in := template.Input{Files: yttFiles}

	opts := template.NewOptions()

	opts.DataValuesFlags = template.DataValuesFlags{
		EnvFromStrings: []string{"DVS"},
		EnvironFunc:    func() []string { return []string{"DVS_str=str"} },
		KVsFromStrings: values,
	}
	opts.Debug = true
	output := opts.RunWithFiles(in, c.logger)
	return Output{&output}, output.Err
}

func (c *client) ConfigMapTemplateFromApplication(appID uint, as store.ApplicationStateChangeStoreInterface) (Output, error) {

	c.logger.Printf("Generating config map params template from app state with id %d \n", appID)
	appState, err := as.GetByID(appID)
	if err != nil {
		return Output{}, err
	}

	arrayType := appState.StorageArrays[0].StorageArrayType

	paths := []string{
		c.templatePath + fmt.Sprintf("templates/configs/values-%s.yaml", arrayType.Name),
		c.templatePath + "templates/configs/driver-config-params.yaml",
	}

	output, err := c.Template(paths, nil)
	if err != nil {
		return Output{}, err
	}
	return output, nil
}

func (c *client) GenerateDynamicSecret(appID uint, as store.ApplicationStateChangeStoreInterface, cf store.ConfigFileStoreInterface) (Output, error) {

	var valuesArray []string
	var paths []string

	c.logger.Printf("Generating dynamic secret template from app state with id %d \n", appID)
	appState, err := as.GetByID(appID)
	if err != nil {
		return Output{}, err
	}
	if appState == nil {
		return Output{}, fmt.Errorf("app state with id %d not found", appID)
	}

	tmpValuesArray := []string{}
	enabledModules := make(map[string]string)
	for _, module := range appState.ModuleTypes {
		tmpValuesArray = append(tmpValuesArray, fmt.Sprintf("%s.enabled=true", module.Name))
		enabledModules[module.Name] = module.Name
	}

	if appState.DriverTypeID == 0 {
		// TODO(Michael): we are assuming the only standalone module is observability, may be modified when other standalone are supported
		if !strings.Contains(appState.ModuleConfiguration, "namespace=") {
			tmpValuesArray = append(tmpValuesArray, fmt.Sprintf("namespace=%s", model.ObservabilityNamespace))
		}
		tmpValuesArray = append(tmpValuesArray, "observability.standalone=true")

	} else {
		if len(appState.StorageArrays) == 0 {
			return Output{}, fmt.Errorf("couldn't find storage arrays for app state with %d", appID)
		}

		storageArray := appState.StorageArrays[0]
		arrayType := storageArray.StorageArrayType

		c.logger.Printf("Array type is %q \n", arrayType.Name)

		paths = append(paths, c.templatePath+"templates/configs/driver-secret.yaml")

		storageArrayType := arrayType.Name

		switch arrayType.Name {

		case model.ArrayTypeUnity:
			valuesArray, err = getUnitySecret(storageArray)
			if err != nil {
				c.logger.Printf("Error during secret generation for %s : %v", arrayType.Name, err)
				return Output{}, err
			}
		case model.ArrayTypePowerMax:
			paths = []string{
				c.templatePath + fmt.Sprintf("templates/configs/driver-secret-%s.yaml", model.ArrayTypePowerMax),
			}
			valuesArray, err = getPowermaxSecret(storageArray)
			if err != nil {
				c.logger.Printf("Error during secret generation for %s : %v", arrayType.Name, err)
				return Output{}, err
			}

			globalConfigStorageArrays := GlobalConfigStorageArrays{
				{ID: storageArray.UniqueID, Endpoint: storageArray.ManagementEndpoint},
			}
			globalConfigStorageArraysBytes, err := yaml.Marshal(globalConfigStorageArrays)
			if err != nil {
				return Output{}, err
			}

			globalConfigManagementServers := GlobalConfigManagementServers{
				{Endpoint: storageArray.ManagementEndpoint},
			}
			globalConfigManagementServersBytes, err := yaml.Marshal(globalConfigManagementServers)
			if err != nil {
				return Output{}, err
			}

			fmt.Printf("%s\n", globalConfigManagementServersBytes)
			fmt.Printf("%s\n", globalConfigStorageArraysBytes)

			valuesArray = append(valuesArray, fmt.Sprintf("globalConfig.storageArrays=%s", globalConfigStorageArraysBytes))
			valuesArray = append(valuesArray, fmt.Sprintf("globalConfig.managementServers=%s", globalConfigManagementServersBytes))
		case model.ArrayTypePowerStore:
			valuesArray, err = getPowerstoreSecret(storageArray)
			if err != nil {
				c.logger.Printf("Error during secret generation for %s : %v", arrayType.Name, err)
				return Output{}, err
			}
		case model.ArrayTypePowerFlex:
			storageArrayType = NamespacePowerflex
			valuesArray, err = getPowerflexSecret(storageArray)
			if err != nil {
				c.logger.Printf("Error during secret generation for %s : %v", arrayType.Name, err)
				return Output{}, err
			}
		case model.ArrayTypePowerScale:
			valuesArray, err = getPowerscaleSecret(storageArray)
			if err != nil {
				c.logger.Printf("Error during secret generation for %s : %v", arrayType.Name, err)
				return Output{}, err
			}
		}

		valuesArray = append(valuesArray, fmt.Sprintf("secret.name=%s-config", storageArrayType))
		valuesArray = append(valuesArray, fmt.Sprintf("secret.namespace=%s", storageArrayType))

		paths = append(paths, c.templatePath+fmt.Sprintf("templates/configs/values-%s.yaml", arrayType.Name))

	}

	paths = append(paths, []string{
		c.templatePath + "templates/authorization-secret.yaml",
		c.templatePath + "templates/configs/values-authorization.yaml",
		c.templatePath + "templates/observability-secret.yaml",
		c.templatePath + "templates/configs/values-observability.yaml",
		c.templatePath + "templates/modules/",
	}...)

	if len(appState.ModuleConfiguration) > 0 {
		re := regexp.MustCompile(`[^\s]+`)
		moduleValues, err := ProcessModuleConfig(re.FindAllString(appState.ModuleConfiguration, -1), enabledModules, cf)
		if err != nil {
			return Output{}, err
		}
		tmpValuesArray = append(tmpValuesArray, moduleValues...)
	}

	valuesArray = append(valuesArray, tmpValuesArray...)
	output, err := c.DynamicTemplate(paths, valuesArray)
	if err != nil {
		return Output{}, err
	}
	return output, output.Err
}

func getUnitySecret(array model.StorageArray) ([]string, error) {

	defaultArray, skipCertificateValidation, err := getBoolValuesFromMetadata(array.MetaData)
	if err != nil {
		return nil, err
	}
	decryptedPassword, err := utils.DecryptPassword(array.Password)
	if err != nil {
		return nil, err
	}

	unityArray := model.UnityStorageArray{
		ArrayID:                   array.UniqueID,
		Username:                  array.Username,
		Password:                  string(decryptedPassword),
		Endpoint:                  array.ManagementEndpoint,
		IsDefault:                 &defaultArray,
		SkipCertificateValidation: &skipCertificateValidation,
	}
	unityArrays := []model.UnityStorageArray{unityArray}
	unityStorageArrayList := model.UnityStorageArrayList{
		StorageArrayList: unityArrays,
	}

	yamlByte, err := yaml.Marshal(unityStorageArrayList)
	if err != nil {
		panic(err)
	}
	storageArrayConfig := string(yamlByte)

	var valuesArray []string
	valuesArray = append(valuesArray, fmt.Sprintf("arrayconfig=%s", storageArrayConfig))
	return valuesArray, nil
}

func getPowermaxSecret(array model.StorageArray) ([]string, error) {

	decryptedPassword, err := utils.DecryptPassword(array.Password)
	if err != nil {
		return nil, err
	}

	var valuesArray []string
	valuesArray = append(valuesArray, fmt.Sprintf("credentials.username=%s", array.Username))
	valuesArray = append(valuesArray, fmt.Sprintf("credentials.password=%s", string(decryptedPassword)))
	return valuesArray, nil
}

func getPowerstoreSecret(array model.StorageArray) ([]string, error) {

	defaultArray, skipCertificateValidation, err := getBoolValuesFromMetadata(array.MetaData)
	if err != nil {
		return nil, err
	}
	nasName := utils.GetValueFromMetadataKey(array.MetaData, constants.KeyNasName)
	blockProtocol := utils.GetValueFromMetadataKey(array.MetaData, constants.KeyBlockProtocol)
	decryptedPassword, err := utils.DecryptPassword(array.Password)
	if err != nil {
		return nil, err
	}

	powerstoreArray := model.PowerstoreArray{
		GlobalID:                  array.UniqueID,
		Username:                  array.Username,
		Password:                  string(decryptedPassword),
		Endpoint:                  array.ManagementEndpoint,
		IsDefault:                 &defaultArray,
		SkipCertificateValidation: &skipCertificateValidation,
		NasName:                   nasName,
		BlockProtocol:             blockProtocol,
	}
	powerstoreArrays := []model.PowerstoreArray{powerstoreArray}
	powerstoreSecret := model.PowerstoreSecret{
		Arrays: powerstoreArrays,
	}
	yamlByte, err := yaml.Marshal(powerstoreSecret)
	if err != nil {
		return nil, err
	}

	arrayConfig := string(yamlByte)

	var valuesArray []string
	valuesArray = append(valuesArray, fmt.Sprintf("arrayconfig=%s", arrayConfig))
	return valuesArray, nil
}

func getPowerflexSecret(array model.StorageArray) ([]string, error) {

	defaultArray, skipCertificateValidation, err := getBoolValuesFromMetadata(array.MetaData)
	if err != nil {
		return nil, err
	}
	mdm := utils.GetValueFromMetadataKey(array.MetaData, constants.KeyMdmID)
	decryptedPassword, err := utils.DecryptPassword(array.Password)
	if err != nil {
		return nil, err
	}

	powerflexArray := model.PowerflexArray{
		SystemID:                  array.UniqueID,
		Username:                  array.Username,
		Password:                  string(decryptedPassword),
		Endpoint:                  array.ManagementEndpoint,
		IsDefault:                 &defaultArray,
		SkipCertificateValidation: &skipCertificateValidation,
		Mdm:                       mdm,
	}
	powerflexArrays := []model.PowerflexArray{powerflexArray}
	yamlByte, err := yaml.Marshal(powerflexArrays)
	if err != nil {
		return nil, err
	}

	arrayConfig := string(yamlByte)

	var valuesArray []string
	valuesArray = append(valuesArray, fmt.Sprintf("arrayconfig=%s", arrayConfig))
	return valuesArray, nil
}

func getPowerscaleSecret(array model.StorageArray) ([]string, error) {

	defaultArray, skipCertificateValidation, err := getBoolValuesFromMetadata(array.MetaData)
	if err != nil {
		return nil, err
	}
	portString := utils.GetValueFromMetadataKey(array.MetaData, constants.KeyPort)
	port, err := strconv.Atoi(portString)
	if err != nil {
		// Set default isilon port as 8080 when parsing fails
		port = 8080
		err = nil
	}

	decryptedPassword, err := utils.DecryptPassword(array.Password)
	if err != nil {
		return nil, err
	}

	isilonCluster := model.IsilonCluster{
		ClusterName:               array.UniqueID,
		Username:                  array.Username,
		Password:                  string(decryptedPassword),
		Endpoint:                  array.ManagementEndpoint,
		EndpointPort:              uint(port),
		IsDefault:                 &defaultArray,
		SkipCertificateValidation: &skipCertificateValidation,
	}

	isilonClusterArray := []model.IsilonCluster{isilonCluster}
	isilonClusters := model.IsilonClusters{
		IsilonCluster: isilonClusterArray,
	}

	yamlByte, err := yaml.Marshal(isilonClusters)
	if err != nil {
		return nil, err
	}

	arrayConfig := string(yamlByte)

	var valuesArray []string
	valuesArray = append(valuesArray, fmt.Sprintf("arrayconfig=%s", arrayConfig))
	return valuesArray, nil
}

func getBoolValuesFromMetadata(metaData string) (bool, bool, error) {
	var defaultArray bool
	var skipCertificateValidation bool
	var err error
	defaultArrayString := utils.GetValueFromMetadataKey(metaData, constants.KeyIsDefault)
	if defaultArrayString == "" {
		// Set isDefault as false if not specified in metadata
		defaultArray = false
	} else {
		defaultArray, err = strconv.ParseBool(defaultArrayString)
		if err != nil {
			return false, false, fmt.Errorf("Error parsing isDefault: %s ", err.Error())
		}
	}
	skipCertificateValidationString := utils.GetValueFromMetadataKey(metaData, constants.KeySkipCertificateValidation)
	if skipCertificateValidationString == "" {
		// Set isDefault as false if not specified in metadata
		skipCertificateValidation = false
	} else {
		skipCertificateValidation, err = strconv.ParseBool(skipCertificateValidationString)
		if err != nil {
			return false, false, fmt.Errorf("Error parsing skipCertificateValidation: %s ", err.Error())
		}
	}
	return defaultArray, skipCertificateValidation, nil
}
