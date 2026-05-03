package like

import (
	"errors"

	"github.com/andrqxa/artshare/internal/controller/apperror"
	modellike "github.com/andrqxa/artshare/internal/model/like"
	likerepository "github.com/andrqxa/artshare/internal/repository/like"
	"github.com/andrqxa/artshare/internal/repository/repoerror"
)

type Controller struct {
	repository *likerepository.Repository
}

func NewController() *Controller {
	return &Controller{
		repository: likerepository.NewRepository(),
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
		return apperror.ErrConflict
	case errors.Is(err, repoerror.ErrNotFound):
		return apperror.ErrNotFound
	default:
		return err
	}
}
