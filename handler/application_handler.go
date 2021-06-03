package handler

import (
	"github.com/dell/csm-deployment/store"
	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ApplicationHandler struct {
	applicationStore            store.ApplicationStoreInterface
	arrayStore                  store.StorageArrayStoreInterface
	taskStore                   store.TaskStoreInterface
	clusterStore                store.ClusterStoreInterface
	applicationStateChangeStore store.ApplicationStateChangeStoreInterface
	moduleStore                 store.ModuleStoreInterface
}

func NewApplicationHandler(is store.ApplicationStoreInterface,
	ts store.TaskStoreInterface,
	cs store.ClusterStoreInterface,
	asc store.ApplicationStateChangeStoreInterface,
	as store.StorageArrayStoreInterface,
	ms store.ModuleStoreInterface) *ApplicationHandler {
	return &ApplicationHandler{
		applicationStore:            is,
		taskStore:                   ts,
		clusterStore:                cs,
		applicationStateChangeStore: asc,
		arrayStore:                  as,
		moduleStore:                 ms,
	}
}

func (h *ApplicationHandler) Register(api *echo.Group) {
	jwtMiddleware := middleware.JWT(utils.JWTSecret)

	applications := api.Group("/applications", jwtMiddleware)
	applications.POST("", h.CreateApplication)
	applications.PUT("/:id", h.UpdateApplication)
	applications.GET("/:id", h.GetApplication)
	applications.DELETE("/:id", h.DeleteApplication)
}
