// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/dell/csm-deployment/k8s"
	"github.com/dell/csm-deployment/ytt"

	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
)

// CreateApplication - Creates an application
// @Summary Create a new application
// @Description Create a new application
// @ID create-application
// @Tags application
// @Accept  json
// @Produce  json
// @Param application body applicationCreateRequest true "Application info for creation"
// @Success 202 {string} string "Accepted"
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /applications [post]
func (h *ApplicationHandler) CreateApplication(c echo.Context) error {
	var application model.Application
	req := &applicationCreateRequest{}
	if err := req.bind(c, &application, h.ModuleTypeStore); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.CriticalSeverity, "", err))
	}

	convertedStorageArrayIDs, err := convertStringsToUint(req.StorageArrays)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.CriticalSeverity, "", err))
	}
	arrays, err := h.arrayStore.GetAllByID(convertedStorageArrayIDs...)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.CriticalSeverity, "", err))
	}
	application.StorageArrays = arrays

	convertedModuleTypeIDs, err := convertStringsToUint(req.ModuleTypes)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.CriticalSeverity, "", err))
	}
	modules, err := h.ModuleTypeStore.GetAllByID(convertedModuleTypeIDs...)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.CriticalSeverity, "", err))
	}
	application.ModuleTypes = modules

	for _, moduleType := range application.ModuleTypes {
		if moduleType.Name == model.ModuleTypeReplication && application.ModuleConfiguration == "" {
			c.Logger().Errorf("error target cluster is not mentioned in the module configuration: %+v", err)
			return c.JSON(http.StatusUnprocessableEntity, utils.NewErrorResponse(http.StatusUnprocessableEntity,
				utils.CriticalSeverity,
				"error target cluster is not mentioned in the module configuration", err))
		}
	}

	convertedClusterID, err := strconv.Atoi(req.ClusterID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, utils.CriticalSeverity, "", err))
	}

	convertedDriverTypeID, err := strconv.Atoi(req.DriverTypeID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, utils.CriticalSeverity, "", err))
	}

	err = h.precheckHandler.Precheck(c, uint(convertedClusterID), uint(convertedDriverTypeID), modules, application.ModuleConfiguration)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, utils.CriticalSeverity, "", err))
	}

	// Persist the application.  The name must be unique.
	if err := h.applicationStore.Create(&application); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.CriticalSeverity, "", err))
	}

	// Create a task with "in progress" status.
	t := model.Task{
		Status:        model.TaskStatusInProgress,
		ApplicationID: application.ID,
	}
	if err := h.taskStore.Create(&t); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.CriticalSeverity, "", err))
	}

	// get a diff of what we will be installing and request approval
	if !h.SkipGoRoutine {
		go func() {
			time.Sleep(10 * time.Second)
			h.captureApplicationDiff(context.Background(), application.ID, t, c)
		}()
	}
	// TODO: if app was deleted outside of server loop then it would deploy it but fail to add to db

	c.Response().Header().Set("Location", fmt.Sprintf("/api/tasks/%d", t.ID))
	return c.NoContent(http.StatusAccepted)
}

func convertStringsToUint(arr []string) ([]uint, error) {
	result := make([]uint, 0)
	for _, value := range arr {
		converted, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		result = append(result, uint(converted))
	}
	return result, nil
}

