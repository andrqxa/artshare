package like

import (
	"fmt"
	"sync"
	"time"

	modellike "github.com/andrqxa/artshare/internal/model/like"
	"github.com/andrqxa/artshare/internal/repository/repoerror"
)

type Repository struct {
	mu    sync.RWMutex
	likes map[string]modellike.Like
}

func NewRepository() *Repository {
	return &Repository{
		likes: make(map[string]modellike.Like),
	}
}

func (r *Repository) Create(input modellike.CreateInput) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := likeKey(input.UserID, input.ArtworkID)
	if _, exists := r.likes[key]; exists {
		return repoerror.ErrConflict
	}
	r.likes[key] = modellike.Like{
		UserID:    input.UserID,
		ArtworkID: input.ArtworkID,
		CreatedAt: time.Now().UTC(),
	}
	return nil
}

func (r *Repository) Delete(input modellike.DeleteInput) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := likeKey(input.UserID, input.ArtworkID)
	if _, exists := r.likes[key]; !exists {
		return repoerror.ErrNotFound
	}
	delete(r.likes, key)
	return nil
}

func (r *Repository) Exists(userID, artworkID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.likes[likeKey(userID, artworkID)]
	return exists
}

func likeKey(userID, artworkID string) string {
	return fmt.Sprintf("%s:%s", userID, artworkID)
}
