package handler

import (
	"github.com/dell/csm-deployment/store"
	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ClusterHandler struct {
	clusterStore store.ClusterStoreInterface
}

func NewClusterHandler(cs store.ClusterStoreInterface) *ClusterHandler {
	return &ClusterHandler{
		clusterStore: cs,
	}
}

func (h *ClusterHandler) Register(api *echo.Group) {
	jwtMiddleware := middleware.JWT(utils.JWTSecret)

	clusters := api.Group("/clusters", jwtMiddleware)
	clusters.GET("/:id", h.GetCluster)
	clusters.POST("", h.CreateCluster)
}
