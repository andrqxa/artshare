package like

import (
	"errors"
	"testing"

	modellike "github.com/andrqxa/artshare/internal/model/like"
	"github.com/andrqxa/artshare/internal/repository/repoerror"
)

func TestRepository_CreateDeleteExists(t *testing.T) {
	repo := NewRepository()

	input := modellike.CreateInput{UserID: "u1", ArtworkID: "a1"}
	if err := repo.Create(input); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if !repo.Exists("u1", "a1") {
		t.Fatalf("Exists: expected true")
	}

	if err := repo.Create(input); !errors.Is(err, repoerror.ErrConflict) {
		t.Fatalf("duplicate Create: got %v want ErrConflict", err)
	}

	if err := repo.Delete(modellike.DeleteInput{UserID: "u1", ArtworkID: "a1"}); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if repo.Exists("u1", "a1") {
		t.Fatalf("Exists after delete: expected false")
	}

	if err := repo.Delete(modellike.DeleteInput{UserID: "u1", ArtworkID: "a1"}); !errors.Is(err, repoerror.ErrNotFound) {
		t.Fatalf("Delete missing: got %v want ErrNotFound", err)
	}
}
