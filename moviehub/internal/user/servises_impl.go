package user

import (
	"context"
	"fmt"
	"moviesapi/internal/shared"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// 1. The Concrete Struct
// This is the "video game console". It holds our dependencies.
// We keep it unexported (lowercase 'u') so outside packages can't mess with it.
type userService struct {
	repo   Repository // This is the "cartridge slot". Dependency Injection happens here!
	config shared.Config
}

// 2. The Constructor
// We call this from main.go. We hand it a Repository, and it gives us back the Service interface.
func NewService(r Repository, cfg shared.Config) Service {
	return &userService{
		repo:   r,
		config: cfg,
	}
}

// ---------------------------------------------------------
// 3. The Methods
// Because these methods are attached to `(s *userService)`,
// they can all access the database by calling `s.repo`!
// ---------------------------------------------------------

func (s *userService) RegisterUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
	// TODO: Hash password using bcrypt

	// Example of using the injected dependency:
	// err := s.repo.CreateUser(ctx, newUser)
	psswd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	hashStr := string(psswd)
	user := &User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: &hashStr,
		GoogleID:     nil,
	}

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// Publish event to Redis Stream
	if err := shared.PublishUserCreatedEvent(ctx, user.ID, user.Username); err != nil {
		shared.Log.Error("failed to publish user_created event", zap.Error(err))
	}

	return user, nil
}

func (s *userService) LoginUser(ctx context.Context, req *LoginRequest) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		// We return a generic error so we don't leak whether the email exists
		return "", fmt.Errorf("invalid credentials")
	}

	if user.PasswordHash == nil {
		return "", fmt.Errorf("this account uses Google Login, please log in with Google")
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.Password))
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := Newjwttoken(user.ID, user.Username, []byte(s.config.SECERETKEY))
	if err != nil {
		return "", fmt.Errorf("failed to generate token")
	}

	return token, nil
}

func (s *userService) LoginGoogle(ctx context.Context, req *GoogleLoginRequest) (string, error) {
	idToken, err := Oidcverifier.Verify(ctx, req.IDToken)
	if err != nil {
		return "", fmt.Errorf("failed to verify id token: %w", err)
	}

	var claims struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		Sub           string `json:"sub"` // Google's unique subject ID
	}

	if err := idToken.Claims(&claims); err != nil {
		return "", fmt.Errorf("failed to extract claims: %w", err)
	}

	if !claims.EmailVerified {
		return "", fmt.Errorf("google email is not verified")
	}

	// Check if user already exists
	user, err := s.repo.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		// User does not exist, so let's register them!
		user = &User{
			Username:     claims.Name,
			Email:        claims.Email,
			PasswordHash: nil, // No password for Google logins
			GoogleID:     &claims.Sub,
		}

		err = s.repo.CreateUser(ctx, user)
		if err != nil {
			return "", fmt.Errorf("failed to create user from google login: %w", err)
		}

		// Publish event to Redis Stream
		if err := shared.PublishUserCreatedEvent(ctx, user.ID, user.Username); err != nil {
			shared.Log.Error("failed to publish user_created event", zap.Error(err))
		}
	}

	// Issue our application's JWT Token
	token, err := Newjwttoken(user.ID, user.Username, []byte(s.config.SECERETKEY))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}
