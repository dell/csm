package handler

import (
	"github.com/dell/csm-deployment/store"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userStore store.UserStoreInterface
}

func New(us store.UserStoreInterface) *UserHandler {
	return &UserHandler{
		userStore: us,
	}
}

func (h *UserHandler) Register(api *echo.Group) {
	adminUsers := api.Group("/users")
	adminUsers.POST("/login", h.Login)
	adminUsers.PATCH("/change-password", h.ChangePasword)
}