// Precheck will run prechecks for drivers and modules
func (p *PrecheckHandler) Precheck(c echo.Context, clusterID uint, driverID uint, modules []model.ModuleType, moduleConfig string) error {
	availableModules := make(map[string]string)

	cluster, err := p.clusterStore.GetByID(clusterID)
	if err != nil {
		return err
	}
	if cluster == nil {
		return fmt.Errorf("not able to find cluster with id %d", clusterID)
	}

	cfs, err := p.configFileStore.GetAll()
	if err != nil {
		return err
	}

	if driverID > 0 {
		driver, err := p.driverStore.GetByID(driverID)
		if err != nil {
			return err
		}
		if driver == nil {
			return fmt.Errorf("not able to find driver with id %d", driverID)
		}

		driverPrechecks := p.precheckGetter.GetDriverPrechecks(driver.StorageArrayType.Name, cluster.ConfigFileData, cluster.ClusterDetails.Nodes, modules, c.Logger())
		for _, precheck := range driverPrechecks {
			c.Logger().Printf("Running precheck: %T for driver of type %s", precheck, driver.StorageArrayType.Name)
			err := precheck.Validate()
			if err != nil {
				return err
			}
		}
		availableModules["csidriver"] = driver.StorageArrayType.Name
	}

	for _, m := range modules {
		availableModules[m.Name] = m.Name
	}

	// Run pre-checks for standalone modules (ex: observability)
	for _, module := range modules {
		modulePrechecks := p.precheckGetter.GetModuleTypePrechecks(module.Name, moduleConfig, cluster.ConfigFileData, cfs, availableModules)
		for _, precheck := range modulePrechecks {
			c.Logger().Printf("Running precheck: %T for %s module", precheck, module.Name)
			err := precheck.Validate()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *ApplicationHandler) captureApplicationDiff(ctx context.Context, applicationID uint, t model.Task, c echo.Context) {
	c.Logger().Infof("Updating task %d", t.ID)
	isReplicationEnabled := false
	isReverseProxyEnabled := false
	reverseProxySecretKeyFilename := ""
	reverseProxySecretCertFilename := ""
	// TODO: Discover the ytt service or use dependency injection.
	h.yttClient.SetOptions(ytt.WithLogger(c.Logger(), false))
	// Retrieve the application.
	application, err := h.applicationStore.GetByID(fmt.Sprintf("%v", applicationID))
	if err != nil {
		c.Logger().Errorf("error getting application: %+v", err)
		return
	}
	var targetConfigData []byte
	var targetFileName string
	sourceConfigData, sourceFileName, done, err := h.getClusterConfigByID(application.ClusterID, c)
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			c.Logger().Errorf("error removing file: %+v", err)
			return
		}
	}(sourceFileName)

	if done {
		return
	}
	// check whether the replication module type is used in the driver installation, if yes we have to connect to the remote cluster.
	for _, moduleType := range application.ModuleTypes {
		module, err := h.ModuleTypeStore.GetByID(moduleType.ID)
		if err != nil {
			c.Logger().Errorf("error getting module: %+v", err)
			return
		}
		if module.Name == model.ModuleTypeReplication {
			isReplicationEnabled = true
			// Fetching target cluster config
			if application.ModuleConfiguration == "" {
				c.Logger().Errorf("Error target cluster is not added in the module configuration: %+v", err)
				return
			}
			tmpArr := strings.Split(application.ModuleConfiguration, "=")
			clusterInInt, err := strconv.ParseUint(tmpArr[1], 10, 64)
			if err != nil {
				c.Logger().Errorf("Error invalid target cluster ID: %+v", err)
				return
			}
			targetConfigData, targetFileName, done, _ = h.getClusterConfigByID(uint(clusterInInt), c)
			if done {
				return
			}
			defer func(name string) {
				err := os.Remove(name)
				if err != nil {
					c.Logger().Errorf("error removing file: %+v", err)
					return
				}
			}(targetFileName)
			// we should create resources only for multi-cluster scenario
			if application.ClusterID != uint(clusterInInt) {
				h.createReplicationNamespace(targetConfigData, c)
				// creating target secret in the source cluster
				h.createReplicationSecrets(sourceConfigData, targetConfigData, c)
			}
		}
		if module.Name == model.ModuleTypeReverseProxy {
			isReverseProxyEnabled = true
			moduleConfigParams := strings.Split(application.ModuleConfiguration, " ")
			reverseProxyKeyFound := false
			reverseProxyCertFound := false
			for _, param := range moduleConfigParams {
				paramValue := strings.Split(strings.Trim(param, "\""), "=")
				key := paramValue[0]
				if key == "reverseProxy.tlsSecretKeyFile" {
					reverseProxyKeyFound = true
					reverseProxySecretKeyFilename = paramValue[1]
				}
				if key == "reverseProxy.tlsSecretCertFile" {
					reverseProxyCertFound = true
					reverseProxySecretCertFilename = paramValue[1]
				}
			}
			if !reverseProxyKeyFound {
				c.Logger().Errorf("reverse proxy secret key file path missing in module configuration")
				return
			}
			if !reverseProxyCertFound {
				c.Logger().Errorf("reverse proxy secret cert file path missing in module configuration")
				return
			}
		}
	}

	applicationStateChange := &model.ApplicationStateChange{
		ApplicationID:       application.ID,
		ClusterID:           application.ClusterID,
		DriverTypeID:        application.DriverTypeID,
		ModuleTypes:         application.ModuleTypes,
		StorageArrays:       application.StorageArrays,
		DriverConfiguration: application.DriverConfiguration,
		ModuleConfiguration: application.ModuleConfiguration,
	}
	// store the requested state change for this application
	if err := h.applicationStateChangeStore.Create(applicationStateChange); err != nil {
		t.Status = model.TaskStatusFailed
		if err := h.taskStore.Update(&t); err != nil {
			c.Logger().Errorf("error creating application state change: %+v", err)
		}
		return
	}

	// Create the namespace manifest
	namespaceOutput, err := h.yttClient.NamespaceTemplateFromApplication(applicationStateChange.ID, h.applicationStateChangeStore)
	if err != nil {
		c.Logger().Errorf("error generating namespace manifest: %+v", err)
		return
	}

	// First create the secret manifest
	secretOutput, err := h.yttClient.GenerateDynamicSecret(applicationStateChange.ID, h.applicationStateChangeStore, h.configFileStore)
	if err != nil {
		c.Logger().Errorf("error generating secret: %+v", err)
		return
	}

	configMapOutput, err := h.yttClient.ConfigMapTemplateFromApplication(applicationStateChange.ID, h.applicationStateChangeStore)
	if err != nil {
		c.Logger().Errorf("error generating secret: %+v", err)
		return
	}

	k8sClient, err := h.runtimeClientFunc(sourceConfigData, c.Logger())
	if err != nil {
		c.Logger().Errorf("error getting k8s client: %+v", err)
		return
	}

	err = k8sClient.CreateNameSpace(ctx, namespaceOutput.AsCombinedBytes())
	if err != nil {
		c.Logger().Errorf("error creating namespace: %+v", err)
		return
	}

	// we will create the replication resources in the source cluster also.
	if isReplicationEnabled {
		h.createReplicationNamespace(sourceConfigData, c)
		// creating source secret in the target cluster
		h.createReplicationSecrets(targetConfigData, sourceConfigData, c)
		err := CreateReplicationController(sourceFileName, targetFileName, c)
		if err != nil {
			c.Logger().Errorf("error creating replication controller: %+v", err)
			return
		}
	}

	if isReverseProxyEnabled {
		// create reverse proxy tls secrets
		h.createReverseProxySecrets(reverseProxySecretKeyFilename, reverseProxySecretCertFilename, sourceConfigData, namespaceOutput.AsCombinedBytes(), c)
	}

	secrets, err := utils.SplitYAML(secretOutput.AsCombinedBytes())
	if err != nil {
		c.Logger().Errorf("error spliting secret yaml: %+v", err)
		return
	}
	for _, secret := range secrets {
		err = k8sClient.CreateSecret(ctx, secret)
		if err != nil {
			c.Logger().Errorf("error creating secret: %+v", err)
			return
		}
	}

	err = k8sClient.CreateConfigMap(ctx, configMapOutput.AsCombinedBytes())
	if err != nil {
		c.Logger().Errorf("error creating configmap: %+v", err)
		return
	}

	// For Unity and PowerScale an empty secret has to be created
	arrayType := applicationStateChange.StorageArrays[0].StorageArrayType.Name
	if arrayType == model.ArrayTypeUnity || arrayType == model.ArrayTypePowerScale {
		// Create the empty secret manifest
		emptySecretOutput, err := h.yttClient.GetEmptySecret(applicationStateChange.ID, h.applicationStateChangeStore)
		if err != nil {
			c.Logger().Errorf("error generating empty secret: %+v", err)
			return
		}
		err = k8sClient.CreateSecret(ctx, emptySecretOutput.AsCombinedBytes())
		if err != nil {
			c.Logger().Errorf("error creating secret: %+v", err)
			return
		}
	}

	output, err := h.yttClient.TemplateFromApplication(applicationStateChange.ID, h.applicationStateChangeStore, h.clusterStore, h.configFileStore)
	if err != nil {
		c.Logger().Errorf("error generating app: %+v", err)
		return
	}

	applicationStateChange.Template = output.AsCombinedBytes()
	if err := h.applicationStateChangeStore.Create(applicationStateChange); err != nil {
		t.Status = model.TaskStatusFailed
		if err := h.taskStore.Update(&t); err != nil {
			c.Logger().Errorf("error updating application state change: %+v", err)
		}
		return
	}

	kappOutput, err := h.kappClient.GetDeployDiff(ctx, output.AsCombinedBytes(), application.Name, sourceFileName)
	if err != nil {
		c.Logger().Errorf("error deploying app: output = %+s, err = %+v", kappOutput, err)
		t.Status = model.TaskStatusFailed
		if err := h.taskStore.Update(&t); err != nil {
			c.Logger().Errorf("error updating task: %+v", err)
		}
		return
	}

	t.Logs = []byte(kappOutput)
	t.Status = model.TaskStatusPrompting
	t.ApplicationID = application.ID
	if err := h.taskStore.Update(&t); err != nil {
		c.Logger().Errorf("error updating task: %+v", err)
	}
	log.Println("Marking task", t.ID, "as prompting")
}

