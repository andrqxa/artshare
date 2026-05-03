package like

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/andrqxa/artshare/internal/controller/apperror"
	"github.com/andrqxa/artshare/internal/controller/like"
	modellike "github.com/andrqxa/artshare/internal/model/like"
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
		UserID:    "00000000-0000-0000-0000-000000000001",
		ArtworkID: chi.URLParam(r, "artworkId"),
	}); err != nil {
		writeControllerError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UnlikeArtwork(w http.ResponseWriter, r *http.Request) {
	if err := h.controller.UnlikeArtwork(modellike.DeleteInput{
		UserID:    "00000000-0000-0000-0000-000000000001",
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
		writeJSON(w, http.StatusConflict, ErrorResponse{Code: "conflict", Message: "request conflicts with current resource state"})
	case errors.Is(err, apperror.ErrNotFound):
		writeJSON(w, http.StatusNotFound, ErrorResponse{Code: "not_found", Message: "requested resource was not found"})
	default:
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Code: "internal_error", Message: "internal server error"})
	}
}
