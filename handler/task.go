package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
)

// GetTask godoc
// @Summary Get an task
// @Description Get an task
// @ID get-task
// @Tags task
// @Accept  json
// @Produce  json
// @Param id path string true "Task ID"
// @Success 200 {object} taskResponse
// @Success 303 {object} taskResponse
// @Failure 400 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTask(c echo.Context) error {
	taskID := c.Param("id")
	task, err := h.taskStore.GetByID(taskID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if task == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	// task requires user intervention, so return it with the appropriate _links to the application
	if task.Status == model.TaskStatusPrompting {

		scheme := utils.GetEnv("SCHEME", "https")
		hostName := utils.GetEnv("HOST", "127.0.0.1")
		port := utils.GetEnv("PORT", "8080")

		// TODO: updating=false is when creating a new application, updating=true is when updating an existing application
		// Need to find a way to determine from the Task if we are creating or updating the application
		approveLink := fmt.Sprintf("%s://%s:%s/api/tasks/%d/approve?updating=false", scheme, hostName, port, task.ApplicationID)
		cancelLink := fmt.Sprintf("%s://%s:%s/api/tasks/%d/cancel?updating=false", scheme, hostName, port, task.ApplicationID)

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

	switch task.ApplicationID {
	case 0:
		return c.JSON(http.StatusOK, newTaskResponse(task))
	default:
		c.Response().Header().Set("Location", fmt.Sprintf("/api/applications/%d", task.ApplicationID))
		return c.JSON(http.StatusSeeOther, newTaskResponse(task))
	}
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
// @Success 202 {object}
// @Failure 400 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /tasks/{id}/approve [post]
func (h *TaskHandler) ApproveStateChange(c echo.Context) error {
	taskID := c.Param("id")
	task, err := h.taskStore.GetByID(taskID)
	if err != nil {
		c.Logger().Errorf("error getting task: %+v", err)
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	// Put the task status back into in-progress.
	task.Status = model.TaskStatusInProgress
	if err := h.taskStore.Update(task); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	GoProcessApplication(h, context.Background(), *task, c)

	c.Response().Header().Set("Location", fmt.Sprintf("/api/tasks/%d", task.ID))
	return c.NoContent(http.StatusAccepted)
}

// GoProcessApplication wrapper to call processApplication as Go routine
var GoProcessApplication = func(h *TaskHandler, ctx context.Context, task model.Task, c echo.Context) {
	go h.processApplication(ctx, task, c)
}

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
// @Success 201 {object}
// @Failure 400 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /tasks/{id}/cancel [post]
func (h *TaskHandler) CancelStateChange(c echo.Context) error {
	taskID := c.Param("id")
	task, err := h.taskStore.GetByID(taskID)
	if err != nil {
		c.Logger().Errorf("error getting task: %+v", err)
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	application, err := h.applicationStore.GetByID(fmt.Sprint(task.ApplicationID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if application == nil {
		c.Logger().Printf("no application with ID %v", task.ApplicationID)
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	applicationStateChange, err := h.applicationStateChangeStore.GetByApplicationID(application.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if applicationStateChange == nil {
		c.Logger().Printf("no application state change found")
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}

	if err := h.applicationStateChangeStore.Delete(applicationStateChange); err != nil {
		c.Logger().Printf("error deleting application state change for application %v", application.ID)
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	// Check the updating query parameter.  If we're not updating (i.e. we're creating), then we can remove
	// the application as well.
	isUpdating := c.QueryParam("updating")
	if isUpdating == "false" {
		application.Name = fmt.Sprintf("_DEL%s_%v", application.Name, time.Now().UnixNano())
		if err := h.applicationStore.Update(application); err != nil {
			c.Logger().Printf("error updating application: %v", application.ID)
			return c.JSON(http.StatusInternalServerError, utils.NewError(err))
		}
		if err := h.applicationStore.Delete(application); err != nil {
			c.Logger().Printf("error deleting application: %v", application.ID)
			return c.JSON(http.StatusInternalServerError, utils.NewError(err))
		}
	}

	task.Status = model.TaskStatusCompleted
	if err := h.taskStore.Update(task); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	return c.NoContent(http.StatusOK)
}

// TODO: This is specific to applications and should probably be moved. Find a way to call the application handler from here.
func (h *TaskHandler) processApplication(ctx context.Context, task model.Task, c echo.Context) {
	time.Sleep(WaitGoProcessApplication)
	c.Logger().Printf("Updating task %d", task.ID)

	// Retrieve the application associated with this task.
	application, err := h.applicationStore.GetByID(fmt.Sprint(task.ApplicationID))
	if err != nil {
		c.Logger().Printf("error getting application: %v", task.ApplicationID)
		c.JSON(http.StatusInternalServerError, utils.NewError(err))
		return
	}
	if application == nil {
		c.Logger().Printf("the application was not found: %v", task.ApplicationID)
		c.JSON(http.StatusNotFound, utils.NotFound())
		return
	}

	// Retrieve the associated application state change.
	applicationStateChange, err := h.applicationStateChangeStore.GetByApplicationID(application.ID)
	if err != nil {
		c.Logger().Printf("error getting the application state change for application: %v", application.ID)
		c.JSON(http.StatusInternalServerError, utils.NewError(err))
		return
	}
	if applicationStateChange == nil {
		c.Logger().Printf("the application state change was not found: %v", application.ID)
		c.JSON(http.StatusNotFound, utils.NotFound())
		return
	}

	// Update the state of the application to reflect the intended state change.
	application.ClusterID = applicationStateChange.ClusterID
	application.ModuleTypes = applicationStateChange.ModuleTypes
	application.DriverTypeID = applicationStateChange.DriverTypeID
	application.DriverConfiguration = applicationStateChange.DriverConfiguration
	application.ModuleConfiguration = applicationStateChange.ModuleConfiguration

	cluster, err := h.clusterStore.GetByID(application.ClusterID)
	if err != nil {
		c.Logger().Errorf("error getting cluster: %+v", err)
		c.JSON(http.StatusInternalServerError, utils.NewError(err))
		return
	}
	configData := cluster.ConfigFileData

	var configFileName string
	if len(configData) != 0 {
		tmpFile, err := ioutil.TempFile("", "config")
		if err != nil {
			c.Logger().Errorf("error creating temp file: %+v", err)
			c.JSON(http.StatusInternalServerError, utils.NewError(err))
			return
		}
		_, err = tmpFile.Write(configData)
		if err != nil {
			c.Logger().Errorf("error writing file: %+v", err)
			c.JSON(http.StatusInternalServerError, utils.NewError(err))
			return
		}
		configFileName = tmpFile.Name()
		defer os.Remove(tmpFile.Name())
	}

	// TODO: not waiting resources to properly come up for now, should be changed later in development
	kappOutput, err := h.kappClient.DeployFromBytes(ctx, applicationStateChange.Template, application.Name, false, configFileName)
	if err != nil {
		c.Logger().Errorf("error deploying app: output = %+v, err = %+v", kappOutput, err)
		task.Status = model.TaskStatusFailed
		if err := h.taskStore.Update(&task); err != nil {
			c.Logger().Errorf("error updating task: %+v", err)
		}
		c.JSON(http.StatusInternalServerError, utils.NewError(err))
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
		c.JSON(http.StatusInternalServerError, utils.NewError(err))
		return
	}

	// Delete the pending application state change.
	if err := h.applicationStateChangeStore.Delete(applicationStateChange); err != nil {
		c.Logger().Printf("error deleting application state change for application %v", application.ID)
		c.JSON(http.StatusInternalServerError, utils.NewError(err))
		return
	}

	task.Status = model.TaskStatusCompleted
	task.ApplicationID = application.ID
	if err := h.taskStore.Update(&task); err != nil {
		log.Printf("error updating task: %+v", err)
	}
	log.Println("Marking task", task.ID, "as finished")
}
