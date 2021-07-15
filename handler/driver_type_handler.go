package handler

import (
	"github.com/dell/csm-deployment/store"
	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// DriverTypeHandler is the handler for Driver Type APIs
type DriverTypeHandler struct {
	driverTypeStore store.DriverTypeStoreInterface
}

// NewDriverTypeHandler creates a new DriverTypeHandler
func NewDriverTypeHandler(as store.DriverTypeStoreInterface) *DriverTypeHandler {
	return &DriverTypeHandler{
		driverTypeStore: as,
	}
}

// Register will register all Driver Type APIs
func (h *DriverTypeHandler) Register(api *echo.Group) {
	jwtMiddleware := middleware.JWT(utils.JWTSecret)

	driverType := api.Group("/driver-types", jwtMiddleware)
	driverType.GET("/:id", h.GetDriverType)
	driverType.GET("", h.ListDriverType)
}
