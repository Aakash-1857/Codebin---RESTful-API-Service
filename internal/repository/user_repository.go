// internal/repository/user_repository.go
package repository

import (
	"context"
	"errors"

	"github.com/aakash-1857/codebin/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{DB: db}
}

// MODIFIED: Insert now creates and returns a new user model.
func (r *UserRepository) Insert(name, email string, passwordHash []byte) (*models.User, error) {
	stmt := `INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)
    RETURNING id, created_at`

	// This is the object that will be returned.
	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	}

	// Scan the auto-generated id and created_at fields back into the user struct.
	err := r.DB.QueryRow(context.Background(), stmt, name, email, passwordHash).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Return the pointer to the fully populated user struct.
	return user, nil
}

// GetByEmail remains the same, but please ensure it has the error handling fix.
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	stmt := `SELECT id, name, email, password_hash, created_at FROM users WHERE email = $1`
	user := &models.User{}
	row := r.DB.QueryRow(context.Background(), stmt, email)

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return user, nil
}