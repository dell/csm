package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dell/csm-deployment/k8s"

	"github.com/dell/csm-deployment/ytt"

	"github.com/dell/csm-deployment/kapp"
	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
)

// CreateApplication godoc
// @Summary Create a new application
// @Description Create a new application
// @ID create-application
// @Tags application
// @Accept  json
// @Produce  json
// @Param application body applicationCreateRequest true "Application info for creation"
// @Success 202 {object} applicationResponse
// @Failure 400 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /applications [post]
func (h *ApplicationHandler) CreateApplication(c echo.Context) error {
	var application model.Application
	req := &applicationCreateRequest{}
	if err := req.bind(c, &application, h.moduleStore); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	arrays, err := h.arrayStore.GetAllByID(req.Application.StorageArrays...)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	application.StorageArrays = arrays

	modules, err := h.moduleStore.GetAllByID(req.Application.ModuleTypes...)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	application.ModuleTypes = modules

	// Persist the application.  The name must be unique.
	if err := h.applicationStore.Create(&application); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	// Create a task with "in progress" status.
	t := model.Task{
		Status: model.TaskStatusInProgress,
	}
	if err := h.taskStore.Create(&t); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	// get a diff of what we will be installing and request approval
	go h.captureApplicationDiff(context.Background(), application.ID, t, c)

	// TODO: if app was deleted outside of server loop then it would deploy it but fail to add to db

	c.Response().Header().Set("Location", fmt.Sprintf("/api/tasks/%d", t.ID))
	return c.NoContent(http.StatusAccepted)
}

func (h *ApplicationHandler) captureApplicationDiff(ctx context.Context, applicationID uint, t model.Task, c echo.Context) {
	time.Sleep(10 * time.Second)
	c.Logger().Printf("Updating task %d", t.ID)

	// Retrieve the application.
	application, err := h.applicationStore.GetByID(fmt.Sprintf("%v", applicationID))
	if err != nil {
		c.Logger().Errorf("error getting application: %+v", err)
		return
	}

	// Retrieve the cluster associated with the application.
	cluster, err := h.clusterStore.GetByID(application.ClusterID)
	if err != nil {
		c.Logger().Errorf("error getting cluster: %+v", err)
		return
	}
	configData := cluster.ConfigFileData

	var configFileName string
	if len(configData) == 0 {
		c.Logger().Errorf("no config data for cluster %v", cluster.ID)
		return
	}

	tmpFile, err := ioutil.TempFile("", "config")
	if err != nil {
		c.Logger().Errorf("error creating temp file: %+v", err)
		return
	}
	_, err = tmpFile.Write(configData)
	if err != nil {
		c.Logger().Errorf("error writing file: %+v", err)
		return
	}
	configFileName = tmpFile.Name()
	defer os.Remove(tmpFile.Name())

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
			log.Printf("error creating application state change: %+v", err)
		}
		return
	}

	// TODO: Discover the ytt service or use dependency injection.
	yttClient := ytt.NewClient(ytt.WithLogger(c.Logger(), false))

	// Create the namespace manifest
	namespaceOutput, err := yttClient.NamespaceTemplateFromApplication(applicationStateChange.ID, h.applicationStateChangeStore)
	if err != nil {
		c.Logger().Errorf("error generating namespace manifest: %+v", err)
		return
	}

	// First create the secret manifest
	secretOutput, err := yttClient.SecretTemplateFromApplication(applicationStateChange.ID, h.applicationStateChangeStore)
	if err != nil {
		c.Logger().Errorf("error generating secret: %+v", err)
		return
	}
	k8sClient, err := k8s.NewControllerRuntimeClient(configData, c.Logger())
	if err != nil {
		c.Logger().Errorf("error getting k8s client: %+v", err)
		return
	}

	err = k8sClient.CreateNameSpace(ctx, namespaceOutput.AsCombinedBytes())
	if err != nil {
		c.Logger().Errorf("error creating namespace: %+v", err)
		return
	}

	err = k8sClient.DeployFromBytes(ctx, secretOutput.AsCombinedBytes())
	if err != nil {
		c.Logger().Errorf("error creating secret: %+v", err)
		return
	}

	output, err := yttClient.TemplateFromApplication(applicationStateChange.ID, h.applicationStateChangeStore, h.clusterStore)
	if err != nil {
		c.Logger().Errorf("error generating app: %+v", err)
		return
	}

	applicationStateChange.Template = output.AsCombinedBytes()
	if err := h.applicationStateChangeStore.Create(applicationStateChange); err != nil {
		t.Status = model.TaskStatusFailed
		if err := h.taskStore.Update(&t); err != nil {
			log.Printf("error updating application state change: %+v", err)
		}
		return
	}

	client := kapp.NewClient("")
	kappOutput, err := client.GetDeployDiff(ctx, output.AsCombinedBytes(), application.Name, configFileName)
	if err != nil {
		c.Logger().Errorf("error deploying app: %+v", err)
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
		log.Printf("error updating task: %+v", err)
	}
	log.Println("Marking task", t.ID, "as prompting")
}

