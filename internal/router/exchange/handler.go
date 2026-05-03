package exchange

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/andrqxa/artshare/internal/controller/apperror"
	"github.com/andrqxa/artshare/internal/controller/exchange"
	modelartwork "github.com/andrqxa/artshare/internal/model/artwork"
	modelexchange "github.com/andrqxa/artshare/internal/model/exchange"
	"github.com/andrqxa/artshare/internal/model/pagination"
	"github.com/go-chi/chi/v5"
)

const maxRequestBodyBytes = 1 << 20

var errMultipleJSONValues = errors.New("request body must contain a single JSON object")

type Handler struct {
	controller *exchange.Controller
}

func NewHandler(controller *exchange.Controller) *Handler {
	return &Handler{controller: controller}
}

func (h *Handler) ListExchangeRequests(w http.ResponseWriter, r *http.Request) {
	page, pageSize, ok := paginationFromRequest(w, r)
	if !ok {
		return
	}
	scope, ok := exchangeScopeFromQuery(w, r)
	if !ok {
		return
	}
	status, ok := exchangeStatusFromQuery(w, r, "status", true)
	if !ok {
		return
	}

	requests, meta := h.controller.ListExchangeRequests(modelexchange.ListFilter{
		UserID:   "00000000-0000-0000-0000-000000000001",
		Scope:    scope,
		Status:   status,
		Page:     page,
		PageSize: pageSize,
	})

	writeJSON(w, http.StatusOK, ExchangeRequestListResponse{
		Data: exchangeRequestSummariesFromModel(requests),
		Meta: paginationMetaFromModel(meta),
	})
}

func (h *Handler) CreateExchangeRequest(w http.ResponseWriter, r *http.Request) {
	var req CreateExchangeRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		writeDecodeError(w, err)
		return
	}

	if req.ArtworkID == "" {
		writeValidationError(w, []FieldError{{Field: "artworkId", Message: "is required"}})
		return
	}

	request, err := h.controller.CreateExchangeRequest(modelexchange.CreateInput{
		ArtworkID:   req.ArtworkID,
		RequesterID: "00000000-0000-0000-0000-000000000001",
		Message:     req.Message,
	})
	if err != nil {
		writeControllerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, exchangeRequestDetailFromModel(request))
}

func (h *Handler) GetExchangeRequestByID(w http.ResponseWriter, r *http.Request) {
	request, err := h.controller.GetExchangeRequestByID(chi.URLParam(r, "exchangeRequestId"))
	if err != nil {
		writeControllerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, exchangeRequestDetailFromModel(request))
}

func (h *Handler) UpdateExchangeRequestStatus(w http.ResponseWriter, r *http.Request) {
	var req UpdateExchangeRequestStatusRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		writeDecodeError(w, err)
		return
	}

	status, ok := exchangeStatusFromValue(w, req.Status, "status", false)
	if !ok {
		return
	}
	request, err := h.controller.UpdateExchangeRequestStatus(modelexchange.UpdateStatusInput{
		ID:      chi.URLParam(r, "exchangeRequestId"),
		Status:  *status,
		Message: req.Message,
	})
	if err != nil {
		writeControllerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, exchangeRequestDetailFromModel(request))
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
	message := "invalid request body"
	if errors.Is(err, errMultipleJSONValues) {
		message = err.Error()
	}

	writeJSON(w, http.StatusBadRequest, ErrorResponse{
		Code:    "bad_request",
		Message: message,
	})
}

func writeJSON(w http.ResponseWriter, status int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

func exchangeScopeFromQuery(w http.ResponseWriter, r *http.Request) (modelexchange.Scope, bool) {
	value := r.URL.Query().Get("scope")
	if value == "" {
		return modelexchange.ScopeAll, true
	}
	switch modelexchange.Scope(value) {
	case modelexchange.ScopeInbox, modelexchange.ScopeOutbox, modelexchange.ScopeAll:
		return modelexchange.Scope(value), true
	default:
		writeValidationError(w, []FieldError{{Field: "scope", Message: "must be inbox, outbox, or all"}})
		return "", false
	}
}

func exchangeStatusFromQuery(w http.ResponseWriter, r *http.Request, name string, allowPending bool) (*modelexchange.Status, bool) {
	value := r.URL.Query().Get(name)
	if value == "" {
		return nil, true
	}
	return exchangeStatusFromValue(w, value, name, allowPending)
}

func exchangeStatusFromValue(w http.ResponseWriter, value string, field string, allowPending bool) (*modelexchange.Status, bool) {
	status := modelexchange.Status(value)
	switch status {
	case modelexchange.StatusAccepted, modelexchange.StatusRejected, modelexchange.StatusCancelled:
		return &status, true
	case modelexchange.StatusPending:
		if allowPending {
			return &status, true
		}
	}
	if allowPending {
		writeValidationError(w, []FieldError{{Field: field, Message: "must be pending, accepted, rejected, or cancelled"}})
	} else {
		writeValidationError(w, []FieldError{{Field: field, Message: "must be accepted, rejected, or cancelled"}})
	}
	return nil, false
}

func writeValidationError(w http.ResponseWriter, details []FieldError) {
	writeJSON(w, http.StatusBadRequest, ErrorResponse{Code: "validation_error", Message: "request validation failed", Details: details})
}

func writeControllerError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apperror.ErrConflict):
		writeJSON(w, http.StatusConflict, ErrorResponse{Code: "conflict", Message: "request conflicts with current resource state"})
	case errors.Is(err, apperror.ErrNotFound):
		writeJSON(w, http.StatusNotFound, ErrorResponse{Code: "not_found", Message: "requested resource was not found"})
	case errors.Is(err, apperror.ErrForbidden):
		writeJSON(w, http.StatusForbidden, ErrorResponse{Code: "forbidden", Message: "permission denied"})
	default:
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Code: "internal_error", Message: "internal server error"})
	}
}

func exchangeRequestSummariesFromModel(requests []modelexchange.Summary) []ExchangeRequestSummary {
	out := make([]ExchangeRequestSummary, 0, len(requests))
	for _, request := range requests {
		out = append(out, ExchangeRequestSummary{
			ID:          request.ID,
			Artwork:     artworkRefFromModel(request.Artwork),
			RequesterID: request.RequesterID,
			OwnerID:     request.OwnerID,
			Status:      string(request.Status),
			CreatedAt:   request.CreatedAt,
			UpdatedAt:   request.UpdatedAt,
		})
	}
	return out
}

func exchangeRequestDetailFromModel(request modelexchange.Detail) ExchangeRequestDetail {
	return ExchangeRequestDetail{
		ID:          request.ID,
		Artwork:     artworkRefFromModel(request.Artwork),
		RequesterID: request.RequesterID,
		OwnerID:     request.OwnerID,
		Status:      string(request.Status),
		CreatedAt:   request.CreatedAt,
		UpdatedAt:   request.UpdatedAt,
		Message:     request.Message,
	}
}

func artworkRefFromModel(artwork modelartwork.OwneredRef) ArtworkOwneredRef {
	return ArtworkOwneredRef{ID: artwork.ID, Title: artwork.Title, OwnerID: artwork.OwnerID, PreviewImageURL: artwork.PreviewImageURL}
}

func paginationMetaFromModel(meta pagination.Meta) PaginationMeta {
	return PaginationMeta{Page: meta.Page, PageSize: meta.PageSize, TotalItems: meta.TotalItems, TotalPages: meta.TotalPages}
}
