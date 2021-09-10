// Copyright (c) 2021 Dell Inc., or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0

package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"k8s.io/client-go/kubernetes"

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
// @Accept json
// @Produce json
// @Param name formData string true "Name of the cluster"
// @Param file formData file true "kube config file"
// @Success 201 {object} clusterResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /clusters [post]
func (h *ClusterHandler) CreateCluster(c echo.Context) error {
	cluster := &model.Cluster{}

	// Read form fields
	name := c.FormValue("name")
	if len(name) == 0 {
		err := errors.New("name is required")
		if err != nil {
			return c.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, utils.CriticalSeverity, "", err))
		}
	}

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, utils.CriticalSeverity, "", err))
	}
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, utils.CriticalSeverity, "", err))
	}

	defer src.Close()
	data, err := io.ReadAll(src)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}

	version, isOpenShift, clientset, err := h.k8sClient.DiscoverK8sDetails(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, utils.CriticalSeverity, "", err))
	}

	orchestratorType := model.OrchestratorTypeKubernetes
	if *isOpenShift {
		orchestratorType = model.OrchestratorTypeOpenshift
	}

	cluster.ClusterName = name
	cluster.ConfigFileData = data
	cluster.K8sVersion = version
	cluster.OrchestratorType = orchestratorType
	cluster.Status = model.ClusterStatusConnected

	if err := h.clusterStore.Create(cluster); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}

	//Update the cluster details asynchronously
	go h.updateClusterDetails(cluster, clientset, c.Logger())

	return c.JSON(http.StatusCreated, newClusterResponse(cluster))
}

// UpdateCluster godoc
// @Summary Update a cluster
// @Description Update a cluster
// @ID update-cluster
// @Tags cluster
// @Accept  json
// @Produce  json
// @Param id path string true "Cluster ID"
// @Param name formData string false "Name of the cluster"
// @Param file formData file false "kube config file"
// @Success 200 {object} clusterResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /clusters/{id} [patch]
func (h *ClusterHandler) UpdateCluster(c echo.Context) error {

	id := c.Param("id")
	clusterID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.ErrorSeverity, "", err))
	}
	cluster, err := h.clusterStore.GetByID(uint(clusterID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}
	if cluster == nil {
		return c.JSON(http.StatusNotFound, utils.NewErrorResponse(http.StatusNotFound, utils.ErrorSeverity, "", err))
	}

	// cluster name is optional during update, if not provided use existing cluster name
	name := c.FormValue("name")
	if len(name) == 0 {
		name = cluster.ClusterName
	}

	// configfile is optional during update, if not provided, use existing config file
	file, err := c.FormFile("file")
	var data []byte
	if err != nil {
		data = cluster.ConfigFileData
	} else {
		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, utils.CriticalSeverity, "", err))
		}
		defer src.Close()
		data, err = io.ReadAll(src)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
		}
	}

	version, isOpenShift, clientset, err := h.k8sClient.DiscoverK8sDetails(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, utils.CriticalSeverity, "", err))
	}

	orchestratorType := model.OrchestratorTypeKubernetes
	if *isOpenShift {
		orchestratorType = model.OrchestratorTypeOpenshift
	}

	cluster.ClusterName = name
	cluster.ConfigFileData = data
	cluster.K8sVersion = version
	cluster.OrchestratorType = orchestratorType
	cluster.Status = model.ClusterStatusConnected

	if err := h.clusterStore.Update(cluster); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}

	//Update the cluster details asynchronously
	go h.updateClusterDetails(cluster, clientset, c.Logger())

	return c.JSON(http.StatusOK, newClusterResponse(cluster))
}

func (h *ClusterHandler) updateClusterDetails(cluster *model.Cluster, clientset *kubernetes.Clientset, logger echo.Logger) {
	dataCollector := k8s.NodeDataCollector{
		ClientSet:                  clientset,
		Logger:                     logger,
		InstallWaitTime:            30 * time.Second,
		HandleTerminatedPodTimeout: 5 * time.Minute,
		HandlePendingPodTimeout:    3 * time.Minute,
		HandlePendingPodsWaitTime:  3 * time.Second,
	}
	nodes, err := dataCollector.Collect()
	if err != nil {
		logger.Error("failed to collect node details", err.Error())
		return
	}
	logger.Info(nodes)

	nodeDetails, err := json.Marshal(nodes)
	if err != nil {
		logger.Error("failed to marshal node details node details", err.Error())
		return
	}

	details := model.ClusterDetails{
		Nodes: string(nodeDetails),
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
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /clusters/{id} [get]
func (h *ClusterHandler) GetCluster(c echo.Context) error {
	id := c.Param("id")
	clusterID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.ErrorSeverity, "", err))
	}
	cluster, err := h.clusterStore.GetByID(uint(clusterID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}
	if cluster == nil {
		return c.JSON(http.StatusNotFound, utils.NewErrorResponse(http.StatusNotFound, utils.ErrorSeverity, "", err))
	}
	return c.JSON(http.StatusOK, newClusterResponse(cluster))
}

// ListClusters godoc
// @Summary List all clusters
// @Description List all clusters
// @ID list-clusters
// @Tags cluster
// @Accept  json
// @Produce  json
// @Param cluster_name query string false "Cluster Name"
// @Success 200 {array} clusterResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /clusters [get]
func (h *ClusterHandler) ListClusters(c echo.Context) error {

	name := c.QueryParam("cluster_name")
	var clusters []model.Cluster
	var err error
	if name != "" {
		clusters, err = h.clusterStore.GetAllByName(name)
	} else {
		clusters, err = h.clusterStore.GetAll()
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}

	resp := make([]*clusterResponse, 0)
	for _, cluster := range clusters {
		resp = append(resp, newClusterResponse(&cluster))
	}
	return c.JSON(http.StatusOK, resp)
}

// DeleteCluster godoc
// @Summary Delete a cluster
// @Description Delete a cluster
// @ID delete-cluster
// @Tags cluster
// @Accept  json
// @Produce  json
// @Param id path string true "Cluster ID"
// @Success 204
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /clusters/{id} [delete]
func (h *ClusterHandler) DeleteCluster(c echo.Context) error {
	clusterID := c.Param("id")
	id, err := strconv.Atoi(clusterID)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.CriticalSeverity, "", err))
	}
	cluster, err := h.clusterStore.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}
	if cluster == nil {
		return c.JSON(http.StatusNotFound, utils.NewErrorResponse(http.StatusNotFound, utils.ErrorSeverity, "", err))
	}
	if err := h.clusterStore.Delete(cluster); err != nil {
		c.Logger().Errorf("error deleting cluster: %+v", err)
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}
	return c.JSON(http.StatusNoContent, nil)
}
