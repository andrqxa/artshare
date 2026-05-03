package user

import "time"

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateCurrentUserRequest struct {
	Email *string `json:"email,omitempty"`
}

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	Artist    *Artist   `json:"artist,omitempty"`
}

type CurrentUserResponse = User

type AuthResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   int    `json:"expiresIn"`
	User        User   `json:"user"`
}

type Artist struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	DisplayName string    `json:"displayName"`
	Bio         *string   `json:"bio,omitempty"`
	AvatarURL   *string   `json:"avatarUrl,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
