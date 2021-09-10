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
	"net/http"
	"os"
	"time"

	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
)

// GetTask godoc
// @Summary Get a task
// @Description Get a task
// @ID get-task
// @Tags task
// @Accept  json
// @Produce  json
// @Param id path string true "Task ID"
// @Success 200 {object} taskResponse
// @Success 303 {object} taskResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTask(c echo.Context) error {
	taskID := c.Param("id")
	task, err := h.taskStore.GetByID(taskID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}
	if task == nil {
		return c.JSON(http.StatusNotFound, utils.NewErrorResponse(http.StatusNotFound, utils.ErrorSeverity, "", err))
	}

	// task requires user intervention, so return it with the appropriate _links to the application
	if task.Status == model.TaskStatusPrompting {

		scheme := utils.GetEnv("SCHEME", "https")
		hostName := utils.GetEnv("API_SERVER_IP", "127.0.0.1")
		port := utils.GetEnv("API_SERVER_PORT", "31313")

		// TODO: updating=false is when creating a new application, updating=true is when updating an existing application
		// Need to find a way to determine from the Task if we are creating or updating the application
		approveLink := fmt.Sprintf("%s://%s:%s/api/tasks/%d/approve?updating=false", scheme, hostName, port, task.Application.ID)
		cancelLink := fmt.Sprintf("%s://%s:%s/api/tasks/%d/cancel?updating=false", scheme, hostName, port, task.Application.ID)

		links := map[string]map[string]string{
			"yes": {
				"href": approveLink,
			},
			"no": {
				"href": cancelLink,
			},
		}
		return c.JSON(http.StatusOK, newTaskResponseWithLinks(task, links))
	}

	return c.JSON(http.StatusOK, newTaskResponse(task))
}

// ListTasks godoc
// @Summary List all tasks
// @Description List all tasks
// @ID list-tasks
// @Tags task
// @Accept  json
// @Produce  json
// @Param application_name query string false "Application Name"
// @Success 200 {array} taskResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /tasks [get]
func (h *TaskHandler) ListTasks(c echo.Context) error {

	applicationName := c.QueryParam("application_name")
	tasks := make([]model.Task, 0)
	var err error

	if applicationName != "" {
		application, err := h.applicationStore.GetByName(applicationName)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
		}
		if application != nil {
			tasks, err = h.taskStore.GetAllByApplication(application.ID)
		}
	} else {
		tasks, err = h.taskStore.GetAll()
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}

	resp := make([]*taskResponse, 0)
	for _, task := range tasks {
		resp = append(resp, newTaskResponse(&task))
	}
	return c.JSON(http.StatusOK, resp)
}

// ApproveStateChange godoc
// @Summary Approve state change for an application
// @Description Approve state change for an application
// @ID approve-state-change-application
// @Tags task
// @Accept  json
// @Produce  json
// @Param id path string true "Task ID"
// @Param updating query boolean false "Task is associated with an Application update operation"
// @Success 202 {string} string "Accepted"
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /tasks/{id}/approve [post]
func (h *TaskHandler) ApproveStateChange(c echo.Context) error {
	taskID := c.Param("id")
	task, err := h.taskStore.GetByID(taskID)
	if err != nil {
		c.Logger().Errorf("error getting task: %+v", err)
		return c.JSON(http.StatusNotFound, utils.NewErrorResponse(http.StatusNotFound, utils.ErrorSeverity, "", err))

	}

	// Put the task status back into in-progress.
	task.Status = model.TaskStatusInProgress
	if err := h.taskStore.Update(task); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}

	GoProcessApplication(h, context.Background(), *task, c)

	c.Response().Header().Set("Location", fmt.Sprintf("/api/tasks/%d", task.ID))
	return c.NoContent(http.StatusAccepted)
}

// GoProcessApplication wrapper to call processApplication as Go routine
var GoProcessApplication = func(h *TaskHandler, ctx context.Context, task model.Task, c echo.Context) {
	go h.processApplication(ctx, task, c)
}

// WaitGoProcessApplication - Place holder for Process wait time
var WaitGoProcessApplication = 10 * time.Second

// CancelStateChange godoc
// @Summary Cancel state change for an application
// @Description Cancel state change for an application
// @ID cancel-state-change-application
// @Tags task
// @Accept  json
// @Produce  json
// @Param id path string true "Task ID"
// @Param updating query boolean false "Task is associated with an Application update operation"
// @Success 200 {string} string "Success"
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /tasks/{id}/cancel [post]
func (h *TaskHandler) CancelStateChange(c echo.Context) error {
	taskID := c.Param("id")
	task, err := h.taskStore.GetByID(taskID)
	if err != nil {
		c.Logger().Errorf("error getting task: %+v", err)
		return c.JSON(http.StatusNotFound, utils.NewErrorResponse(http.StatusNotFound, utils.ErrorSeverity, "", err))
	}

	application, err := h.applicationStore.GetByID(fmt.Sprint(task.Application.ID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}
	if application == nil {
		c.Logger().Printf("no application with ID %v", task.Application.ID)
		return c.JSON(http.StatusNotFound, utils.NewErrorResponse(http.StatusNotFound, utils.ErrorSeverity, "", err))
	}

	applicationStateChange, err := h.applicationStateChangeStore.GetByApplicationID(application.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}
	if applicationStateChange == nil {
		c.Logger().Printf("no application state change found")
		return c.JSON(http.StatusNotFound, utils.NewErrorResponse(http.StatusNotFound, utils.ErrorSeverity, "", err))
	}

	if err := h.applicationStateChangeStore.Delete(applicationStateChange); err != nil {
		errMgs := fmt.Sprintf("error deleting application state change for application %v", application.ID)
		c.Logger().Print(errMgs)
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, errMgs, err))
	}

	// Check the updating query parameter.  If we're not updating (i.e. we're creating), then we can remove
	// the application as well.
	isUpdating := c.QueryParam("updating")
	if isUpdating == "false" {
		application.Name = fmt.Sprintf("_DEL%s_%v", application.Name, time.Now().UnixNano())
		if err := h.applicationStore.Update(application); err != nil {
			errMgs := fmt.Sprintf("error updating application: %v", application.ID)
			c.Logger().Print(errMgs)
			return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, errMgs, err))
		}
		if err := h.applicationStore.Delete(application); err != nil {
			errMgs := fmt.Sprintf("error deleting application: %v", application.ID)
			c.Logger().Print(errMgs)
			return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, errMgs, err))
		}
	}

	task.Status = model.TaskStatusCompleted
	if err := h.taskStore.Update(task); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}

	return c.NoContent(http.StatusOK)
}

