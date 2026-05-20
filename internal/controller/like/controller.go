package like

import (
	"errors"

	"github.com/andrqxa/artshare/internal/controller/apperror"
	modellike "github.com/andrqxa/artshare/internal/model/like"
	"github.com/andrqxa/artshare/internal/repository/repoerror"
)

type Repository interface {
	Create(input modellike.CreateInput) error
	Delete(input modellike.DeleteInput) error
	Exists(userID, artworkID string) bool
}

type Controller struct {
	repository Repository
}

func NewController(repository Repository) *Controller {
	return &Controller{
		repository: repository,
	}
}

func (c *Controller) LikeArtwork(input modellike.CreateInput) error {
	return controllerError(c.repository.Create(input))
}

func (c *Controller) UnlikeArtwork(input modellike.DeleteInput) error {
	return controllerError(c.repository.Delete(input))
}

func (c *Controller) IsLiked(userID, artworkID string) bool {
	return c.repository.Exists(userID, artworkID)
}

func controllerError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, repoerror.ErrConflict):
		return ErrLikeConflict
	case errors.Is(err, repoerror.ErrNotFound):
		return ErrLikeNotFound
	default:
		return apperror.WrapInternal(err)
	}
}
