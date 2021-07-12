package handler

import (
	"net/http"
	"strconv"

	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
)

// GetDriverType godoc
// @Summary Get a driver type
// @Description Get a driver type
// @ID get-driver-type
// @Tags driver-type
// @Accept  json
// @Produce  json
// @Param id path string true "Driver Type ID"
// @Success 200 {object} driverResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /driver-types/{id} [get]
func (h *DriverTypeHandler) GetDriverType(c echo.Context) error {
	arrayID := c.Param("id")
	id, err := strconv.Atoi(arrayID)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewErrorResponse(http.StatusUnprocessableEntity, utils.CriticalSeverity, "", err))
	}
	driverType, err := h.driverTypeStore.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}
	if driverType == nil {
		return c.JSON(http.StatusNotFound, utils.NewErrorResponse(http.StatusNotFound, utils.ErrorSeverity, "", err))
	}
	return c.JSON(http.StatusOK, newDriverResponse(driverType))
}

// ListDriverTypes godoc
// @Summary List all driver types
// @Description List all driver types
// @ID list-driver-types
// @Tags driver-type
// @Accept  json
// @Produce  json
// @Success 200 {array} driverResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Security ApiKeyAuth
// @Router /driver-types [get]
func (h *DriverTypeHandler) ListDriverType(c echo.Context) error {
	driverTypes, err := h.driverTypeStore.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, utils.CriticalSeverity, "", err))
	}
	resp := make([]*driverResponse, 0)
	for _, arr := range driverTypes {
		resp = append(resp, newDriverResponse(&arr))
	}
	return c.JSON(http.StatusOK, resp)
}