// TODO: This is specific to applications and should probably be moved. Find a way to call the application handler from here.
func (h *TaskHandler) processApplication(ctx context.Context, task model.Task, c echo.Context) {
	time.Sleep(WaitGoProcessApplication)
	c.Logger().Printf("Updating task %d", task.ID)

	// Retrieve the application associated with this task.
	application, err := h.applicationStore.GetByID(fmt.Sprint(task.Application.ID))
	if err != nil {
		errMgs := fmt.Sprintf("error getting application: %v", task.Application.ID)
		c.Logger().Print(errMgs)
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, errMgs, err))
		return
	}
	if application == nil {
		c.Logger().Printf("the application was not found: %v", task.Application.ID)
		c.JSON(http.StatusNotFound, utils.NewErrorResponse(http.StatusNotFound, utils.ErrorSeverity, "", err))
		return
	}

	// Retrieve the associated application state change.
	applicationStateChange, err := h.applicationStateChangeStore.GetByApplicationID(application.ID)
	if err != nil {
		errMgs := fmt.Sprintf("error getting the application state change for application: %v", application.ID)
		c.Logger().Print(errMgs)
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, errMgs, err))
		return
	}
	if applicationStateChange == nil {
		c.Logger().Printf("the application state change was not found: %v", application.ID)
		c.JSON(http.StatusNotFound, utils.NewErrorResponse(http.StatusNotFound, utils.ErrorSeverity, "", err))
		return
	}

	// Update the state of the application to reflect the intended state change.
	application.StorageArrays = applicationStateChange.StorageArrays
	application.ClusterID = applicationStateChange.ClusterID
	application.ModuleTypes = applicationStateChange.ModuleTypes
	application.DriverTypeID = applicationStateChange.DriverTypeID
	application.DriverConfiguration = applicationStateChange.DriverConfiguration
	application.ModuleConfiguration = applicationStateChange.ModuleConfiguration

	cluster, err := h.clusterStore.GetByID(application.ClusterID)
	if err != nil {
		c.Logger().Errorf("error getting cluster: %+v", err)
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "error getting cluster", err))
		return
	}
	configData := cluster.ConfigFileData

	var configFileName string
	if len(configData) != 0 {
		tmpFile, err := ioutil.TempFile("", "config")
		if err != nil {
			c.Logger().Errorf("error creating temp file: %+v", err)
			c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "error creating temp file", err))
			return
		}
		_, err = tmpFile.Write(configData)
		if err != nil {
			c.Logger().Errorf("error writing file: %+v", err)
			c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "error writing file", err))
			return
		}
		configFileName = tmpFile.Name()
		defer os.Remove(tmpFile.Name())
	}

	// TODO: not waiting resources to properly come up for now, should be changed later in development
	kappOutput, err := h.kappClient.DeployFromBytes(ctx, applicationStateChange.Template, application.Name, false, configFileName)
	if err != nil {
		c.Logger().Errorf("error deploying app: output = %+s, err = %+v", kappOutput, err)
		task.Status = model.TaskStatusFailed
		if err := h.taskStore.Update(&task); err != nil {
			c.Logger().Errorf("error updating task: %+v", err)
		}
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, fmt.Sprintf("error deploying app: output = %+v", kappOutput), err))
		return
	}

	// At this point, ytt has compiled our YAML and its been successfully applied
	// to the cluster using kapp deploy.
	// The next step is to finalize the Application resource and put the Task into
	// a Completed state.
	application.ApplicationOutput = kappOutput
	// TODO: The call to create is actually an Upsert operation, so it doesn't
	// matter if the application already exists (it would be saved). Perhaps change
	// this to "Save" instead.
	if err := h.applicationStore.Create(application); err != nil {
		task.Status = model.TaskStatusFailed
		if err := h.taskStore.Update(&task); err != nil {
			c.Logger().Printf("error creating application: %+v", err)
		}
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
		return
	}

	// Delete the pending application state change.
	if err := h.applicationStateChangeStore.Delete(applicationStateChange); err != nil {
		errMgs := fmt.Sprintf("error deleting application state change for application %v", application.ID)
		c.Logger().Print(errMgs)
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, errMgs, err))
		return
	}

	task.Status = model.TaskStatusCompleted
	task.ApplicationID = application.ID
	if err := h.taskStore.Update(&task); err != nil {
		c.Logger().Errorf("error updating task: %+v", err)
	}
	c.Logger().Infof("Marking task %d as finished", task.ID)
}
