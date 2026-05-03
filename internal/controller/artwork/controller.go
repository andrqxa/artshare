package artwork

import (
	"errors"

	"github.com/andrqxa/artshare/internal/controller/apperror"
	modelartwork "github.com/andrqxa/artshare/internal/model/artwork"
	"github.com/andrqxa/artshare/internal/model/pagination"
	artworkrepository "github.com/andrqxa/artshare/internal/repository/artwork"
	"github.com/andrqxa/artshare/internal/repository/repoerror"
)

type Controller struct {
	repository *artworkrepository.Repository
}

func NewController() *Controller {
	return &Controller{
		repository: artworkrepository.NewRepository(),
	}
}

func (c *Controller) ListArtworks(filter modelartwork.ListFilter) ([]modelartwork.Summary, pagination.Meta) {
	return c.repository.List(filter)
}

func (c *Controller) CreateArtwork(input modelartwork.CreateInput) (modelartwork.Detail, error) {
	artwork, err := c.repository.Create(input)
	return artwork, controllerError(err)
}

func (c *Controller) GetArtworkByID(id string) (modelartwork.Detail, error) {
	artwork, err := c.repository.FindByID(id)
	return artwork, controllerError(err)
}

func (c *Controller) UpdateArtwork(input modelartwork.UpdateInput) (modelartwork.Detail, error) {
	artwork, err := c.repository.Update(input)
	return artwork, controllerError(err)
}

func (c *Controller) DeleteArtwork(id string) error {
	return controllerError(c.repository.Delete(id))
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
