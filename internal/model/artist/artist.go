package artist

import "time"

type Artist struct {
	ID          string
	UserID      string
	DisplayName string
	Bio         *string
	AvatarURL   *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Summary struct {
	ID           string
	DisplayName  string
	AvatarURL    *string
	Bio          *string
	ArtworkCount int
}

type Detail struct {
	ID           string
	UserID       string
	DisplayName  string
	Bio          *string
	AvatarURL    *string
	ArtworkCount int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CreateInput struct {
	UserID      string
	DisplayName string
	Bio         *string
	AvatarURL   *string
}

type UpdateInput struct {
	DisplayName *string
	Bio         *string
	AvatarURL   *string
}
