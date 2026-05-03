package user

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/mail"
	"strings"

	"github.com/andrqxa/artshare/internal/controller/apperror"
	controlleruser "github.com/andrqxa/artshare/internal/controller/user"
	modeluser "github.com/andrqxa/artshare/internal/model/user"
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

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		writeDecodeError(w, err)
		return
	}

	if errs := validateRegisterRequest(req); len(errs) > 0 {
		writeValidationError(w, errs)
		return
	}

	currentUser, token, err := h.controller.RegisterUser(modeluser.RegisterInput{
		Email:    strings.TrimSpace(req.Email),
		Password: req.Password,
		Role:     modeluser.Role(req.Role),
	})
	if err != nil {
		writeControllerError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, authResponseFromModel(currentUser, token))
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		writeDecodeError(w, err)
		return
	}

	if errs := validateLoginRequest(req); len(errs) > 0 {
		writeValidationError(w, errs)
		return
	}

	currentUser, token, err := h.controller.LoginUser(modeluser.LoginInput{
		Email:    strings.TrimSpace(req.Email),
		Password: req.Password,
	})
	if err != nil {
		writeControllerError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, authResponseFromModel(currentUser, token))
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

func validateRegisterRequest(req RegisterRequest) []FieldError {
	var errs []FieldError
	email := strings.TrimSpace(req.Email)
	if _, err := mail.ParseAddress(email); err != nil {
		errs = append(errs, FieldError{Field: "email", Message: "must be a valid email address"})
	}
	if len(req.Password) < 8 {
		errs = append(errs, FieldError{Field: "password", Message: "must be at least 8 characters"})
	}
	if req.Role != string(modeluser.RoleViewer) && req.Role != string(modeluser.RoleArtist) {
		errs = append(errs, FieldError{Field: "role", Message: "must be viewer or artist"})
	}
	return errs
}

func validateLoginRequest(req LoginRequest) []FieldError {
	var errs []FieldError
	if _, err := mail.ParseAddress(strings.TrimSpace(req.Email)); err != nil {
		errs = append(errs, FieldError{Field: "email", Message: "must be a valid email address"})
	}
	if req.Password == "" {
		errs = append(errs, FieldError{Field: "password", Message: "is required"})
	}
	return errs
}

func writeValidationError(w http.ResponseWriter, details []FieldError) {
	writeJSON(w, http.StatusBadRequest, ErrorResponse{
		Code:    "validation_error",
		Message: "request validation failed",
		Details: details,
	})
}

func writeControllerError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apperror.ErrConflict):
		writeJSON(w, http.StatusConflict, ErrorResponse{Code: "conflict", Message: "request conflicts with current resource state"})
	case errors.Is(err, apperror.ErrUnauthorized):
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Code: "unauthorized", Message: "authentication failed"})
	default:
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Code: "internal_error", Message: "internal server error"})
	}
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

func authResponseFromModel(currentUser modeluser.CurrentUser, token modeluser.AuthToken) AuthResponse {
	return AuthResponse{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		ExpiresIn:   token.ExpiresIn,
		User:        currentUserResponseFromModel(currentUser),
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
