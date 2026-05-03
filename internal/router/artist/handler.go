package artist

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/andrqxa/artshare/internal/controller/artist"
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

	writeNotImplemented(w, "createCurrentArtist")
}

func (h *Handler) GetCurrentArtist(w http.ResponseWriter, r *http.Request) {
	writeNotImplemented(w, "getCurrentArtist")
}

func (h *Handler) UpdateCurrentArtist(w http.ResponseWriter, r *http.Request) {
	var req UpdateArtistRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		writeDecodeError(w, err)
		return
	}

	writeNotImplemented(w, "updateCurrentArtist")
}

func (h *Handler) ListArtists(w http.ResponseWriter, r *http.Request) {
	page, pageSize := paginationFromRequest(r)
	writeJSON(w, http.StatusOK, ArtistListResponse{
		Data: []ArtistSummary{},
		Meta: PaginationMeta{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: 0,
			TotalPages: 0,
		},
	})
}

func (h *Handler) GetArtistByID(w http.ResponseWriter, r *http.Request) {
	_ = chi.URLParam(r, "artistId")
	writeNotImplemented(w, "getArtistById")
}

func paginationFromRequest(r *http.Request) (int, int) {
	page := intQueryParam(r, "page", 1)
	pageSize := intQueryParam(r, "pageSize", 20)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func intQueryParam(r *http.Request, name string, fallback int) int {
	value := r.URL.Query().Get(name)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
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

func writeNotImplemented(w http.ResponseWriter, operation string) {
	writeJSON(w, http.StatusNotImplemented, ErrorResponse{
		Code:    "not_implemented",
		Message: fmt.Sprintf("%s is not implemented yet", operation),
	})
}

func writeJSON(w http.ResponseWriter, status int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}
