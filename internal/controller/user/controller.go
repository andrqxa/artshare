package user

import (
	modeluser "artshare/internal/model/user"
	"sync"
	"time"
)

type Controller struct {
	mu          sync.RWMutex
	currentUser modeluser.CurrentUser
}

func NewController() *Controller {
	return &Controller{
		currentUser: modeluser.CurrentUser{
			ID:        "00000000-0000-0000-0000-000000000001",
			Email:     "viewer@artshare.local",
			Role:      modeluser.RoleViewer,
			CreatedAt: time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	}
}

func (c *Controller) GetCurrentUser() modeluser.CurrentUser {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.currentUser
}

func (c *Controller) UpdateCurrentUser(input modeluser.UpdateCurrentUserInput) modeluser.CurrentUser {
	c.mu.Lock()
	defer c.mu.Unlock()

	if input.Email != nil {
		c.currentUser.Email = *input.Email
	}

	return c.currentUser
}
