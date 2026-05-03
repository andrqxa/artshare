package like

import (
	"fmt"
	"sync"
	"time"

	"github.com/andrqxa/artshare/internal/controller/apperror"
	modellike "github.com/andrqxa/artshare/internal/model/like"
)

type Controller struct {
	mu    sync.RWMutex
	likes map[string]modellike.Like
}

func NewController() *Controller {
	return &Controller{
		likes: make(map[string]modellike.Like),
	}
}

func (c *Controller) LikeArtwork(input modellike.CreateInput) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := likeKey(input.UserID, input.ArtworkID)
	if _, exists := c.likes[key]; exists {
		return apperror.ErrConflict
	}
	c.likes[key] = modellike.Like{
		UserID:    input.UserID,
		ArtworkID: input.ArtworkID,
		CreatedAt: time.Now().UTC(),
	}
	return nil
}

func (c *Controller) UnlikeArtwork(input modellike.DeleteInput) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := likeKey(input.UserID, input.ArtworkID)
	if _, exists := c.likes[key]; !exists {
		return apperror.ErrNotFound
	}
	delete(c.likes, key)
	return nil
}

func (c *Controller) IsLiked(userID, artworkID string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, exists := c.likes[likeKey(userID, artworkID)]
	return exists
}

func likeKey(userID, artworkID string) string {
	return fmt.Sprintf("%s:%s", userID, artworkID)
}
