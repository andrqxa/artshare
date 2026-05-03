package user

import (
	controlleruser "artshare/internal/controller/user"
	modeluser "artshare/internal/model/user"
	"encoding/json"
	"errors"
	"io"
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

func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, currentUserResponseFromModel(h.controller.GetCurrentUser()))
}

func (h *Handler) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	var req UpdateCurrentUserRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		writeDecodeError(w, err)
		return
	}

	if req.Email != nil {
		email := strings.TrimSpace(*req.Email)
		if _, err := mail.ParseAddress(email); err != nil {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{
				Code:    "validation_error",
				Message: "request validation failed",
				Details: []FieldError{
					{Field: "email", Message: "must be a valid email address"},
				},
			})
			return
		}
		req.Email = &email
	}

	currentUser := h.controller.UpdateCurrentUser(modeluser.UpdateCurrentUserInput{
		Email: req.Email,
	})
	writeJSON(w, http.StatusOK, currentUserResponseFromModel(currentUser))
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return err
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errMultipleJSONValues
	}

	return nil
}

func writeDecodeError(w http.ResponseWriter, err error) {
	message := "invalid request body"
	if errors.Is(err, errMultipleJSONValues) {
		message = err.Error()
	}

	writeJSON(w, http.StatusBadRequest, ErrorResponse{
		Code:    "bad_request",
		Message: message,
	})
}

func currentUserResponseFromModel(currentUser modeluser.CurrentUser) CurrentUserResponse {
	return CurrentUserResponse{
		ID:        currentUser.ID,
		Email:     currentUser.Email,
		Role:      string(currentUser.Role),
		CreatedAt: currentUser.CreatedAt,
		Artist:    artistFromModel(currentUser.ArtistProfile),
	}
}

func artistFromModel(artist *modeluser.ArtistProfile) *Artist {
	if artist == nil {
		return nil
	}

	return &Artist{
		ID:          artist.ID,
		UserID:      artist.UserID,
		DisplayName: artist.DisplayName,
		Bio:         artist.Bio,
		AvatarURL:   artist.AvatarURL,
		CreatedAt:   artist.CreatedAt,
		UpdatedAt:   artist.UpdatedAt,
	}
}

func writeJSON(w http.ResponseWriter, status int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}
