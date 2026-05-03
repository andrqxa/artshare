package user

import (
	"time"

	modelartist "github.com/andrqxa/artshare/internal/model/artist"
)

type Role string

const (
	RoleViewer Role = "viewer"
	RoleArtist Role = "artist"
)

type ArtistProfile = modelartist.Artist

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

type RegisterInput struct {
	Email    string
	Password string
	Role     Role
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthToken struct {
	AccessToken string
	TokenType   string
	ExpiresIn   int
}
