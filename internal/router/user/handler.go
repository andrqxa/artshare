package user

import (
<<<<<<< HEAD
	controlleruser "artshare/internal/controller/user"
	modeluser "artshare/internal/model/user"
	"encoding/json"
	"errors"
	"io"
=======
	"artshare/internal/controller/user"
	"encoding/json"
>>>>>>> develop.bak
	"net/http"
	"net/mail"
	"strings"
)

const maxRequestBodyBytes = 1 << 20

var errMultipleJSONValues = errors.New("request body must contain a single JSON object")

type Handler struct {
	controller *controlleruser.Controller
}

func NewHandler(controller *controlleruser.Controller) *Handler {
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
