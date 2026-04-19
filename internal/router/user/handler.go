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
	if !hasBearerToken(r) {
		writeError(w, http.StatusUnauthorized, "unauthorized", "authentication is required or the token is invalid")
		return
	}

	writeJSON(w, http.StatusOK, toCurrentUserResponse(h.controller.GetCurrentUser()))
}

func (h *Handler) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	if !hasBearerToken(r) {
		writeError(w, http.StatusUnauthorized, "unauthorized", "authentication is required or the token is invalid")
		return
	}

	var req UpdateCurrentUserRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		writeRequestDecodeError(w, err)
		return
	}

	if req.Email != nil {
		email := strings.TrimSpace(*req.Email)
		if email == "" || !isValidEmail(email) {
			writeError(
				w,
				http.StatusBadRequest,
				"validation_error",
				"the request payload is invalid",
				FieldError{Field: "email", Message: "must be a valid email address"},
			)
			return
		}

		req.Email = &email
	}

	currentUser := h.controller.UpdateCurrentUser(modeluser.UpdateCurrentUserInput{
		Email: req.Email,
	})

	writeJSON(w, http.StatusOK, toCurrentUserResponse(currentUser))
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) error {
	if contentType := r.Header.Get("Content-Type"); contentType != "" && !strings.HasPrefix(contentType, "application/json") {
		return errors.New("content type must be application/json")
	}

	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxRequestBodyBytes))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return err
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return errMultipleJSONValues
	}

	return nil
}

func hasBearerToken(r *http.Request) bool {
	parts := strings.Fields(r.Header.Get("Authorization"))
	return len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") && parts[1] != ""
}

func isValidEmail(email string) bool {
	address, err := mail.ParseAddress(email)
	return err == nil && address.Address == email
}

func toCurrentUserResponse(currentUser modeluser.CurrentUser) CurrentUserResponse {
	response := CurrentUserResponse{
		ID:            currentUser.ID,
		Email:         currentUser.Email,
		Role:          string(currentUser.Role),
		CreatedAt:     currentUser.CreatedAt,
		ArtistProfile: nil,
	}

	if currentUser.ArtistProfile != nil {
		response.ArtistProfile = &ArtistProfileResponse{
			ID:          currentUser.ArtistProfile.ID,
			UserID:      currentUser.ArtistProfile.UserID,
			DisplayName: currentUser.ArtistProfile.DisplayName,
			Bio:         currentUser.ArtistProfile.Bio,
			AvatarURL:   currentUser.ArtistProfile.AvatarURL,
			CreatedAt:   currentUser.ArtistProfile.CreatedAt,
			UpdatedAt:   currentUser.ArtistProfile.UpdatedAt,
		}
	}

	return response
}

func writeRequestDecodeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, io.EOF):
		writeError(w, http.StatusBadRequest, "validation_error", "request body is required")
	case errors.Is(err, errMultipleJSONValues):
		writeError(w, http.StatusBadRequest, "validation_error", "request body must contain a single JSON object")
	default:
		writeError(w, http.StatusBadRequest, "validation_error", "the request payload is invalid")
	}
}

func writeError(w http.ResponseWriter, status int, code, message string, details ...FieldError) {
	response := ErrorResponse{
		Code:    code,
		Message: message,
	}

	if len(details) > 0 {
		response.Details = details
	}

	writeJSON(w, status, response)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
