package handler

import (
	"github.com/dell/csm-deployment/store"
	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type StorageArrayHandler struct {
	arrayStore store.StorageArrayStoreInterface
}

func NewStorageArrayHandler(as store.StorageArrayStoreInterface) *StorageArrayHandler {
	return &StorageArrayHandler{
		arrayStore: as,
	}
}

func (h *StorageArrayHandler) Register(api *echo.Group) {
	jwtMiddleware := middleware.JWT(utils.JWTSecret)

	storageArrays := api.Group("/storage-arrays", jwtMiddleware)
	storageArrays.POST("", h.CreateStorageArray)
	storageArrays.GET("", h.ListStorageArrays)
	storageArrays.GET("/:id", h.GetStorageArray)
	storageArrays.DELETE("/:id", h.DeleteStorageArray)
	storageArrays.PUT("/:id", h.UpdateStorageArray)
}
