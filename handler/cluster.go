package handler

import (
	"io"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"strconv"

	"github.com/dell/csm-deployment/k8s"

	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
)

// CreateCluster godoc
// @Summary Create a new cluster
// @Description Create a new cluster
// @ID create-cluster
// @Tags cluster
// @Accept  json
// @Produce  json
// @Param   name formData string true  "Name of the cluster"
// @Param   file formData file true  "kube config file"
// @Success 201 {object} clusterResponse
// @Failure 400 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /clusters [post]
func (h *ClusterHandler) CreateCluster(c echo.Context) error {
	// Read form fields
	name := c.FormValue("name")

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	defer src.Close()
	data, err := io.ReadAll(src)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	version, isOpenShit, clientset, err := k8s.DiscoverK8sDetails(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	orchestratorType := model.OrchestratorTypeKubernetes
	if *isOpenShit {
		orchestratorType = model.OrchestratorTypeOpenshift
	}

	cluster := model.Cluster{
		ClusterName:      name,
		ConfigFileData:   data,
		K8sVersion:       version,
		OrchestratorType: orchestratorType,
		Status:           model.ClusterStatusConnected,
	}

	if err := h.clusterStore.Create(&cluster); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	//Update the cluster details asynchronously
	go h.updateClusterDetails(&cluster, clientset, c.Logger())

	return c.JSON(http.StatusCreated, newClusterResponse(&cluster))
}

func (h *ClusterHandler) updateClusterDetails(cluster *model.Cluster, clientset *kubernetes.Clientset, logger echo.Logger) {
	dataCollector := k8s.NodeDataCollector{
		ClientSet: clientset,
		Logger:    logger,
	}
	nodes, err := dataCollector.Collect()
	if err != nil {
		logger.Error("failed to collect node details", err.Error())
		return
	}
	logger.Info(nodes)
	// serialize this list into comma separated strings
	serializedNodes := ""
	for i, node := range nodes {
		if i == 0 {
			serializedNodes = node
		} else {
			serializedNodes = serializedNodes + "," + node
		}
	}
	details := model.ClusterDetails{
		Nodes: serializedNodes,
	}
	err = h.clusterStore.UpdateClusterDetails(cluster, &details)
	if err != nil {
		logger.Error("failed to update cluster details", err.Error())
		return
	}
	logger.Info("Successfully collected node details")
}

// GetCluster godoc
// @Summary Get a cluster
// @Description Get a cluster
// @ID get-cluster
// @Tags cluster
// @Accept  json
// @Produce  json
// @Param id path string true "Cluster ID"
// @Success 200 {object} clusterResponse
// @Failure 400 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /clusters/{id} [get]
func (h *ClusterHandler) GetCluster(c echo.Context) error {
	id := c.Param("id")
	clusterID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	cluster, err := h.clusterStore.GetByClusterID(uint(clusterID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if cluster == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	return c.JSON(http.StatusOK, newClusterResponse(cluster))
}
