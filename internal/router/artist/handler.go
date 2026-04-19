package artist

type Handler struct {
	controller *artist.Controller
}

func NewHandler(controller *artist.Controller) *Handler {
	return &Handler{controller: controller}
}
