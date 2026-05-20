package artist

import (
	"errors"

	"github.com/andrqxa/artshare/internal/controller/apperror"
	modelartist "github.com/andrqxa/artshare/internal/model/artist"
	"github.com/andrqxa/artshare/internal/model/pagination"
	"github.com/andrqxa/artshare/internal/repository/repoerror"
)

type Repository interface {
	Create(input modelartist.CreateInput) (modelartist.Detail, error)
	FindByUserID(userID string) (modelartist.Detail, error)
	UpdateByUserID(userID string, input modelartist.UpdateInput) (modelartist.Detail, error)
	List(filter modelartist.ListFilter) ([]modelartist.Summary, pagination.Meta)
	FindByID(id string) (modelartist.Detail, error)
}

type Controller struct {
	repository Repository
}

func NewController(repository Repository) *Controller {
	return &Controller{
		repository: repository,
	}
}

func (c *Controller) CreateCurrentArtist(input modelartist.CreateInput) (modelartist.Detail, error) {
	artist, err := c.repository.Create(input)
	return artist, controllerError(err)
}

func (c *Controller) GetCurrentArtist(userID string) (modelartist.Detail, error) {
	artist, err := c.repository.FindByUserID(userID)
	return artist, controllerError(err)
}

func (c *Controller) UpdateCurrentArtist(userID string, input modelartist.UpdateInput) (modelartist.Detail, error) {
	artist, err := c.repository.UpdateByUserID(userID, input)
	return artist, controllerError(err)
}

func (c *Controller) ListArtists(filter modelartist.ListFilter) ([]modelartist.Summary, pagination.Meta) {
	return c.repository.List(filter)
}

func (c *Controller) GetArtistByID(id string) (modelartist.Detail, error) {
	artist, err := c.repository.FindByID(id)
	return artist, controllerError(err)
}

func controllerError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, repoerror.ErrConflict):
		return ErrArtistAlreadyExists
	case errors.Is(err, repoerror.ErrNotFound):
		return ErrArtistNotFound
	default:
		return apperror.WrapInternal(err)
	}
}
