package user

import "time"

type UpdateCurrentUserRequest struct {
	Email *string `json:"email,omitempty"`
}

type CurrentUserResponse struct {
	ID            string                 `json:"id"`
	Email         string                 `json:"email"`
	Role          string                 `json:"role"`
	CreatedAt     time.Time              `json:"createdAt"`
	ArtistProfile *ArtistProfileResponse `json:"artistProfile"`
}

type ArtistProfileResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	DisplayName string    `json:"displayName"`
	Bio         *string   `json:"bio"`
	AvatarURL   *string   `json:"avatarUrl"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
