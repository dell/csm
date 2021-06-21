package handler

import (
	"github.com/dell/csm-deployment/store"
	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"k8s.io/client-go/kubernetes"
)

//go:generate mockgen -destination=mocks/k8s_client_interface.go -package=mocks github.com/dell/csm-deployment/handler K8sClientInterface
type K8sClientInterface interface {
	DiscoverK8sDetails(data []byte) (string, *bool, *kubernetes.Clientset, error)
}

type ClusterHandler struct {
	clusterStore store.ClusterStoreInterface
	k8sClient    K8sClientInterface
}

func NewClusterHandler(cs store.ClusterStoreInterface, k8sClient K8sClientInterface) *ClusterHandler {
	return &ClusterHandler{
		clusterStore: cs,
		k8sClient:    k8sClient,
	}
}

func (h *ClusterHandler) Register(api *echo.Group) {
	jwtMiddleware := middleware.JWT(utils.JWTSecret)

	clusters := api.Group("/clusters", jwtMiddleware)
	clusters.GET("/:id", h.GetCluster)
	clusters.POST("", h.CreateCluster)
	clusters.GET("", h.ListClusters)
	clusters.DELETE("/:id", h.DeleteCluster)
	clusters.PATCH("/:id", h.UpdateCluster)
}
