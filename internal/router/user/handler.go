package user

import (
	"artshare/internal/controller/user"
	"net/http"
)

type Handler struct {
	controller *user.Controller
}

func NewHandler(controller *user.Controller) *Handler {
	return &Handler{
		controller: controller,
	}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {

	// TODO: Implement logic with controller call

	// распаковуваем ДТО из запроса
}

func (h *Handler) UpdateUser(writer http.ResponseWriter, request *http.Request) {

}

//TODO: Implement other handlers (GetUser, UpdateUser, DeleteUser, etc.)
