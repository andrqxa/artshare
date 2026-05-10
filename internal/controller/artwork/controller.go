package artwork

import (
	"errors"

	"github.com/andrqxa/artshare/internal/controller/apperror"
	modelartwork "github.com/andrqxa/artshare/internal/model/artwork"
	"github.com/andrqxa/artshare/internal/model/pagination"
	"github.com/andrqxa/artshare/internal/repository/repoerror"
)

type Repository interface {
	List(filter modelartwork.ListFilter) ([]modelartwork.Summary, pagination.Meta)
	Create(input modelartwork.CreateInput) (modelartwork.Detail, error)
	FindByID(id string) (modelartwork.Detail, error)
	Update(input modelartwork.UpdateInput) (modelartwork.Detail, error)
	Delete(id string) error
}

type Controller struct {
	repository Repository
}

func NewController(repository Repository) *Controller {
	return &Controller{
		repository: repository,
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
