package user

import (
	"artshare/internal/controller/user"
	"encoding/json"
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
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	h.controller.CreateUser()

	writeJSON(w, http.StatusCreated, CreateUserResponse{
		TokenType: "Bearer",
		User: User{
			Email: req.Email,
			Role:  req.Role,
		},
	})
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, User{})
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	writeJSON(w, http.StatusOK, User{
		Email: req.Email,
		Role:  req.Role,
	})
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}
