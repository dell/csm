package handler

import (
	"github.com/dell/csm-deployment/model"
	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

// CreateStorageArray godoc
// @Summary Create a new storage array
// @Description Create a new storage array
// @ID create-storage-array
// @Tags storage-array
// @Accept  json
// @Produce  json
// @Param storageArray body storageArrayCreateRequest true "Storage Array info for creation"
// @Success 202 {object} storageArrayResponse
// @Failure 400 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /storageArrays [post]
func (h *StorageArrayHandler) CreateStorageArray(c echo.Context) error {
	var storageArray model.StorageArray
	req := &storageArrayCreateRequest{}
	if err := req.bind(c, &storageArray); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	arrayType, err := h.arrayStore.GetTypeByTypeName(strings.ToLower(req.StorageArray.StorageArrayType))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	storageArray.StorageArrayTypeID = arrayType.ID

	if err := h.arrayStore.Create(&storageArray); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return c.JSON(http.StatusCreated, newStorageArrayResponse(&storageArray))
}

// ListStorageArrays godoc
// @Summary List all storage arrays
// @Description List all storage arrays
// @ID list-storage-arrays
// @Tags storage-array
// @Accept  json
// @Produce  json
// @Success 202 {object} []storageArrayResponse
// @Failure 400 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /storageArrays [get]
func (h *StorageArrayHandler) ListStorageArrays(c echo.Context) error {
	arrays, err := h.arrayStore.GetAllByID()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	var resp []*storageArrayResponse
	for _, arr := range arrays {
		resp = append(resp, newStorageArrayResponse(&arr))
	}
	return c.JSON(http.StatusCreated, resp)
}

// GetStorageArray godoc
// @Summary Get storage array
// @Description Get storage array
// @ID get-storage-array
// @Tags storage-array
// @Accept  json
// @Produce  json
// @Param id path string true "Storage Array ID"
// @Success 200 {object} storageArrayResponse
// @Failure 400 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /storageArrays/{id} [get]
func (h *StorageArrayHandler) GetStorageArray(c echo.Context) error {
	arrayID := c.Param("id")
	id, err := strconv.Atoi(arrayID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	storageArray, err := h.arrayStore.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if storageArray == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	return c.JSON(http.StatusOK, newStorageArrayResponse(storageArray))
}
