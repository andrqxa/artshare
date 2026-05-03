package artist

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/andrqxa/artshare/internal/controller/apperror"
	modelartist "github.com/andrqxa/artshare/internal/model/artist"
	"github.com/andrqxa/artshare/internal/model/pagination"
)

type Controller struct {
	mu          sync.RWMutex
	nextID      int
	byID        map[string]modelartist.Detail
	byUserID    map[string]string
	currentUser string
}

func NewController() *Controller {
	return &Controller{
		nextID:      1,
		byID:        make(map[string]modelartist.Detail),
		byUserID:    make(map[string]string),
		currentUser: "00000000-0000-0000-0000-000000000001",
	}
}

func (c *Controller) CreateCurrentArtist(input modelartist.CreateInput) (modelartist.Detail, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	userID := input.UserID
	if userID == "" {
		userID = c.currentUser
	}
	if _, exists := c.byUserID[userID]; exists {
		return modelartist.Detail{}, apperror.ErrConflict
	}

	now := time.Now().UTC()
	artist := modelartist.Detail{
		ID:          nextUUID(c.nextID),
		UserID:      userID,
		DisplayName: input.DisplayName,
		Bio:         input.Bio,
		AvatarURL:   input.AvatarURL,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	c.nextID++
	c.byID[artist.ID] = artist
	c.byUserID[userID] = artist.ID

	return artist, nil
}

func (c *Controller) GetCurrentArtist(userID string) (modelartist.Detail, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if userID == "" {
		userID = c.currentUser
	}
	artistID, exists := c.byUserID[userID]
	if !exists {
		return modelartist.Detail{}, apperror.ErrNotFound
	}

	return c.byID[artistID], nil
}

func (c *Controller) UpdateCurrentArtist(userID string, input modelartist.UpdateInput) (modelartist.Detail, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if userID == "" {
		userID = c.currentUser
	}
	artistID, exists := c.byUserID[userID]
	if !exists {
		return modelartist.Detail{}, apperror.ErrNotFound
	}

	artist := c.byID[artistID]
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
	c.byID[artist.ID] = artist

	return artist, nil
}

func (c *Controller) ListArtists(filter modelartist.ListFilter) ([]modelartist.Summary, pagination.Meta) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	page, pageSize := normalizePagination(filter.Page, filter.PageSize)
	query := strings.TrimSpace(strings.ToLower(valueOrEmpty(filter.Query)))

	items := make([]modelartist.Summary, 0, len(c.byID))
	for _, artist := range c.byID {
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

func (c *Controller) GetArtistByID(id string) (modelartist.Detail, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	artist, exists := c.byID[id]
	if !exists {
		return modelartist.Detail{}, apperror.ErrNotFound
	}
	return artist, nil
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