func (h *ApplicationHandler) getClusterConfigByID(clusterID uint, c echo.Context) ([]byte, string, bool, error) {
	// Retrieve the cluster associated with the application.
	cluster, err := h.clusterStore.GetByID(clusterID)
	if err != nil {
		c.Logger().Errorf("error getting cluster: %+v", err)
		return nil, "", true, nil
	}
	configData := cluster.ConfigFileData

	if len(configData) == 0 {
		c.Logger().Errorf("no config data for cluster %v", cluster.ID)
		return nil, "", true, nil
	}

	tmpFile, err := ioutil.TempFile("", "config")
	if err != nil {
		c.Logger().Errorf("error creating temp file: %+v", err)
		return nil, "", true, nil
	}
	_, err = tmpFile.Write(configData)
	if err != nil {
		c.Logger().Errorf("error writing file: %+v", err)
		return nil, "", true, nil
	}
	return configData, tmpFile.Name(), false, err
}

// ListApplications - List all applications
// @Summary List all applications
// @Description List all applications
// @ID list-applications
// @Tags application
// @Accept  json
// @Produce  json
// @Param name query string false "Application Name"
// @Success 200 {array} applicationResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /applications [get]
func (h *ApplicationHandler) ListApplications(c echo.Context) error {
	name := c.QueryParam("name")
	applications := make([]model.Application, 0)
	var err error
	if name != "" {
		application, err := h.applicationStore.GetByName(name)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
		}
		if application != nil {
			applications = append(applications, *application)
		}
	} else {
		applications, err = h.applicationStore.GetAll()
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}

	resp := make([]*applicationResponse, 0)
	for _, application := range applications {
		resp = append(resp, newApplicationResponse(&application))
	}
	return c.JSON(http.StatusOK, resp)
}

