package handler

import (
	"github.com/dell/csm-deployment/store"
	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// ApplicationHandler - Place holder for Api Interfaces
type ApplicationHandler struct {
	applicationStore            store.ApplicationStoreInterface
	arrayStore                  store.StorageArrayStoreInterface
	taskStore                   store.TaskStoreInterface
	clusterStore                store.ClusterStoreInterface
	applicationStateChangeStore store.ApplicationStateChangeStoreInterface
	ModuleTypeStore             store.ModuleTypeStoreInterface
	driverStore                 store.DriverTypeStoreInterface
	k8sClient                   K8sClientInterface
}

// NewApplicationHandler -  returns a new ApplicationHandler
func NewApplicationHandler(is store.ApplicationStoreInterface,
	ts store.TaskStoreInterface,
	cs store.ClusterStoreInterface,
	asc store.ApplicationStateChangeStoreInterface,
	as store.StorageArrayStoreInterface,
	ms store.ModuleTypeStoreInterface,
	ds store.DriverTypeStoreInterface,
	k8sClient K8sClientInterface) *ApplicationHandler {
	return &ApplicationHandler{
		applicationStore:            is,
		taskStore:                   ts,
		clusterStore:                cs,
		applicationStateChangeStore: asc,
		arrayStore:                  as,
		ModuleTypeStore:             ms,
		driverStore:                 ds,
		k8sClient:                   k8sClient,
	}
}

//Register new Application Handler
func (h *ApplicationHandler) Register(api *echo.Group) {
	jwtMiddleware := middleware.JWT(utils.JWTSecret)

	applications := api.Group("/applications", jwtMiddleware)
	applications.POST("", h.CreateApplication)
	applications.PUT("/:id", h.UpdateApplication)
	applications.GET("/:id", h.GetApplication)
	applications.DELETE("/:id", h.DeleteApplication)
}
