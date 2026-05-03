package artist

import "time"

type CreateArtistRequest struct {
	DisplayName string  `json:"displayName"`
	Bio         *string `json:"bio,omitempty"`
	AvatarURL   *string `json:"avatarUrl,omitempty"`
}

type UpdateArtistRequest struct {
	DisplayName *string `json:"displayName,omitempty"`
	Bio         *string `json:"bio,omitempty"`
	AvatarURL   *string `json:"avatarUrl,omitempty"`
}

type ArtistSummary struct {
	ID           string  `json:"id"`
	DisplayName  string  `json:"displayName"`
	AvatarURL    *string `json:"avatarUrl,omitempty"`
	Bio          *string `json:"bio,omitempty"`
	ArtworkCount int     `json:"artworkCount"`
}

type ArtistDetail struct {
	ID           string    `json:"id"`
	DisplayName  string    `json:"displayName"`
	AvatarURL    *string   `json:"avatarUrl,omitempty"`
	Bio          *string   `json:"bio,omitempty"`
	ArtworkCount int       `json:"artworkCount"`
	UserID       string    `json:"userId"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type ArtistListResponse struct {
	Data []ArtistSummary `json:"data"`
	Meta PaginationMeta  `json:"meta"`
}

type PaginationMeta struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}
