package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	kappMocks "github.com/dell/csm-deployment/kapp/mocks"
	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/router"
	storeMocks "github.com/dell/csm-deployment/store/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const (
	ApplicationID = uint(456)
	TaskID        = uint(123)
	ClusterID     = uint(789)
)

func Test_GetTask(t *testing.T) {
	tests := map[string]func(t *testing.T) (int, *TaskHandler, string, *gomock.Controller){
		"success": func(*testing.T) (int, *TaskHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			getTaskResponseJSON := `{"id":123,"status":"task-status","application_id":0,"logs":"task-logs","_links":null}`

			taskStore := storeMocks.NewMockTaskStoreInterface(ctrl)
			task := model.Task{
				Status:        "task-status",
				ApplicationID: 0,
				Logs:          []byte("task-logs"),
			}
			task.ID = TaskID
			taskStore.EXPECT().GetByID(fmt.Sprint(task.ID)).Times(1).Return(&task, nil)

			handler := &TaskHandler{taskStore: taskStore}

			return http.StatusOK, handler, getTaskResponseJSON, ctrl
		},

		"success SeeOther": func(*testing.T) (int, *TaskHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			getTaskResponseJSON := `{"id":123,"status":"task-status","application_id":456,"logs":"task-logs","_links":null}`

			taskStore := storeMocks.NewMockTaskStoreInterface(ctrl)
			task := model.Task{
				Status:        "task-status",
				ApplicationID: ApplicationID,
				Logs:          []byte("task-logs"),
			}
			task.ID = TaskID
			taskStore.EXPECT().GetByID(fmt.Sprint(task.ID)).Times(1).Return(&task, nil)

			handler := &TaskHandler{taskStore: taskStore}

			return http.StatusSeeOther, handler, getTaskResponseJSON, ctrl
		},

		"success prompting": func(*testing.T) (int, *TaskHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			getTaskResponseJSON := `{"id":123,"status":"prompting","application_id":456,"logs":"task-logs","_links":{"no":{"href":"https://127.0.0.1:8080/api/tasks/456/cancel?updating=false"},"yes":{"href":"https://127.0.0.1:8080/api/tasks/456/approve?updating=false"}}}`

			taskStore := storeMocks.NewMockTaskStoreInterface(ctrl)
			task := model.Task{
				Status:        "prompting",
				ApplicationID: ApplicationID,
				Logs:          []byte("task-logs"),
			}
			task.ID = TaskID
			taskStore.EXPECT().GetByID(fmt.Sprint(task.ID)).Times(1).Return(&task, nil)

			handler := &TaskHandler{taskStore: taskStore}

			return http.StatusOK, handler, getTaskResponseJSON, ctrl
		},

		"error getting task": func(*testing.T) (int, *TaskHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			getTaskResponseJSON := `{"http_status_code":500,"messages":[{"code":500,"message":"Internal Server Error","message_l10n":"hello-world error","Arguments":null,"severity":"CRITICAL"}]}`

			taskStore := storeMocks.NewMockTaskStoreInterface(ctrl)
			taskStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, errors.New("hello-world error"))

			handler := &TaskHandler{taskStore: taskStore}

			return http.StatusInternalServerError, handler, getTaskResponseJSON, ctrl
		},

		"task not found": func(*testing.T) (int, *TaskHandler, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			getTaskResponseJSON := ``

			taskStore := storeMocks.NewMockTaskStoreInterface(ctrl)
			taskStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, nil)

			handler := &TaskHandler{taskStore: taskStore}

			return http.StatusNotFound, handler, getTaskResponseJSON, ctrl
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			expectedStatus, handler, expectedResponse, ctrl := tc(t)

			e := router.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/tasks/:id")
			c.SetParamNames("id")
			c.SetParamValues(fmt.Sprint((TaskID)))

			assert.NoError(t, handler.GetTask(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}

func Test_ApproveStateChange(t *testing.T) {

	tests := map[string]func(t *testing.T) (int, *TaskHandler, string, string, *gomock.Controller){
		"success": func(*testing.T) (int, *TaskHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			getTaskResponseJSON := ""
			taskStore := storeMocks.NewMockTaskStoreInterface(ctrl)
			task := model.Task{
				Status:        "task-status",
				ApplicationID: ApplicationID,
				Logs:          []byte("task-logs"),
			}
			task.ID = TaskID
			taskStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&task, nil)
			taskStore.EXPECT().Update(gomock.Any()).Times(2).Return(nil)

			applicationStore := storeMocks.NewMockApplicationStoreInterface(ctrl)
			application := model.Application{}
			application.ID = ApplicationID
			applicationStore.EXPECT().GetByID(fmt.Sprint(application.ID)).Times(1).Return(&application, nil)
			applicationStore.EXPECT().Create(gomock.Any()).Times(1).Return(nil)

			applicationStateChangeStore := storeMocks.NewMockApplicationStateChangeStoreInterface(ctrl)
			applicationStateChange := model.ApplicationStateChange{}
			applicationStateChange.ID = ApplicationID
			applicationStateChangeStore.EXPECT().GetByApplicationID(gomock.Any()).Times(1).Return(&applicationStateChange, nil)
			applicationStateChangeStore.EXPECT().Delete(gomock.Any()).Times(1).Return(nil)

			clusterStore := storeMocks.NewMockClusterStoreInterface(ctrl)
			cluster := model.Cluster{}
			cluster.ID = ClusterID
			cluster.ConfigFileData = []byte("test-config-data")
			clusterStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&cluster, nil)

			kappClient := kappMocks.NewMockInterface(ctrl)
			kappClient.EXPECT().DeployFromBytes(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return("testing kapp output", nil)

			handler := &TaskHandler{
				taskStore:                   taskStore,
				applicationStore:            applicationStore,
				applicationStateChangeStore: applicationStateChangeStore,
				clusterStore:                clusterStore,
				kappClient:                  kappClient,
			}

			return http.StatusAccepted, handler, "123", getTaskResponseJSON, ctrl
		},

		"task not found": func(*testing.T) (int, *TaskHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			getTaskResponseJSON := ``

			taskStore := storeMocks.NewMockTaskStoreInterface(ctrl)
			taskStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, errors.New("hello-world error"))

			handler := &TaskHandler{taskStore: taskStore}

			return http.StatusNotFound, handler, "123", getTaskResponseJSON, ctrl
		},
		"error deploying form bytes": func(*testing.T) (int, *TaskHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			getTaskResponseJSON := ``
			taskStore := storeMocks.NewMockTaskStoreInterface(ctrl)
			task := model.Task{
				Status:        "task-status",
				ApplicationID: ApplicationID,
				Logs:          []byte("task-logs"),
			}
			task.ID = TaskID
			taskStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&task, nil)
			taskStore.EXPECT().Update(gomock.Any()).Times(2).Return(nil)

			applicationStore := storeMocks.NewMockApplicationStoreInterface(ctrl)
			application := model.Application{}
			application.ID = ApplicationID
			applicationStore.EXPECT().GetByID(fmt.Sprint(application.ID)).Times(1).Return(&application, nil)

			applicationStateChangeStore := storeMocks.NewMockApplicationStateChangeStoreInterface(ctrl)
			applicationStateChange := model.ApplicationStateChange{}
			applicationStateChange.ID = ApplicationID
			applicationStateChangeStore.EXPECT().GetByApplicationID(gomock.Any()).Times(1).Return(&applicationStateChange, nil)

			clusterStore := storeMocks.NewMockClusterStoreInterface(ctrl)
			cluster := model.Cluster{}
			cluster.ID = ClusterID
			cluster.ConfigFileData = []byte("test-config-data")
			clusterStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&cluster, nil)

			kappClient := kappMocks.NewMockInterface(ctrl)
			kappClient.EXPECT().DeployFromBytes(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return("", errors.New("hello-world error"))

			handler := &TaskHandler{
				taskStore:                   taskStore,
				applicationStore:            applicationStore,
				applicationStateChangeStore: applicationStateChangeStore,
				clusterStore:                clusterStore,
				kappClient:                  kappClient,
			}

			return http.StatusInternalServerError, handler, "123", getTaskResponseJSON, ctrl
		},
		"error saving/creating deployment": func(*testing.T) (int, *TaskHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			getTaskResponseJSON := ``
			taskStore := storeMocks.NewMockTaskStoreInterface(ctrl)
			task := model.Task{
				Status:        "task-status",
				ApplicationID: ApplicationID,
				Logs:          []byte("task-logs"),
			}
			task.ID = TaskID
			taskStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&task, nil)
			taskStore.EXPECT().Update(gomock.Any()).Times(2).Return(nil)

			applicationStore := storeMocks.NewMockApplicationStoreInterface(ctrl)
			application := model.Application{}
			application.ID = ApplicationID
			applicationStore.EXPECT().GetByID(fmt.Sprint(application.ID)).Times(1).Return(&application, nil)
			applicationStore.EXPECT().Create(gomock.Any()).Times(1).Return(errors.New("hello-world error"))

			applicationStateChangeStore := storeMocks.NewMockApplicationStateChangeStoreInterface(ctrl)
			applicationStateChange := model.ApplicationStateChange{}
			applicationStateChange.ID = ApplicationID
			applicationStateChangeStore.EXPECT().GetByApplicationID(gomock.Any()).Times(1).Return(&applicationStateChange, nil)

			clusterStore := storeMocks.NewMockClusterStoreInterface(ctrl)
			cluster := model.Cluster{}
			cluster.ID = ClusterID
			cluster.ConfigFileData = []byte("test-config-data")
			clusterStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&cluster, nil)

			kappClient := kappMocks.NewMockInterface(ctrl)
			kappClient.EXPECT().DeployFromBytes(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return("", nil)

			handler := &TaskHandler{
				taskStore:                   taskStore,
				applicationStore:            applicationStore,
				applicationStateChangeStore: applicationStateChangeStore,
				clusterStore:                clusterStore,
				kappClient:                  kappClient,
			}

			return http.StatusInternalServerError, handler, "123", getTaskResponseJSON, ctrl
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			oldGoProcessApplication := GoProcessApplication
			oldWaitGoProcessApplication := WaitGoProcessApplication
			defer func() {
				fmt.Println("setting")
				GoProcessApplication = oldGoProcessApplication
				WaitGoProcessApplication = oldWaitGoProcessApplication
			}()

			WaitGoProcessApplication = 1 * time.Second
			GoProcessApplication = func(h *TaskHandler, ctx context.Context, task model.Task, c echo.Context) {
				var wg sync.WaitGroup

				wg.Add(1)
				go func() {
					defer wg.Done()
					h.processApplication(ctx, task, c)
				}()
				wg.Wait()
			}

			expectedStatus, handler, taskID, expectedResponse, ctrl := tc(t)

			e := router.New()
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/tasks/:id/approve")
			c.SetParamNames("id")
			c.SetParamValues(taskID)

			assert.NoError(t, handler.ApproveStateChange(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}

func Test_CancelStateChange(t *testing.T) {
	tests := map[string]func(t *testing.T) (int, *TaskHandler, string, string, *gomock.Controller){
		"success": func(*testing.T) (int, *TaskHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			getTaskResponseJSON := ""
			taskStore := storeMocks.NewMockTaskStoreInterface(ctrl)
			task := model.Task{
				Status:        "task-status",
				ApplicationID: ApplicationID,
				Logs:          []byte("task-logs"),
			}
			task.ID = TaskID
			taskStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&task, nil)
			taskStore.EXPECT().Update(gomock.Any()).Times(1).Return(nil)

			applicationStore := storeMocks.NewMockApplicationStoreInterface(ctrl)
			application := model.Application{}
			application.ID = ApplicationID
			applicationStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&application, nil)
			applicationStore.EXPECT().Update(gomock.Any()).Times(1).Return(nil)
			applicationStore.EXPECT().Delete(gomock.Any()).Times(1).Return(nil)

			applicationStateChangeStore := storeMocks.NewMockApplicationStateChangeStoreInterface(ctrl)
			applicationStateChange := model.ApplicationStateChange{}
			applicationStateChange.ID = ApplicationID
			applicationStateChangeStore.EXPECT().GetByApplicationID(gomock.Any()).Times(1).Return(&applicationStateChange, nil)
			applicationStateChangeStore.EXPECT().Delete(gomock.Any()).Times(1).Return(nil)

			handler := &TaskHandler{
				taskStore:                   taskStore,
				applicationStore:            applicationStore,
				applicationStateChangeStore: applicationStateChangeStore,
			}

			return http.StatusOK, handler, "123", getTaskResponseJSON, ctrl
		},

		"task not found": func(*testing.T) (int, *TaskHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			getTaskResponseJSON := ``

			taskStore := storeMocks.NewMockTaskStoreInterface(ctrl)
			taskStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, errors.New("hello-world error"))

			handler := &TaskHandler{taskStore: taskStore}

			return http.StatusNotFound, handler, "123", getTaskResponseJSON, ctrl
		},
		"application not found": func(*testing.T) (int, *TaskHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			getTaskResponseJSON := ``
			taskStore := storeMocks.NewMockTaskStoreInterface(ctrl)
			task := model.Task{
				Status:        "task-status",
				ApplicationID: ApplicationID,
				Logs:          []byte("task-logs"),
			}
			task.ID = TaskID
			taskStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&task, nil)

			applicationStore := storeMocks.NewMockApplicationStoreInterface(ctrl)
			applicationStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(nil, nil)

			handler := &TaskHandler{
				taskStore:        taskStore,
				applicationStore: applicationStore,
			}

			return http.StatusNotFound, handler, "123", getTaskResponseJSON, ctrl
		},
		"application state change not found": func(*testing.T) (int, *TaskHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			getTaskResponseJSON := ``
			taskStore := storeMocks.NewMockTaskStoreInterface(ctrl)
			task := model.Task{
				Status:        "task-status",
				ApplicationID: ApplicationID,
				Logs:          []byte("task-logs"),
			}
			task.ID = TaskID
			taskStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&task, nil)

			applicationStore := storeMocks.NewMockApplicationStoreInterface(ctrl)
			application := model.Application{}
			application.ID = ApplicationID
			applicationStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&application, nil)

			applicationStateChangeStore := storeMocks.NewMockApplicationStateChangeStoreInterface(ctrl)
			applicationStateChangeStore.EXPECT().GetByApplicationID(gomock.Any()).Times(1).Return(nil, nil)

			handler := &TaskHandler{
				taskStore:                   taskStore,
				applicationStore:            applicationStore,
				applicationStateChangeStore: applicationStateChangeStore,
			}

			return http.StatusNotFound, handler, "123", getTaskResponseJSON, ctrl
		},
		"error deleting application state change": func(*testing.T) (int, *TaskHandler, string, string, *gomock.Controller) {
			ctrl := gomock.NewController(t)
			getTaskResponseJSON := ``
			taskStore := storeMocks.NewMockTaskStoreInterface(ctrl)
			task := model.Task{
				Status:        "task-status",
				ApplicationID: ApplicationID,
				Logs:          []byte("task-logs"),
			}
			task.ID = TaskID
			taskStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&task, nil)

			applicationStore := storeMocks.NewMockApplicationStoreInterface(ctrl)
			application := model.Application{}
			application.ID = ApplicationID
			applicationStore.EXPECT().GetByID(gomock.Any()).Times(1).Return(&application, nil)

			applicationStateChangeStore := storeMocks.NewMockApplicationStateChangeStoreInterface(ctrl)
			applicationStateChange := model.ApplicationStateChange{}
			applicationStateChange.ID = ApplicationID
			applicationStateChangeStore.EXPECT().GetByApplicationID(gomock.Any()).Times(1).Return(&applicationStateChange, nil)
			applicationStateChangeStore.EXPECT().Delete(gomock.Any()).Times(1).Return(errors.New("hello-world error"))

			handler := &TaskHandler{
				taskStore:                   taskStore,
				applicationStore:            applicationStore,
				applicationStateChangeStore: applicationStateChangeStore,
			}

			return http.StatusInternalServerError, handler, "123", getTaskResponseJSON, ctrl
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			expectedStatus, handler, taskID, expectedResponse, ctrl := tc(t)

			e := router.New()
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/tasks/:id/cancel")
			c.SetParamNames("id")
			c.SetParamValues(taskID)
			c.QueryParams().Add("updating", "false")

			assert.NoError(t, handler.CancelStateChange(c))
			assert.Equal(t, expectedStatus, rec.Code)
			if expectedResponse != "" {
				trimmedResponse := strings.TrimSpace(rec.Body.String())
				assert.Equal(t, expectedResponse, trimmedResponse)
			}
			ctrl.Finish()
		})
	}
}
