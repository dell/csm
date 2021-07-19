package handler

import (
	"github.com/dell/csm-deployment/store"
	"github.com/labstack/echo/v4"
)

// UserHandler constains the store interface for User
type UserHandler struct {
	userStore store.UserStoreInterface
}

// New create an handler for User
func New(us store.UserStoreInterface) *UserHandler {
	return &UserHandler{
		userStore: us,
	}
}

// Register rosters all the API endpoints for users
func (h *UserHandler) Register(api *echo.Group) {
	adminUsers := api.Group("/users")
	adminUsers.POST("/login", h.Login)
	adminUsers.PATCH("/change-password", h.ChangePasword)
}
