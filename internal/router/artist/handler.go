package artist

import "artshare/internal/controller/artist"

type Handler struct {
	controller *artist.Controller
}

func NewHandler(controller *artist.Controller) *Handler {
	return &Handler{controller: controller}
}