// GetApplication - Retrieves application by ID
// @Summary Get an application
// @Description Get an application
// @ID get-application
// @Tags application
// @Accept  json
// @Produce  json
// @Param id path string true "Application ID"
// @Success 200 {object} applicationResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /applications/{id} [get]
func (h *ApplicationHandler) GetApplication(c echo.Context) error {
	applicationID := c.Param("id")
	application, err := h.applicationStore.GetByID(applicationID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.CriticalSeverity, "", err))
	}
	if application == nil {
		return c.JSON(http.StatusNotFound, utils.NewErrorResponse(http.StatusNotFound, utils.ErrorSeverity, "", err))
	}
	return c.JSON(http.StatusOK, newApplicationResponse(application))
}

// DeleteApplication - Deletes an application
// @Summary Delete an application
// @Description Delete an application
// @ID delete-application
// @Tags application
// @Accept  json
// @Produce  json
// @Param id path string true "Application ID"
// @Success 204
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /applications/{id} [delete]
func (h *ApplicationHandler) DeleteApplication(c echo.Context) error {
	applicationID := c.Param("id")
	application, err := h.applicationStore.GetByID(applicationID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.CriticalSeverity, "", err))
	}
	if application == nil {
		return c.JSON(http.StatusNotFound, utils.NewErrorResponse(http.StatusNotFound, utils.ErrorSeverity, "", err))
	}

	// TODO: shouldn't delete request be async and create task too?
	cluster, err := h.clusterStore.GetByID(application.ClusterID)
	if err != nil {
		c.Logger().Errorf("error getting cluster: %+v", err)
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.CriticalSeverity, "", err))
	}
	configData := cluster.ConfigFileData

	var configFileName string
	if len(configData) != 0 {
		tmpFile, err := ioutil.TempFile("", "config")
		if err != nil {
			c.Logger().Errorf("error creating temp file: %+v", err)
			return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.CriticalSeverity, "", err))
		}
		_, err = tmpFile.Write(configData)
		if err != nil {
			c.Logger().Errorf("error writing file: %+v", err)
			return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.CriticalSeverity, "", err))
		}
		configFileName = tmpFile.Name()
		defer func(name string) {
			err := os.Remove(name)
			if err != nil {
				c.Logger().Errorf("error removing file: %+v", err)
			}
		}(tmpFile.Name())
	}

	_, err = h.kappClient.Delete(c.Request().Context(), application.Name, configFileName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.CriticalSeverity, "", err))
	}

	if err := h.applicationStore.Delete(application); err != nil {
		c.Logger().Errorf("error deleting application: %+v", err)
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.CriticalSeverity, "", err))
	}

	// Deleting source cluster namespace - replication
	sourceClient, err := h.runtimeClientFunc(configData, c.Logger())
	if err != nil {
		c.Logger().Errorf("error getting k8s client: %+v", err)
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.CriticalSeverity, "", err))
	}
	ctx := context.Background()
	err = sourceClient.DeleteNameSpaceByName(ctx, model.ReplicationNamespace)
	if err != nil {
		c.Logger().Errorf("error deleting namespace: %+v", err)
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.CriticalSeverity, "", err))
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *ApplicationHandler) createReplicationNamespace(configData []byte, c echo.Context) {
	k8sClient, err := h.runtimeClientFunc(configData, c.Logger())
	if err != nil {
		c.Logger().Errorf("error getting k8s client: %+v", err)
		return
	}
	ctx := context.Background()
	err = k8sClient.CreateNameSpaceFromName(ctx, model.ReplicationNamespace)
	if err != nil {
		c.Logger().Errorf("error creating namespace: %+v", err)
		return
	}
}

