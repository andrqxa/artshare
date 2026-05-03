package like

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andrqxa/artshare/internal/controller/like"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	controller *like.Controller
}

func NewHandler(controller *like.Controller) *Handler {
	return &Handler{controller: controller}
}

func (h *Handler) LikeArtwork(w http.ResponseWriter, r *http.Request) {
	_ = chi.URLParam(r, "artworkId")
	writeNotImplemented(w, "likeArtwork")
}

func (h *Handler) UnlikeArtwork(w http.ResponseWriter, r *http.Request) {
	_ = chi.URLParam(r, "artworkId")
	writeNotImplemented(w, "unlikeArtwork")
}

func writeNotImplemented(w http.ResponseWriter, operation string) {
	writeJSON(w, http.StatusNotImplemented, ErrorResponse{
		Code:    "not_implemented",
		Message: fmt.Sprintf("%s is not implemented yet", operation),
	})
}

func writeJSON(w http.ResponseWriter, status int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}
