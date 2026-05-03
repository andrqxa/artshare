package user

import (
	"fmt"
	"sync"
	"time"

	"github.com/andrqxa/artshare/internal/controller/apperror"
	modeluser "github.com/andrqxa/artshare/internal/model/user"
)

type Controller struct {
	mu          sync.RWMutex
	nextID      int
	users       map[string]registeredUser
	currentUser modeluser.CurrentUser
}

type registeredUser struct {
	User     modeluser.CurrentUser
	Password string
}

func NewController() *Controller {
	currentUser := modeluser.CurrentUser{
		ID:        "00000000-0000-0000-0000-000000000001",
		Email:     "viewer@artshare.local",
		Role:      modeluser.RoleViewer,
		CreatedAt: time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	return &Controller{
		nextID:      2,
		currentUser: currentUser,
		users: map[string]registeredUser{
			currentUser.Email: {
				User:     currentUser,
				Password: "password",
			},
		},
	}
}

func (c *Controller) RegisterUser(input modeluser.RegisterInput) (modeluser.CurrentUser, modeluser.AuthToken, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.users[input.Email]; exists {
		return modeluser.CurrentUser{}, modeluser.AuthToken{}, apperror.ErrConflict
	}

	user := modeluser.CurrentUser{
		ID:        nextUUID(c.nextID),
		Email:     input.Email,
		Role:      input.Role,
		CreatedAt: time.Now().UTC(),
	}
	c.nextID++
	c.users[user.Email] = registeredUser{
		User:     user,
		Password: input.Password,
	}
	c.currentUser = user

	return user, authTokenForUser(user), nil
}

func (c *Controller) LoginUser(input modeluser.LoginInput) (modeluser.CurrentUser, modeluser.AuthToken, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	registered, exists := c.users[input.Email]
	if !exists || registered.Password != input.Password {
		return modeluser.CurrentUser{}, modeluser.AuthToken{}, apperror.ErrUnauthorized
	}

	c.currentUser = registered.User
	return registered.User, authTokenForUser(registered.User), nil
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
		registered := c.users[c.currentUser.Email]
		delete(c.users, c.currentUser.Email)
		c.currentUser.Email = *input.Email
		registered.User = c.currentUser
		c.users[c.currentUser.Email] = registeredUser{
			User:     registered.User,
			Password: registered.Password,
		}
	}

	return c.currentUser
}

func authTokenForUser(user modeluser.CurrentUser) modeluser.AuthToken {
	return modeluser.AuthToken{
		AccessToken: fmt.Sprintf("dev-token-%s", user.ID),
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}
}

func nextUUID(id int) string {
	return fmt.Sprintf("00000000-0000-0000-0000-%012d", id)
}
