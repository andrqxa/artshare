package user

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type CreateUserResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   int    `json:"expiresIn"`
	User        User   `json:"user"`
}

type User struct { // TODO: CHeck!!!
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
