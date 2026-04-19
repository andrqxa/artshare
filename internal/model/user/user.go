package user

import "time"

type Role string

const (
	RoleViewer Role = "viewer"
	RoleArtist Role = "artist"
)

type ArtistProfile struct {
	ID          string
	UserID      string
	DisplayName string
	Bio         *string
	AvatarURL   *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CurrentUser struct {
	ID            string
	Email         string
	Role          Role
	CreatedAt     time.Time
	ArtistProfile *ArtistProfile
}

type UpdateCurrentUserInput struct {
	Email *string
}
