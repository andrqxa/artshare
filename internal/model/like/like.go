package like

import "time"

type Like struct {
	UserID    string
	ArtworkID string
	CreatedAt time.Time
}

type CreateInput struct {
	UserID    string
	ArtworkID string
}

type DeleteInput struct {
	UserID    string
	ArtworkID string
}
