package artist

import (
	"errors"
	"testing"

	modelartist "github.com/andrqxa/artshare/internal/model/artist"
	"github.com/andrqxa/artshare/internal/repository/repoerror"
)

const userID = "00000000-0000-0000-0000-000000000001"

func TestRepository_CreateAndFind(t *testing.T) {
	repo := NewRepository()

	created, err := repo.Create(modelartist.CreateInput{
		UserID:      userID,
		DisplayName: "Marina",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID == "" {
		t.Fatalf("expected ID, got empty")
	}

	byUser, err := repo.FindByUserID(userID)
	if err != nil {
		t.Fatalf("FindByUserID: %v", err)
	}
	if byUser.ID != created.ID {
		t.Fatalf("FindByUserID returned %q, want %q", byUser.ID, created.ID)
	}

	byID, err := repo.FindByID(created.ID)
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if byID.UserID != userID {
		t.Fatalf("FindByID returned UserID %q, want %q", byID.UserID, userID)
	}
}

func TestRepository_CreateConflict(t *testing.T) {
	repo := NewRepository()
	if _, err := repo.Create(modelartist.CreateInput{UserID: userID, DisplayName: "A"}); err != nil {
		t.Fatalf("first create: %v", err)
	}

	_, err := repo.Create(modelartist.CreateInput{UserID: userID, DisplayName: "B"})
	if !errors.Is(err, repoerror.ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}

func TestRepository_FindMissing(t *testing.T) {
	repo := NewRepository()

	if _, err := repo.FindByUserID("missing"); !errors.Is(err, repoerror.ErrNotFound) {
		t.Fatalf("FindByUserID missing: got %v, want ErrNotFound", err)
	}
	if _, err := repo.FindByID("missing"); !errors.Is(err, repoerror.ErrNotFound) {
		t.Fatalf("FindByID missing: got %v, want ErrNotFound", err)
	}
	if _, err := repo.UpdateByUserID("missing", modelartist.UpdateInput{}); !errors.Is(err, repoerror.ErrNotFound) {
		t.Fatalf("UpdateByUserID missing: got %v, want ErrNotFound", err)
	}
}

func TestRepository_UpdateByUserID(t *testing.T) {
	repo := NewRepository()
	created, err := repo.Create(modelartist.CreateInput{UserID: userID, DisplayName: "Old"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	newName := "New"
	updated, err := repo.UpdateByUserID(userID, modelartist.UpdateInput{DisplayName: &newName})
	if err != nil {
		t.Fatalf("UpdateByUserID: %v", err)
	}
	if updated.DisplayName != newName {
		t.Fatalf("DisplayName: got %q want %q", updated.DisplayName, newName)
	}
	if !updated.UpdatedAt.After(created.UpdatedAt) {
		t.Fatalf("expected UpdatedAt to advance, got %v vs %v", updated.UpdatedAt, created.UpdatedAt)
	}
}

func TestRepository_ListFilterAndPagination(t *testing.T) {
	repo := NewRepository()

	users := []struct {
		id   string
		name string
	}{
		{"u1", "Alice"},
		{"u2", "Bob"},
		{"u3", "Carol"},
	}
	for _, u := range users {
		if _, err := repo.Create(modelartist.CreateInput{UserID: u.id, DisplayName: u.name}); err != nil {
			t.Fatalf("seed %s: %v", u.id, err)
		}
	}

	all, meta := repo.List(modelartist.ListFilter{Page: 1, PageSize: 10})
	if meta.TotalItems != 3 {
		t.Fatalf("totalItems: got %d want 3", meta.TotalItems)
	}
	if len(all) != 3 {
		t.Fatalf("returned: got %d want 3", len(all))
	}
	if all[0].DisplayName != "Alice" {
		t.Fatalf("expected alphabetical ordering, got %q first", all[0].DisplayName)
	}

	query := "bob"
	filtered, meta := repo.List(modelartist.ListFilter{Query: &query, Page: 1, PageSize: 10})
	if meta.TotalItems != 1 || len(filtered) != 1 || filtered[0].DisplayName != "Bob" {
		t.Fatalf("query filter: got %+v meta=%+v", filtered, meta)
	}

	page2, meta := repo.List(modelartist.ListFilter{Page: 2, PageSize: 2})
	if meta.TotalPages != 2 {
		t.Fatalf("totalPages: got %d want 2", meta.TotalPages)
	}
	if len(page2) != 1 {
		t.Fatalf("page 2 size: got %d want 1", len(page2))
	}
}
