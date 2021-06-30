package handler

import (
	"github.com/dell/csm-deployment/kapp"
	"github.com/dell/csm-deployment/store"
	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type TaskHandler struct {
	taskStore                   store.TaskStoreInterface
	applicationStore            store.ApplicationStoreInterface
	applicationStateChangeStore store.ApplicationStateChangeStoreInterface
	clusterStore                store.ClusterStoreInterface
	kappClient                  kapp.Interface
}

func NewTaskHandler(ts store.TaskStoreInterface, as store.ApplicationStoreInterface, asc store.ApplicationStateChangeStoreInterface, cs store.ClusterStoreInterface, kapp kapp.Interface) *TaskHandler {
	return &TaskHandler{
		taskStore:                   ts,
		applicationStore:            as,
		applicationStateChangeStore: asc,
		clusterStore:                cs,
		kappClient:                  kapp,
	}
}

func (h *TaskHandler) Register(api *echo.Group) {
	jwtMiddleware := middleware.JWT(utils.JWTSecret)

	tasks := api.Group("/tasks", jwtMiddleware)
	tasks.GET("/:id", h.GetTask)
	tasks.POST("/:id/approve", h.ApproveStateChange)
	tasks.POST("/:id/cancel", h.CancelStateChange)
}