// GetRuntimeClient will return a k8s runtime client
var GetRuntimeClient = func(data []byte, logger echo.Logger) (k8s.ControllerRuntimeInterface, error) {
	return k8s.NewControllerRuntimeClient(data, logger)
}

func (h *ApplicationHandler) createReplicationSecrets(sourceConfigData []byte, targetConfigData []byte, c echo.Context) {
	k8sClient, err := h.runtimeClientFunc(sourceConfigData, c.Logger())
	if err != nil {
		c.Logger().Errorf("error getting k8s client: %+v", err)
		return
	}
	ctx := context.Background()
	err = k8sClient.CreateNameSpaceFromName(ctx, model.ReplicationNamespace)
	if err != nil {
		c.Logger().Errorf("error creating namespace: %+v", err)
		return
	}

	err = k8sClient.CreateSecretFromName(ctx, "replication-creds", "dell-replication-controller", targetConfigData)
	if err != nil {
		c.Logger().Errorf("error creating secret: %+v", err)
		return
	}
}

func (h *ApplicationHandler) createReverseProxySecrets(key, cert string, sourceConfigData, namespaceData []byte, c echo.Context) {
	k8sClient, err := h.runtimeClientFunc(sourceConfigData, c.Logger())
	if err != nil {
		c.Logger().Errorf("error getting k8s client: %+v", err)
		return
	}
	ctx := context.Background()

	keyBytes, err := h.configFileStore.GetAllByName(key)
	if err != nil {
		c.Logger().Errorf("error loading tls secret key: %+v", err)
		return
	}

	certBytes, err := h.configFileStore.GetAllByName(cert)
	if err != nil {
		c.Logger().Errorf("error loading tls secret cert: %+v", err)
		return
	}

	err = k8sClient.CreateTLSSecretFromName(ctx, "csirevproxy-tls-secret", namespaceData, keyBytes[0].ConfigFileData, certBytes[0].ConfigFileData)
	if err != nil {
		c.Logger().Errorf("error creating secret: %+v", err)
		return
	}
}

