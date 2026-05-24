package artist

import (
	"errors"
	"testing"

	modelartist "github.com/andrqxa/artshare/internal/model/artist"
	"github.com/andrqxa/artshare/internal/model/pagination"
	"github.com/andrqxa/artshare/internal/repository/repoerror"
)

type stubRepository struct {
	create       func(modelartist.CreateInput) (modelartist.Detail, error)
	findByUserID func(string) (modelartist.Detail, error)
	findByID     func(string) (modelartist.Detail, error)
	update       func(string, modelartist.UpdateInput) (modelartist.Detail, error)
	list         func(modelartist.ListFilter) ([]modelartist.Summary, pagination.Meta)
}

func (s *stubRepository) Create(input modelartist.CreateInput) (modelartist.Detail, error) {
	return s.create(input)
}
func (s *stubRepository) FindByUserID(userID string) (modelartist.Detail, error) {
	return s.findByUserID(userID)
}
func (s *stubRepository) FindByID(id string) (modelartist.Detail, error) {
	return s.findByID(id)
}
func (s *stubRepository) UpdateByUserID(userID string, input modelartist.UpdateInput) (modelartist.Detail, error) {
	return s.update(userID, input)
}
func (s *stubRepository) List(filter modelartist.ListFilter) ([]modelartist.Summary, pagination.Meta) {
	return s.list(filter)
}

func TestController_CreateMapsConflict(t *testing.T) {
	repo := &stubRepository{create: func(modelartist.CreateInput) (modelartist.Detail, error) {
		return modelartist.Detail{}, repoerror.ErrConflict
	}}
	c := NewController(repo)

	_, err := c.CreateCurrentArtist(modelartist.CreateInput{UserID: "u1", DisplayName: "x"})
	if !errors.Is(err, ErrArtistAlreadyExists) {
		t.Fatalf("got %v want ErrArtistAlreadyExists", err)
	}
}

func TestController_GetMapsNotFound(t *testing.T) {
	repo := &stubRepository{findByUserID: func(string) (modelartist.Detail, error) {
		return modelartist.Detail{}, repoerror.ErrNotFound
	}}
	c := NewController(repo)

	_, err := c.GetCurrentArtist("u1")
	if !errors.Is(err, ErrArtistNotFound) {
		t.Fatalf("got %v want ErrArtistNotFound", err)
	}
}

func TestController_GetByIDMapsNotFound(t *testing.T) {
	repo := &stubRepository{findByID: func(string) (modelartist.Detail, error) {
		return modelartist.Detail{}, repoerror.ErrNotFound
	}}
	c := NewController(repo)

	_, err := c.GetArtistByID("id")
	if !errors.Is(err, ErrArtistNotFound) {
		t.Fatalf("got %v want ErrArtistNotFound", err)
	}
}

func TestController_UpdateMapsNotFound(t *testing.T) {
	repo := &stubRepository{update: func(string, modelartist.UpdateInput) (modelartist.Detail, error) {
		return modelartist.Detail{}, repoerror.ErrNotFound
	}}
	c := NewController(repo)

	_, err := c.UpdateCurrentArtist("u1", modelartist.UpdateInput{})
	if !errors.Is(err, ErrArtistNotFound) {
		t.Fatalf("got %v want ErrArtistNotFound", err)
	}
}

func TestController_CreatePassesInputThrough(t *testing.T) {
	var captured modelartist.CreateInput
	repo := &stubRepository{create: func(input modelartist.CreateInput) (modelartist.Detail, error) {
		captured = input
		return modelartist.Detail{ID: "1", UserID: input.UserID, DisplayName: input.DisplayName}, nil
	}}
	c := NewController(repo)

	if _, err := c.CreateCurrentArtist(modelartist.CreateInput{UserID: "u1", DisplayName: "Marina"}); err != nil {
		t.Fatalf("CreateCurrentArtist: %v", err)
	}
	if captured.UserID != "u1" || captured.DisplayName != "Marina" {
		t.Fatalf("got %+v want UserID=u1 DisplayName=Marina", captured)
	}
}
