package handler

import (
	"github.com/dell/csm-deployment/store"
	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// ModuleTypeHandler is the handler for Module Type APIs
type ModuleTypeHandler struct {
	moduleTypeStore store.ModuleStoreInterface
}

// NewModuleTypeHandler creates a new ModuleTypeHandler
func NewModuleTypeHandler(as store.ModuleStoreInterface) *ModuleTypeHandler {
	return &ModuleTypeHandler{
		moduleTypeStore: as,
	}
}

// Register will register all Module Type APIs
func (h *ModuleTypeHandler) Register(api *echo.Group) {
	jwtMiddleware := middleware.JWT(utils.JWTSecret)

	moduleType := api.Group("/module-types", jwtMiddleware)
	moduleType.GET("/:id", h.GetModuleType)
	moduleType.GET("", h.ListModuleType)
}
