package user

import (
	"errors"
	"fmt"

	"github.com/andrqxa/artshare/internal/controller/apperror"
	modeluser "github.com/andrqxa/artshare/internal/model/user"
	"github.com/andrqxa/artshare/internal/repository/repoerror"
)

type Repository interface {
	Create(input modeluser.RegisterInput) (modeluser.CurrentUser, error)
	FindByEmail(email string) (modeluser.CurrentUser, string, error)
	GetCurrent() modeluser.CurrentUser
	SetCurrent(user modeluser.CurrentUser)
	UpdateCurrent(input modeluser.UpdateCurrentUserInput) modeluser.CurrentUser
}

type Controller struct {
	repository Repository
}

func NewController(repository Repository) *Controller {
	return &Controller{
		repository: repository,
	}
}

func (c *Controller) RegisterUser(input modeluser.RegisterInput) (modeluser.CurrentUser, modeluser.AuthToken, error) {
	user, err := c.repository.Create(input)
	if err != nil {
		return modeluser.CurrentUser{}, modeluser.AuthToken{}, controllerError(err)
	}

	return user, authTokenForUser(user), nil
}

func (c *Controller) LoginUser(input modeluser.LoginInput) (modeluser.CurrentUser, modeluser.AuthToken, error) {
	user, password, err := c.repository.FindByEmail(input.Email)
	if err != nil || password != input.Password {
		return modeluser.CurrentUser{}, modeluser.AuthToken{}, apperror.ErrUnauthorized
	}

	c.repository.SetCurrent(user)
	return user, authTokenForUser(user), nil
}

func (c *Controller) GetCurrentUser() modeluser.CurrentUser {
	return c.repository.GetCurrent()
}

func (c *Controller) UpdateCurrentUser(input modeluser.UpdateCurrentUserInput) modeluser.CurrentUser {
	return c.repository.UpdateCurrent(input)
}

func controllerError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, repoerror.ErrConflict):
		return apperror.ErrConflict
	case errors.Is(err, repoerror.ErrNotFound):
		return apperror.ErrNotFound
	default:
		return err
	}
}

func authTokenForUser(user modeluser.CurrentUser) modeluser.AuthToken {
	return modeluser.AuthToken{
		AccessToken: fmt.Sprintf("dev-token-%s", user.ID),
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}
}
