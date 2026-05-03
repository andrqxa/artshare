package artist

import (
	"errors"

	"github.com/andrqxa/artshare/internal/controller/apperror"
	modelartist "github.com/andrqxa/artshare/internal/model/artist"
	"github.com/andrqxa/artshare/internal/model/pagination"
	artistrepository "github.com/andrqxa/artshare/internal/repository/artist"
	"github.com/andrqxa/artshare/internal/repository/repoerror"
)

const currentUserID = "00000000-0000-0000-0000-000000000001"

type Controller struct {
	repository  *artistrepository.Repository
	currentUser string
}

func NewController() *Controller {
	return &Controller{
		repository:  artistrepository.NewRepository(),
		currentUser: currentUserID,
	}
}

func (c *Controller) CreateCurrentArtist(input modelartist.CreateInput) (modelartist.Detail, error) {
	if input.UserID == "" {
		input.UserID = c.currentUser
	}

	artist, err := c.repository.Create(input)
	return artist, controllerError(err)
}

func (c *Controller) GetCurrentArtist(userID string) (modelartist.Detail, error) {
	if userID == "" {
		userID = c.currentUser
	}

	artist, err := c.repository.FindByUserID(userID)
	return artist, controllerError(err)
}

func (c *Controller) UpdateCurrentArtist(userID string, input modelartist.UpdateInput) (modelartist.Detail, error) {
	if userID == "" {
		userID = c.currentUser
	}

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
		return apperror.ErrConflict
	case errors.Is(err, repoerror.ErrNotFound):
		return apperror.ErrNotFound
	default:
		return err
	}
}
