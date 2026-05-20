package artwork

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/andrqxa/artshare/internal/controller/apperror"
	"github.com/andrqxa/artshare/internal/controller/artwork"
	"github.com/andrqxa/artshare/internal/middleware/currentuser"
	modelartwork "github.com/andrqxa/artshare/internal/model/artwork"
	"github.com/andrqxa/artshare/internal/model/pagination"
	routererror "github.com/andrqxa/artshare/internal/router/error"
	"github.com/go-chi/chi/v5"
)

const maxRequestBodyBytes = 1 << 20

var errMultipleJSONValues = errors.New("request body must contain a single JSON object")

type Handler struct {
	controller *artwork.Controller
}

func NewHandler(controller *artwork.Controller) *Handler {
	return &Handler{controller: controller}
}

func (h *Handler) ListArtworks(w http.ResponseWriter, r *http.Request) {
	page, pageSize, ok := paginationFromRequest(w, r)
	if !ok {
		return
	}
	status, ok := artworkStatusFromQuery(w, r, "status")
	if !ok {
		return
	}
	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "newest"
	}
	if sort != "newest" && sort != "oldest" {
		writeValidationError(w, []FieldError{{Field: "sort", Message: "must be newest or oldest"}})
		return
	}

	artworks, meta := h.controller.ListArtworks(modelartwork.ListFilter{
		Query:    optionalQuery(r, "q"),
		ArtistID: optionalQuery(r, "artistId"),
		Status:   status,
		Sort:     sort,
		Page:     page,
		PageSize: pageSize,
	})

	writeJSON(w, http.StatusOK, ArtworkListResponse{
		Data: artworkSummariesFromModel(artworks),
		Meta: paginationMetaFromModel(meta),
	})
}

func (h *Handler) CreateArtwork(w http.ResponseWriter, r *http.Request) {
	var req CreateArtworkRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		writeDecodeError(w, err)
		return
	}

	if errs := validateCreateArtworkRequest(req); len(errs) > 0 {
		writeValidationError(w, errs)
		return
	}

	artworkDetail, err := h.controller.CreateArtwork(modelartwork.CreateInput{
		ArtistID:                currentuser.FromContext(r.Context()),
		Title:                   strings.TrimSpace(req.Title),
		Description:             req.Description,
		Tags:                    req.Tags,
		Images:                  imagesToModel(req.Images),
		AvailabilityForExchange: req.AvailabilityForExchange,
	})
	if err != nil {
		writeControllerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, artworkDetailFromModel(artworkDetail))
}

func (h *Handler) GetArtworkByID(w http.ResponseWriter, r *http.Request) {
	artworkDetail, err := h.controller.GetArtworkByID(chi.URLParam(r, "artworkId"))
	if err != nil {
		writeControllerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, artworkDetailFromModel(artworkDetail))
}

func (h *Handler) UpdateArtwork(w http.ResponseWriter, r *http.Request) {
	var req UpdateArtworkRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		writeDecodeError(w, err)
		return
	}

	if errs := validateUpdateArtworkRequest(req); len(errs) > 0 {
		writeValidationError(w, errs)
		return
	}

	status := (*modelartwork.Status)(nil)
	if req.Status != nil {
		parsed := modelartwork.Status(*req.Status)
		status = &parsed
	}
	artworkDetail, err := h.controller.UpdateArtwork(modelartwork.UpdateInput{
		ID:                      chi.URLParam(r, "artworkId"),
		Title:                   trimmedString(req.Title),
		Description:             req.Description,
		Tags:                    req.Tags,
		Images:                  imagesToModel(req.Images),
		AvailabilityForExchange: req.AvailabilityForExchange,
		Status:                  status,
	})
	if err != nil {
		writeControllerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, artworkDetailFromModel(artworkDetail))
}