// GetApplication godoc
// @Summary Get an application
// @Description Get an application
// @ID get-application
// @Tags application
// @Accept  json
// @Produce  json
// @Param id path string true "Application ID"
// @Success 200 {object} applicationResponse
// @Failure 400 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /applications/{id} [get]
func (h *ApplicationHandler) GetApplication(c echo.Context) error {
	applicationID := c.Param("id")
	application, err := h.applicationStore.GetByID(applicationID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if application == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	return c.JSON(http.StatusOK, newApplicationResponse(application))
}

// UpdateApplication godoc
// @Summary Update an application
// @Description Update an application
// @ID update-application
// @Tags application
// @Accept  json
// @Produce  json
// @Param id path string true "Application ID"
// @Success 200 {object} applicationResponse
// @Failure 400 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /applications/{id} [put]
func (h *ApplicationHandler) UpdateApplication(c echo.Context) error {
	applicationID := c.Param("id")
	application, err := h.applicationStore.GetByID(applicationID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if application == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	// Create a task with "in progress" status.
	t := model.Task{
		Status: model.TaskStatusInProgress,
	}
	if err := h.taskStore.Create(&t); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	// Capture the diff and add it to the task which will need to be approved/cancelled to proceed with application update
	go h.captureApplicationDiff(context.Background(), application.ID, t, c)

	c.Response().Header().Set("Location", fmt.Sprintf("/api/tasks/%d", t.ID))
	return c.NoContent(http.StatusAccepted)
}

// DeleteApplication godoc
// @Summary Delete an application
// @Description Delete an application
// @ID delete-application
// @Tags application
// @Accept  json
// @Produce  json
// @Param id path string true "Application ID"
// @Success 200 {object} applicationResponse
// @Failure 400 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /applications/{id} [delete]
func (h *ApplicationHandler) DeleteApplication(c echo.Context) error {
	applicationID := c.Param("id")
	application, err := h.applicationStore.GetByID(applicationID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if application == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	// TODO: shouldn't delete request be async and create task too?
	cluster, err := h.clusterStore.GetByID(application.ClusterID)
	if err != nil {
		c.Logger().Errorf("error getting cluster: %+v", err)
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	configData := cluster.ConfigFileData

	var configFileName string
	if len(configData) != 0 {
		tmpFile, err := ioutil.TempFile("", "config")
		if err != nil {
			c.Logger().Errorf("error creating temp file: %+v", err)
			return c.JSON(http.StatusInternalServerError, utils.NewError(err))
		}
		_, err = tmpFile.Write(configData)
		if err != nil {
			c.Logger().Errorf("error writing file: %+v", err)
			return c.JSON(http.StatusInternalServerError, utils.NewError(err))
		}
		configFileName = tmpFile.Name()
		defer os.Remove(tmpFile.Name())
	}

	client := kapp.NewClient("")
	kappOutput, err := client.Delete(c.Request().Context(), application.Name, configFileName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	if err := h.applicationStore.Delete(application); err != nil {
		c.Logger().Errorf("error deleting application: %+v", err)
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.JSON(http.StatusOK, kappOutput)
}
