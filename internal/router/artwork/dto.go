package artwork

import "time"

type ArtworkImage struct {
	URL     string  `json:"url"`
	AltText *string `json:"altText,omitempty"`
}

type ArtworkOwner struct {
	ID          string  `json:"id"`
	DisplayName string  `json:"displayName"`
	AvatarURL   *string `json:"avatarUrl,omitempty"`
}

type ArtworkSummary struct {
	ID                      string       `json:"id"`
	Title                   string       `json:"title"`
	Status                  string       `json:"status"`
	Owner                   ArtworkOwner `json:"owner"`
	PreviewImageURL         *string      `json:"previewImageUrl,omitempty"`
	LikeCount               int          `json:"likeCount"`
	LikedByMe               *bool        `json:"likedByMe,omitempty"`
	AvailabilityForExchange bool         `json:"availabilityForExchange"`
	CreatedAt               time.Time    `json:"createdAt"`
	PublishedAt             *time.Time   `json:"publishedAt,omitempty"`
}

type ArtworkDetail struct {
	ID                      string         `json:"id"`
	Title                   string         `json:"title"`
	Status                  string         `json:"status"`
	Owner                   ArtworkOwner   `json:"owner"`
	PreviewImageURL         *string        `json:"previewImageUrl,omitempty"`
	LikeCount               int            `json:"likeCount"`
	LikedByMe               *bool          `json:"likedByMe,omitempty"`
	AvailabilityForExchange bool           `json:"availabilityForExchange"`
	CreatedAt               time.Time      `json:"createdAt"`
	PublishedAt             *time.Time     `json:"publishedAt,omitempty"`
	Description             string         `json:"description"`
	Tags                    []string       `json:"tags,omitempty"`
	Images                  []ArtworkImage `json:"images"`
	UpdatedAt               time.Time      `json:"updatedAt"`
}

type CreateArtworkRequest struct {
	Title                   string         `json:"title"`
	Description             string         `json:"description"`
	Tags                    []string       `json:"tags,omitempty"`
	Images                  []ArtworkImage `json:"images"`
	AvailabilityForExchange bool           `json:"availabilityForExchange,omitempty"`
}

type UpdateArtworkRequest struct {
	Title                   *string        `json:"title,omitempty"`
	Description             *string        `json:"description,omitempty"`
	Tags                    []string       `json:"tags,omitempty"`
	Images                  []ArtworkImage `json:"images,omitempty"`
	AvailabilityForExchange *bool          `json:"availabilityForExchange,omitempty"`
	Status                  *string        `json:"status,omitempty"`
}

type ArtworkListResponse struct {
	Data []ArtworkSummary `json:"data"`
	Meta PaginationMeta   `json:"meta"`
}

type PaginationMeta struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}
