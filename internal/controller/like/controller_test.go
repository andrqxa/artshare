package like

import (
	"errors"
	"testing"

	modellike "github.com/andrqxa/artshare/internal/model/like"
	"github.com/andrqxa/artshare/internal/repository/repoerror"
)

type stubRepository struct {
	createErr error
	deleteErr error
	exists    bool
}

func (s *stubRepository) Create(modellike.CreateInput) error {
	return s.createErr
}
func (s *stubRepository) Delete(modellike.DeleteInput) error {
	return s.deleteErr
}
func (s *stubRepository) Exists(string, string) bool {
	return s.exists
}

func TestController_LikeMapsConflict(t *testing.T) {
	c := NewController(&stubRepository{createErr: repoerror.ErrConflict})
	if err := c.LikeArtwork(modellike.CreateInput{}); !errors.Is(err, ErrLikeConflict) {
		t.Fatalf("got %v want ErrLikeConflict", err)
	}
}

func TestController_UnlikeMapsNotFound(t *testing.T) {
	c := NewController(&stubRepository{deleteErr: repoerror.ErrNotFound})
	if err := c.UnlikeArtwork(modellike.DeleteInput{}); !errors.Is(err, ErrLikeNotFound) {
		t.Fatalf("got %v want ErrLikeNotFound", err)
	}
}

func TestController_IsLikedDelegates(t *testing.T) {
	c := NewController(&stubRepository{exists: true})
	if !c.IsLiked("u", "a") {
		t.Fatalf("expected true")
	}
}
