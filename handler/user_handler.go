package handler

import (
	"github.com/dell/csm-deployment/store"
	"github.com/labstack/echo/v4"
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
	adminUsers := api.Group("/users")
	adminUsers.POST("/login", h.Login)
	adminUsers.PATCH("/update", h.UpdateUser)
}
