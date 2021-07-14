package handler

import (
	"github.com/dell/csm-deployment/store"
	"github.com/dell/csm-deployment/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handler struct {
	userStore store.UserStoreInterface
}

func New(us store.UserStoreInterface) *Handler {
	return &Handler{
		userStore: us,
	}
}

func (h *Handler) Register(api *echo.Group) {
	jwtMiddleware := middleware.JWT(utils.JWTSecret)

	adminUsers := api.Group("/users", jwtMiddleware)
	adminUsers.POST("/login", h.Login)
}
