package user

import (
	"time"
)

// User represents the database schema for a user in PostgreSQL
type User struct {
	ID           string    `json:"id" db:"id"` // Using string for UUID
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash *string   `json:"-" db:"password_hash"` // Pointer because it can be null (Google Auth)
	GoogleID     *string   `json:"google_id,omitempty" db:"google_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// CreateUserRequest represents the expected payload for registering a new user via Basic Auth
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// GoogleLoginRequest represents the payload from the client containing the Google OpenID token
type GoogleLoginRequest struct {
	IDToken string `json:"id_token" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse is the clean payload we send back to clients
type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
