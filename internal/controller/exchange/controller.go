package exchange

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/andrqxa/artshare/internal/controller/apperror"
	modelartwork "github.com/andrqxa/artshare/internal/model/artwork"
	modelexchange "github.com/andrqxa/artshare/internal/model/exchange"
	"github.com/andrqxa/artshare/internal/model/pagination"
)

type Controller struct {
	mu     sync.RWMutex
	nextID int
	byID   map[string]modelexchange.Detail
}

func NewController() *Controller {
	return &Controller{
		nextID: 1,
		byID:   make(map[string]modelexchange.Detail),
	}
}

func (c *Controller) ListExchangeRequests(filter modelexchange.ListFilter) ([]modelexchange.Summary, pagination.Meta) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	page, pageSize := normalizePagination(filter.Page, filter.PageSize)
	scope := filter.Scope
	if scope == "" {
		scope = modelexchange.ScopeAll
	}

	items := make([]modelexchange.Summary, 0, len(c.byID))
	for _, request := range c.byID {
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

func (c *Controller) CreateExchangeRequest(input modelexchange.CreateInput) (modelexchange.Detail, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().UTC()
	request := modelexchange.Detail{
		ID: nextUUID(c.nextID),
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
	c.nextID++
	c.byID[request.ID] = request

	return request, nil
}

func (c *Controller) GetExchangeRequestByID(id string) (modelexchange.Detail, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	request, exists := c.byID[id]
	if !exists {
		return modelexchange.Detail{}, apperror.ErrNotFound
	}
	return request, nil
}

func (c *Controller) UpdateExchangeRequestStatus(input modelexchange.UpdateStatusInput) (modelexchange.Detail, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	request, exists := c.byID[input.ID]
	if !exists {
		return modelexchange.Detail{}, apperror.ErrNotFound
	}
	if request.Status != modelexchange.StatusPending {
		return modelexchange.Detail{}, apperror.ErrConflict
	}

	request.Status = input.Status
	if input.Message != nil {
		request.Message = input.Message
	}
	request.UpdatedAt = time.Now().UTC()
	c.byID[request.ID] = request

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
	return pagination.Meta{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
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
