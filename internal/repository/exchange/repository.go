package exchange

import (
	"fmt"
	"sort"
	"sync"
	"time"

	modelartwork "github.com/andrqxa/artshare/internal/model/artwork"
	modelexchange "github.com/andrqxa/artshare/internal/model/exchange"
	"github.com/andrqxa/artshare/internal/model/pagination"
	"github.com/andrqxa/artshare/internal/repository/repoerror"
)

type Repository struct {
	mu     sync.RWMutex
	nextID int
	byID   map[string]modelexchange.Detail
}

func NewRepository() *Repository {
	return &Repository{
		nextID: 1,
		byID:   make(map[string]modelexchange.Detail),
	}
}

func (r *Repository) List(filter modelexchange.ListFilter) ([]modelexchange.Summary, pagination.Meta) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	page, pageSize := normalizePagination(filter.Page, filter.PageSize)
	scope := filter.Scope
	if scope == "" {
		scope = modelexchange.ScopeAll
	}

	items := make([]modelexchange.Summary, 0, len(r.byID))
	for _, request := range r.byID {
		if filter.UserID != "" {
			if scope == modelexchange.ScopeInbox && request.OwnerID != filter.UserID {
				continue
			}
			if scope == modelexchange.ScopeOutbox && request.RequesterID != filter.UserID {
				continue
			}
			if scope == modelexchange.ScopeAll && request.OwnerID != filter.UserID && request.RequesterID != filter.UserID {
				continue
			}
		}
		if filter.Status != nil && request.Status != *filter.Status {
			continue
		}
		items = append(items, summaryFromDetail(request))
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})

	total := len(items)
	items = pageSlice(items, page, pageSize)
	return items, paginationMeta(page, pageSize, total)
}

func (r *Repository) Create(input modelexchange.CreateInput) (modelexchange.Detail, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	request := modelexchange.Detail{
		ID: nextUUID(r.nextID),
		Artwork: modelartwork.OwneredRef{
			ID:      input.ArtworkID,
			Title:   "Artwork",
			OwnerID: "00000000-0000-0000-0000-000000000001",
		},
		RequesterID: input.RequesterID,
		OwnerID:     "00000000-0000-0000-0000-000000000001",
		Status:      modelexchange.StatusPending,
		Message:     input.Message,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	r.nextID++
	r.byID[request.ID] = request

	return request, nil
}

func (r *Repository) FindByID(id string) (modelexchange.Detail, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	request, exists := r.byID[id]
	if !exists {
		return modelexchange.Detail{}, repoerror.ErrNotFound
	}
	return request, nil
}

func (r *Repository) UpdateStatus(input modelexchange.UpdateStatusInput) (modelexchange.Detail, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	request, exists := r.byID[input.ID]
	if !exists {
		return modelexchange.Detail{}, repoerror.ErrNotFound
	}
	request.Status = input.Status
	if input.Message != nil {
		request.Message = input.Message
	}
	request.UpdatedAt = time.Now().UTC()
	r.byID[request.ID] = request

	return request, nil
}

func summaryFromDetail(request modelexchange.Detail) modelexchange.Summary {
	return modelexchange.Summary{
		ID:          request.ID,
		Artwork:     request.Artwork,
		RequesterID: request.RequesterID,
		OwnerID:     request.OwnerID,
		Status:      request.Status,
		CreatedAt:   request.CreatedAt,
		UpdatedAt:   request.UpdatedAt,
	}
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
