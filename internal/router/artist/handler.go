package artist

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/andrqxa/artshare/internal/controller/apperror"
	"github.com/andrqxa/artshare/internal/controller/artist"
	"github.com/andrqxa/artshare/internal/middleware/currentuser"
	modelartist "github.com/andrqxa/artshare/internal/model/artist"
	"github.com/andrqxa/artshare/internal/model/pagination"
	routererror "github.com/andrqxa/artshare/internal/router/error"
	"github.com/go-chi/chi/v5"
)

const maxRequestBodyBytes = 1 << 20

var errMultipleJSONValues = errors.New("request body must contain a single JSON object")

type Handler struct {
	controller *artist.Controller
}

func NewHandler(controller *artist.Controller) *Handler {
	return &Handler{controller: controller}
}

func (h *Handler) CreateCurrentArtist(w http.ResponseWriter, r *http.Request) {
	var req CreateArtistRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		writeDecodeError(w, err)
		return
	}

	if errs := validateCreateArtistRequest(req); len(errs) > 0 {
		writeValidationError(w, errs)
		return
	}

	artistModel, err := h.controller.CreateCurrentArtist(modelartist.CreateInput{
		UserID:      currentuser.FromContext(r.Context()),
		DisplayName: strings.TrimSpace(req.DisplayName),
		Bio:         req.Bio,
		AvatarURL:   req.AvatarURL,
	})
	if err != nil {
		writeControllerError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, artistDetailFromModel(artistModel))
}

func (h *Handler) GetCurrentArtist(w http.ResponseWriter, r *http.Request) {
	artistModel, err := h.controller.GetCurrentArtist(currentuser.FromContext(r.Context()))
	if err != nil {
		writeControllerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, artistDetailFromModel(artistModel))
}

func (h *Handler) UpdateCurrentArtist(w http.ResponseWriter, r *http.Request) {
	var req UpdateArtistRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		writeDecodeError(w, err)
		return
	}

	if errs := validateUpdateArtistRequest(req); len(errs) > 0 {
		writeValidationError(w, errs)
		return
	}

	artistModel, err := h.controller.UpdateCurrentArtist(currentuser.FromContext(r.Context()), modelartist.UpdateInput{
		DisplayName: trimmedString(req.DisplayName),
		Bio:         req.Bio,
		AvatarURL:   req.AvatarURL,
	})
	if err != nil {
		writeControllerError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, artistDetailFromModel(artistModel))
}

func (h *Handler) ListArtists(w http.ResponseWriter, r *http.Request) {
	page, pageSize, ok := paginationFromRequest(w, r)
	if !ok {
		return
	}
	query := optionalQuery(r, "q")
	artists, meta := h.controller.ListArtists(modelartist.ListFilter{
		Query:    query,
		Page:     page,
		PageSize: pageSize,
	})

	writeJSON(w, http.StatusOK, ArtistListResponse{
		Data: artistSummariesFromModel(artists),
		Meta: paginationMetaFromModel(meta),
	})
}

func (h *Handler) GetArtistByID(w http.ResponseWriter, r *http.Request) {
	artistModel, err := h.controller.GetArtistByID(chi.URLParam(r, "artistId"))
	if err != nil {
		writeControllerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, artistDetailFromModel(artistModel))
}

func paginationFromRequest(w http.ResponseWriter, r *http.Request) (int, int, bool) {
	page, err := intQueryParam(r, "page", 1)
	if err != nil || page < 1 {
		writeValidationError(w, []FieldError{{Field: "page", Message: "must be an integer greater than or equal to 1"}})
		return 0, 0, false
	}
	pageSize, err := intQueryParam(r, "pageSize", 20)
	if err != nil || pageSize < 1 {
		writeValidationError(w, []FieldError{{Field: "pageSize", Message: "must be an integer between 1 and 100"}})
		return 0, 0, false
	}
	if pageSize > 100 {
		writeValidationError(w, []FieldError{{Field: "pageSize", Message: "must be an integer between 1 and 100"}})
		return 0, 0, false
	}
	return page, pageSize, true
}

func intQueryParam(r *http.Request, name string, fallback int) (int, error) {
	value := r.URL.Query().Get(name)
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return parsed, nil
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
	if errors.Is(err, errMultipleJSONValues) {
		writeJSON(w, http.StatusBadRequest, routererror.ErrMultipleJSONValues)
		return
	}
	writeJSON(w, http.StatusBadRequest, routererror.ErrBadRequest)
}

func writeJSON(w http.ResponseWriter, status int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

func validateCreateArtistRequest(req CreateArtistRequest) []FieldError {
	displayName := strings.TrimSpace(req.DisplayName)
	if len(displayName) < 2 || len(displayName) > 80 {
		return []FieldError{{Field: "displayName", Message: "must be between 2 and 80 characters"}}
	}
	return nil
}

func validateUpdateArtistRequest(req UpdateArtistRequest) []FieldError {
	if req.DisplayName != nil {
		displayName := strings.TrimSpace(*req.DisplayName)
		if len(displayName) < 2 || len(displayName) > 80 {
			return []FieldError{{Field: "displayName", Message: "must be between 2 and 80 characters"}}
		}
	}
	return nil
}

func trimmedString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	return &trimmed
}

func optionalQuery(r *http.Request, name string) *string {
	value := strings.TrimSpace(r.URL.Query().Get(name))
	if value == "" {
		return nil
	}
	return &value
}

func writeValidationError(w http.ResponseWriter, details []FieldError) {
	writeJSON(w, http.StatusBadRequest, routererror.Validation(details))
}

func writeControllerError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apperror.ErrConflict):
		writeJSON(w, http.StatusConflict, ErrArtistConflict)
	case errors.Is(err, apperror.ErrNotFound):
		writeJSON(w, http.StatusNotFound, ErrArtistNotFound)
	default:
		writeJSON(w, http.StatusInternalServerError, routererror.ErrInternal)
	}
}

func artistSummariesFromModel(artists []modelartist.Summary) []ArtistSummary {
	out := make([]ArtistSummary, 0, len(artists))
	for _, artist := range artists {
		out = append(out, ArtistSummary{
			ID:           artist.ID,
			DisplayName:  artist.DisplayName,
			AvatarURL:    artist.AvatarURL,
			Bio:          artist.Bio,
			ArtworkCount: artist.ArtworkCount,
		})
	}
	return out
}

func artistDetailFromModel(artist modelartist.Detail) ArtistDetail {
	return ArtistDetail{
		ID:           artist.ID,
		DisplayName:  artist.DisplayName,
		AvatarURL:    artist.AvatarURL,
		Bio:          artist.Bio,
		ArtworkCount: artist.ArtworkCount,
		UserID:       artist.UserID,
		CreatedAt:    artist.CreatedAt,
		UpdatedAt:    artist.UpdatedAt,
	}
}

func paginationMetaFromModel(meta pagination.Meta) PaginationMeta {
	return PaginationMeta{Page: meta.Page, PageSize: meta.PageSize, TotalItems: meta.TotalItems, TotalPages: meta.TotalPages}
}
