package artwork

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	modelartwork "github.com/andrqxa/artshare/internal/model/artwork"
	"github.com/andrqxa/artshare/internal/model/pagination"
	"github.com/andrqxa/artshare/internal/repository/repoerror"
)

type Repository struct {
	mu     sync.RWMutex
	nextID int
	byID   map[string]modelartwork.Artwork
}

func NewRepository() *Repository {
	return &Repository{
		nextID: 1,
		byID:   make(map[string]modelartwork.Artwork),
	}
}

func (r *Repository) List(filter modelartwork.ListFilter) ([]modelartwork.Summary, pagination.Meta) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	page, pageSize := normalizePagination(filter.Page, filter.PageSize)
	query := strings.TrimSpace(strings.ToLower(valueOrEmpty(filter.Query)))

	items := make([]modelartwork.Summary, 0, len(r.byID))
	for _, artwork := range r.byID {
		if query != "" && !strings.Contains(strings.ToLower(artwork.Title), query) {
			continue
		}
		if filter.ArtistID != nil && artwork.ArtistID != *filter.ArtistID {
			continue
		}
		if filter.Status != nil && artwork.Status != *filter.Status {
			continue
		}
		items = append(items, summaryFromArtwork(artwork))
	}
	sort.Slice(items, func(i, j int) bool {
		if filter.Sort == "oldest" {
			return items[i].CreatedAt.Before(items[j].CreatedAt)
		}
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})

	total := len(items)
	items = pageSlice(items, page, pageSize)
	return items, paginationMeta(page, pageSize, total)
}

func (r *Repository) Create(input modelartwork.CreateInput) (modelartwork.Detail, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	artwork := modelartwork.Artwork{
		ID:                      nextUUID(r.nextID),
		ArtistID:                input.ArtistID,
		Title:                   input.Title,
		Description:             input.Description,
		Status:                  modelartwork.StatusDraft,
		Tags:                    input.Tags,
		Images:                  input.Images,
		PreviewImageURL:         previewImageURL(input.Images),
		AvailabilityForExchange: input.AvailabilityForExchange,
		CreatedAt:               now,
		UpdatedAt:               now,
	}
	r.nextID++
	r.byID[artwork.ID] = artwork

	return detailFromArtwork(artwork), nil
}

func (r *Repository) FindByID(id string) (modelartwork.Detail, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	artwork, exists := r.byID[id]
	if !exists {
		return modelartwork.Detail{}, repoerror.ErrNotFound
	}
	return detailFromArtwork(artwork), nil
}

func (r *Repository) Update(input modelartwork.UpdateInput) (modelartwork.Detail, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	artwork, exists := r.byID[input.ID]
	if !exists {
		return modelartwork.Detail{}, repoerror.ErrNotFound
	}
	if input.Title != nil {
		artwork.Title = *input.Title
	}
	if input.Description != nil {
		artwork.Description = *input.Description
	}
	if input.Tags != nil {
		artwork.Tags = input.Tags
	}
	if input.Images != nil {
		artwork.Images = input.Images
		artwork.PreviewImageURL = previewImageURL(input.Images)
	}
	if input.AvailabilityForExchange != nil {
		artwork.AvailabilityForExchange = *input.AvailabilityForExchange
	}
	if input.Status != nil {
		artwork.Status = *input.Status
		if *input.Status == modelartwork.StatusPublished && artwork.PublishedAt == nil {
			now := time.Now().UTC()
			artwork.PublishedAt = &now
		}
	}
	artwork.UpdatedAt = time.Now().UTC()
	r.byID[artwork.ID] = artwork

	return detailFromArtwork(artwork), nil
}

func (r *Repository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byID[id]; !exists {
		return repoerror.ErrNotFound
	}
	delete(r.byID, id)
	return nil
}

func summaryFromArtwork(artwork modelartwork.Artwork) modelartwork.Summary {
	return modelartwork.Summary{
		ID:                      artwork.ID,
		Title:                   artwork.Title,
		Status:                  artwork.Status,
		Owner:                   ownerFromArtwork(artwork),
		PreviewImageURL:         artwork.PreviewImageURL,
		LikeCount:               artwork.LikeCount,
		LikedByMe:               artwork.LikedByMe,
		AvailabilityForExchange: artwork.AvailabilityForExchange,
		CreatedAt:               artwork.CreatedAt,
		PublishedAt:             artwork.PublishedAt,
	}
}

func detailFromArtwork(artwork modelartwork.Artwork) modelartwork.Detail {
	return modelartwork.Detail{
		ID:                      artwork.ID,
		Title:                   artwork.Title,
		Description:             artwork.Description,
		Status:                  artwork.Status,
		Owner:                   ownerFromArtwork(artwork),
		PreviewImageURL:         artwork.PreviewImageURL,
		LikeCount:               artwork.LikeCount,
		LikedByMe:               artwork.LikedByMe,
		AvailabilityForExchange: artwork.AvailabilityForExchange,
		Tags:                    artwork.Tags,
		Images:                  artwork.Images,
		CreatedAt:               artwork.CreatedAt,
		UpdatedAt:               artwork.UpdatedAt,
		PublishedAt:             artwork.PublishedAt,
	}
}

func ownerFromArtwork(artwork modelartwork.Artwork) modelartwork.Owner {
	return modelartwork.Owner{ID: artwork.ArtistID, DisplayName: "Artist"}
}

func previewImageURL(images []modelartwork.Image) *string {
	if len(images) == 0 {
		return nil
	}
	return &images[0].URL
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func normalizePagination(page, pageSize int) (int, int) {
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

func paginationMeta(page, pageSize, totalItems int) pagination.Meta {
	totalPages := 0
	if totalItems > 0 {
		totalPages = (totalItems + pageSize - 1) / pageSize
	}
	return pagination.Meta{Page: page, PageSize: pageSize, TotalItems: totalItems, TotalPages: totalPages}
}

func pageSlice[T any](items []T, page, pageSize int) []T {
	start := (page - 1) * pageSize
	if start >= len(items) {
		return []T{}
	}
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}

func nextUUID(id int) string {
	return fmt.Sprintf("00000000-0000-0000-0000-%012d", id)
}
