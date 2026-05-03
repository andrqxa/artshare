package exchange

import "time"

type CreateExchangeRequest struct {
	ArtworkID string  `json:"artworkId"`
	Message   *string `json:"message,omitempty"`
}

type UpdateExchangeRequestStatusRequest struct {
	Status  string  `json:"status"`
	Message *string `json:"message,omitempty"`
}

type ArtworkOwneredRef struct {
	ID              string  `json:"id"`
	Title           string  `json:"title"`
	OwnerID         string  `json:"ownerId"`
	PreviewImageURL *string `json:"previewImageUrl,omitempty"`
}

type ExchangeRequestSummary struct {
	ID          string            `json:"id"`
	Artwork     ArtworkOwneredRef `json:"artwork"`
	RequesterID string            `json:"requesterId"`
	OwnerID     string            `json:"ownerId"`
	Status      string            `json:"status"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

type ExchangeRequestDetail struct {
	ID          string            `json:"id"`
	Artwork     ArtworkOwneredRef `json:"artwork"`
	RequesterID string            `json:"requesterId"`
	OwnerID     string            `json:"ownerId"`
	Status      string            `json:"status"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	Message     *string           `json:"message,omitempty"`
}

type ExchangeRequestListResponse struct {
	Data []ExchangeRequestSummary `json:"data"`
	Meta PaginationMeta           `json:"meta"`
}

type PaginationMeta struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}
