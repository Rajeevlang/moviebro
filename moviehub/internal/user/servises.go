package user

import (
	"context"
)

type Service interface {
	// RegisterUser should hash the password, call repo.CreateUser, and then emit the Kafka event
	RegisterUser(ctx context.Context, req *CreateUserRequest) (*User, error)

	// LoginUser should verify the password hash and return a JWT token string
	LoginUser(ctx context.Context, req *LoginRequest) (string, error)

	LoginGoogle(ctx context.Context, req *GoogleLoginRequest) (string, error)
}
