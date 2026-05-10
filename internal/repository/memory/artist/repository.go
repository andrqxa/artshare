package artist

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	modelartist "github.com/andrqxa/artshare/internal/model/artist"
	"github.com/andrqxa/artshare/internal/model/pagination"
	"github.com/andrqxa/artshare/internal/repository/repoerror"
)

type Repository struct {
	mu       sync.RWMutex
	nextID   int
	byID     map[string]modelartist.Detail
	byUserID map[string]string
}

func NewRepository() *Repository {
	return &Repository{
		nextID:   1,
		byID:     make(map[string]modelartist.Detail),
		byUserID: make(map[string]string),
	}
}

func (r *Repository) Create(input modelartist.CreateInput) (modelartist.Detail, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byUserID[input.UserID]; exists {
		return modelartist.Detail{}, repoerror.ErrConflict
	}

	now := time.Now().UTC()
	artist := modelartist.Detail{
		ID:          nextUUID(r.nextID),
		UserID:      input.UserID,
		DisplayName: input.DisplayName,
		Bio:         input.Bio,
		AvatarURL:   input.AvatarURL,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	r.nextID++
	r.byID[artist.ID] = artist
	r.byUserID[artist.UserID] = artist.ID

	return artist, nil
}

func (r *Repository) FindByUserID(userID string) (modelartist.Detail, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	artistID, exists := r.byUserID[userID]
	if !exists {
		return modelartist.Detail{}, repoerror.ErrNotFound
	}
	return r.byID[artistID], nil
}

func (r *Repository) FindByID(id string) (modelartist.Detail, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	artist, exists := r.byID[id]
	if !exists {
		return modelartist.Detail{}, repoerror.ErrNotFound
	}
	return artist, nil
}

func (r *Repository) UpdateByUserID(userID string, input modelartist.UpdateInput) (modelartist.Detail, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	artistID, exists := r.byUserID[userID]
	if !exists {
		return modelartist.Detail{}, repoerror.ErrNotFound
	}

	artist := r.byID[artistID]
	if input.DisplayName != nil {
		artist.DisplayName = *input.DisplayName
	}
	if input.Bio != nil {
		artist.Bio = input.Bio
	}
	if input.AvatarURL != nil {
		artist.AvatarURL = input.AvatarURL
	}
	artist.UpdatedAt = time.Now().UTC()
	r.byID[artist.ID] = artist

	return artist, nil
}

func (r *Repository) List(filter modelartist.ListFilter) ([]modelartist.Summary, pagination.Meta) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	page, pageSize := normalizePagination(filter.Page, filter.PageSize)
	query := strings.TrimSpace(strings.ToLower(valueOrEmpty(filter.Query)))

	items := make([]modelartist.Summary, 0, len(r.byID))
	for _, artist := range r.byID {
		if query != "" && !strings.Contains(strings.ToLower(artist.DisplayName), query) {
			continue
		}
		items = append(items, modelartist.Summary{
			ID:           artist.ID,
			DisplayName:  artist.DisplayName,
			AvatarURL:    artist.AvatarURL,
			Bio:          artist.Bio,
			ArtworkCount: artist.ArtworkCount,
		})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].DisplayName < items[j].DisplayName
	})

	total := len(items)
	items = pageSlice(items, page, pageSize)
	return items, paginationMeta(page, pageSize, total)
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