// CreateReplicationController creates the replication controller
var CreateReplicationController = func(sourceConfig string, targetConfig string, c echo.Context) error {
	baseURL, err := os.Getwd()
	if err != nil {
		return err
	}
	pathSeparator := string(os.PathSeparator)
	repctlBinaryName := "repctl"
	if runtime.GOOS == "windows" {
		repctlBinaryName = repctlBinaryName + ".exe"
	}
	repctlBinary := baseURL + pathSeparator + "templates" + pathSeparator + repctlBinaryName
	cmd := exec.CommandContext(context.Background(), repctlBinary, "cluster", "add", "-f", sourceConfig+","+targetConfig, "-n", "source,target")
	out, err := cmd.CombinedOutput()
	if err != nil {
		c.Logger().Infof(fmt.Sprintf("%s", out))
		return err
	}
	c.Logger().Infof(fmt.Sprintf("%s", out))

	// To create CRDs
	crdPath := baseURL + pathSeparator + "templates" + pathSeparator + "replication-controller" + pathSeparator + "replicationcrds.all.yaml"
	crdCmd := exec.CommandContext(context.Background(), repctlBinary, "create", "-f", crdPath)
	out, err = crdCmd.CombinedOutput()
	if err != nil {
		c.Logger().Infof(fmt.Sprintf("%s", out))
		return err
	}
	c.Logger().Infof(fmt.Sprintf("%s", out))

	// To create controller
	controllerTemplatePath := baseURL + pathSeparator + "templates" + pathSeparator + "replication-controller" + pathSeparator + "controller.yaml"
	ctrlCmd := exec.CommandContext(context.Background(), repctlBinary, "create", "-f", controllerTemplatePath)
	out, err = ctrlCmd.CombinedOutput()
	if err != nil {
		c.Logger().Infof(fmt.Sprintf("%s", out))
		return err
	}
	c.Logger().Infof(fmt.Sprintf("%s", out))

	// To inject configs to the clusters
	configsCmd := exec.CommandContext(context.Background(), repctlBinary, "cluster", "inject")
	out, err = configsCmd.CombinedOutput()
	if err != nil {
		c.Logger().Infof(fmt.Sprintf("%s", out))
		return err
	}
	c.Logger().Infof(fmt.Sprintf("%s", out))

	return err
}
