package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type X struct {
	conn *pgxpool.Pool
}

func Newpoolstruct(conn *pgxpool.Pool) *X {
	return &X{
		conn: conn,
	}
}

func (x *X) CreateUser(ctx context.Context, u *User) error {

	query := `
		INSERT INTO users (username, email, password_hash, google_id)
		VALUES (@username, @email, @password_hash, @google_id)
		RETURNING id, username, email, password_hash, google_id, created_at, updated_at
	`

	args := pgx.NamedArgs{
		"username":      u.Username,
		"email":         u.Email,
		"password_hash": u.PasswordHash,
		"google_id":     u.GoogleID,
	}

	var createdUser User

	// Execute the insert and instantly scan the newly generated ID and Timestamps back into a User struct
	err := x.conn.QueryRow(ctx, query, args).Scan(
		&createdUser.ID,
		&createdUser.Username,
		&createdUser.Email,
		&createdUser.PasswordHash,
		&createdUser.GoogleID,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil

}

func (x *X) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, username, email, password_hash, google_id, created_at, updated_at
		FROM users 
		WHERE email = @email
	`

	args := pgx.NamedArgs{
		"email": email,
	}

	var u User
	err := x.conn.QueryRow(ctx, query, args).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.PasswordHash,
		&u.GoogleID,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		// Handle the specific case where the email doesn't exist
		if errors.Is(err, pgx.ErrNoRows) {
			return &User{}, fmt.Errorf("no user found with email %s", email)
		}
		return &User{}, fmt.Errorf("database error: %w", err)
	}

	return &u, nil

}
func (x *X) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `
		SELECT id, username, email, password_hash, google_id, created_at, updated_at
		FROM users 
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	var u User
	err := x.conn.QueryRow(ctx, query, args).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.PasswordHash,
		&u.GoogleID,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &User{}, fmt.Errorf("user not found")
		}
		return &User{}, fmt.Errorf("database error: %w", err)
	}

	return &u, nil

}
