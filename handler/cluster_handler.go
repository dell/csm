package handler

import (
	"github.com/dell/csm-deployment/store"
	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"k8s.io/client-go/kubernetes"
)

// K8sClientInterface is an interface for support k8s API operations
//go:generate mockgen -destination=mocks/k8s_client_interface.go -package=mocks github.com/dell/csm-deployment/handler K8sClientInterface
type K8sClientInterface interface {
	DiscoverK8sDetails(data []byte) (string, *bool, *kubernetes.Clientset, error)
}

// ClusterHandler is the handler for Cluster APIs
type ClusterHandler struct {
	clusterStore store.ClusterStoreInterface
	k8sClient    K8sClientInterface
}

// NewClusterHandler creates a new ClusterHandler
func NewClusterHandler(cs store.ClusterStoreInterface, k8sClient K8sClientInterface) *ClusterHandler {
	return &ClusterHandler{
		clusterStore: cs,
		k8sClient:    k8sClient,
	}
}

// Register will register all Cluster APIs
func (h *ClusterHandler) Register(api *echo.Group) {
	jwtMiddleware := middleware.JWT(utils.JWTSecret)

	clusters := api.Group("/clusters", jwtMiddleware)
	clusters.GET("/:id", h.GetCluster)
	clusters.POST("", h.CreateCluster)
	clusters.GET("", h.ListClusters)
	clusters.DELETE("/:id", h.DeleteCluster)
	clusters.PATCH("/:id", h.UpdateCluster)
}
