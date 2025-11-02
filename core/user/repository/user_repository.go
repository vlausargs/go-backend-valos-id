package repository

import (
	"database/sql"
	"fmt"
	"time"

	"go-backend-valos-id/core/user/model"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(user *model.User) error {
	query := `
		INSERT INTO users (username, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := r.db.QueryRow(query, user.Username, user.Email, user.Password, now, now).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByID retrieves a user by their ID
func (r *UserRepository) GetUserByID(id int) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = $1`

	err := r.db.Get(user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by their email
func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE email = $1`

	err := r.db.Get(user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetAllUsers retrieves all users from the database
func (r *UserRepository) GetAllUsers() ([]model.User, error) {
	var users []model.User
	query := `SELECT id, username, email, created_at, updated_at FROM users ORDER BY created_at DESC`

	err := r.db.Select(&users, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	return users, nil
}

// UpdateUser updates an existing user
func (r *UserRepository) UpdateUser(user *model.User) error {
	query := `
		UPDATE users
		SET username = $2, email = $3, updated_at = $4
		WHERE id = $1`

	user.UpdatedAt = time.Now()

	result, err := r.db.Exec(query, user.ID, user.Username, user.Email, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(userID int, hashedPassword string) error {
	query := `UPDATE users SET password = $2, updated_at = $3 WHERE id = $1`

	result, err := r.db.Exec(query, userID, hashedPassword, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// DeleteUser deletes a user by their ID
func (r *UserRepository) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// UserExists checks if a user exists by email
func (r *UserRepository) UserExists(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	err := r.db.Get(&exists, query, email)
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}

	return exists, nil
}

// GetUsersWithPagination retrieves users with pagination
func (r *UserRepository) GetUsersWithPagination(limit, offset int) ([]model.User, error) {
	var users []model.User
	query := `
		SELECT id, username, email, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	err := r.db.Select(&users, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get users with pagination: %w", err)
	}

	return users, nil
}

// CountUsers returns the total number of users
func (r *UserRepository) CountUsers() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users`

	err := r.db.Get(&count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}