func (h *Handler) DeleteArtwork(w http.ResponseWriter, r *http.Request) {
	if err := h.controller.DeleteArtwork(chi.URLParam(r, "artworkId")); err != nil {
		writeControllerError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
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

func validateCreateArtworkRequest(req CreateArtworkRequest) []FieldError {
	var errs []FieldError
	if strings.TrimSpace(req.Title) == "" || len(req.Title) > 200 {
		errs = append(errs, FieldError{Field: "title", Message: "must be between 1 and 200 characters"})
	}
	if len(req.Description) > 5000 {
		errs = append(errs, FieldError{Field: "description", Message: "must be at most 5000 characters"})
	}
	if len(req.Images) < 1 || len(req.Images) > 10 {
		errs = append(errs, FieldError{Field: "images", Message: "must contain between 1 and 10 images"})
	}
	if len(req.Tags) > 20 {
		errs = append(errs, FieldError{Field: "tags", Message: "must contain at most 20 tags"})
	}
	return errs
}

func validateUpdateArtworkRequest(req UpdateArtworkRequest) []FieldError {
	var errs []FieldError
	if req.Title != nil && (strings.TrimSpace(*req.Title) == "" || len(*req.Title) > 200) {
		errs = append(errs, FieldError{Field: "title", Message: "must be between 1 and 200 characters"})
	}
	if req.Description != nil && len(*req.Description) > 5000 {
		errs = append(errs, FieldError{Field: "description", Message: "must be at most 5000 characters"})
	}
	if req.Images != nil && (len(req.Images) < 1 || len(req.Images) > 10) {
		errs = append(errs, FieldError{Field: "images", Message: "must contain between 1 and 10 images"})
	}
	if len(req.Tags) > 20 {
		errs = append(errs, FieldError{Field: "tags", Message: "must contain at most 20 tags"})
	}
	if req.Status != nil && *req.Status != string(modelartwork.StatusDraft) && *req.Status != string(modelartwork.StatusPublished) {
		errs = append(errs, FieldError{Field: "status", Message: "must be draft or published"})
	}
	return errs
}

func artworkStatusFromQuery(w http.ResponseWriter, r *http.Request, name string) (*modelartwork.Status, bool) {
	value := r.URL.Query().Get(name)
	if value == "" {
		return nil, true
	}
	if value != string(modelartwork.StatusDraft) && value != string(modelartwork.StatusPublished) {
		writeValidationError(w, []FieldError{{Field: name, Message: "must be draft or published"}})
		return nil, false
	}
	status := modelartwork.Status(value)
	return &status, true
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
	case errors.Is(err, apperror.ErrNotFound):
		writeJSON(w, http.StatusNotFound, ErrArtworkNotFound)
	case errors.Is(err, apperror.ErrForbidden):
		writeJSON(w, http.StatusForbidden, ErrArtworkForbidden)
	default:
		writeJSON(w, http.StatusInternalServerError, routererror.ErrInternal)
	}
}

func imagesToModel(images []ArtworkImage) []modelartwork.Image {
	if images == nil {
		return nil
	}
	out := make([]modelartwork.Image, 0, len(images))
	for _, image := range images {
		out = append(out, modelartwork.Image{URL: image.URL, AltText: image.AltText})
	}
	return out
}

func imagesFromModel(images []modelartwork.Image) []ArtworkImage {
	out := make([]ArtworkImage, 0, len(images))
	for _, image := range images {
		out = append(out, ArtworkImage{URL: image.URL, AltText: image.AltText})
	}
	return out
}

func artworkSummariesFromModel(artworks []modelartwork.Summary) []ArtworkSummary {
	out := make([]ArtworkSummary, 0, len(artworks))
	for _, artwork := range artworks {
		out = append(out, ArtworkSummary{
			ID:                      artwork.ID,
			Title:                   artwork.Title,
			Status:                  string(artwork.Status),
			Owner:                   artworkOwnerFromModel(artwork.Owner),
			PreviewImageURL:         artwork.PreviewImageURL,
			LikeCount:               artwork.LikeCount,
			LikedByMe:               artwork.LikedByMe,
			AvailabilityForExchange: artwork.AvailabilityForExchange,
			CreatedAt:               artwork.CreatedAt,
			PublishedAt:             artwork.PublishedAt,
		})
	}
	return out
}

func artworkDetailFromModel(artwork modelartwork.Detail) ArtworkDetail {
	return ArtworkDetail{
		ID:                      artwork.ID,
		Title:                   artwork.Title,
		Status:                  string(artwork.Status),
		Owner:                   artworkOwnerFromModel(artwork.Owner),
		PreviewImageURL:         artwork.PreviewImageURL,
		LikeCount:               artwork.LikeCount,
		LikedByMe:               artwork.LikedByMe,
		AvailabilityForExchange: artwork.AvailabilityForExchange,
		CreatedAt:               artwork.CreatedAt,
		PublishedAt:             artwork.PublishedAt,
		Description:             artwork.Description,
		Tags:                    artwork.Tags,
		Images:                  imagesFromModel(artwork.Images),
		UpdatedAt:               artwork.UpdatedAt,
	}
}

func artworkOwnerFromModel(owner modelartwork.Owner) ArtworkOwner {
	return ArtworkOwner{ID: owner.ID, DisplayName: owner.DisplayName, AvatarURL: owner.AvatarURL}
}

func paginationMetaFromModel(meta pagination.Meta) PaginationMeta {
	return PaginationMeta{Page: meta.Page, PageSize: meta.PageSize, TotalItems: meta.TotalItems, TotalPages: meta.TotalPages}
}
