package like

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/andrqxa/artshare/internal/controller/apperror"
	"github.com/andrqxa/artshare/internal/controller/like"
	"github.com/andrqxa/artshare/internal/middleware/currentuser"
	modellike "github.com/andrqxa/artshare/internal/model/like"
	routererror "github.com/andrqxa/artshare/internal/router/error"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	controller *like.Controller
}

func NewHandler(controller *like.Controller) *Handler {
	return &Handler{controller: controller}
}

func (h *Handler) LikeArtwork(w http.ResponseWriter, r *http.Request) {
	if err := h.controller.LikeArtwork(modellike.CreateInput{
		UserID:    currentuser.FromContext(r.Context()),
		ArtworkID: chi.URLParam(r, "artworkId"),
	}); err != nil {
		writeControllerError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UnlikeArtwork(w http.ResponseWriter, r *http.Request) {
	if err := h.controller.UnlikeArtwork(modellike.DeleteInput{
		UserID:    currentuser.FromContext(r.Context()),
		ArtworkID: chi.URLParam(r, "artworkId"),
	}); err != nil {
		writeControllerError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

func writeControllerError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apperror.ErrConflict):
		writeJSON(w, http.StatusConflict, ErrLikeConflict)
	case errors.Is(err, apperror.ErrNotFound):
		writeJSON(w, http.StatusNotFound, ErrLikeNotFound)
	default:
		writeJSON(w, http.StatusInternalServerError, routererror.ErrInternal)
	}
}
