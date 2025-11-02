package repository

import (
	"context"
	"database/sql"
	"time"

	"go-backend-valos-id/core/internal/repository"
	"go-backend-valos-id/core/user/model"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool    *pgxpool.Pool
	queries *repository.Queries
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool:    pool,
		queries: repository.New(pool),
	}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(user *model.User) error {
	ctx := context.Background()
	now := time.Now()

	// Convert time.Time to pgtype.Timestamptz for pgx
	timestamptz := pgtype.Timestamptz{
		Time:  now,
		Valid: true,
	}

	params := repository.CreateUserParams{
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: timestamptz,
		UpdatedAt: timestamptz,
	}

	result, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		// Check for unique constraint violation
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "users_email_unique" || pgErr.ConstraintName == "users_email_key" {
				return &pgconn.PgError{
					Code:    "23505",
					Message: "user with this email already exists",
				}
			}
			if pgErr.ConstraintName == "users_username_unique" || pgErr.ConstraintName == "users_username_key" {
				return &pgconn.PgError{
					Code:    "23505",
					Message: "user with this username already exists",
				}
			}
		}
		return err
	}

	user.ID = result.ID
	user.CreatedAt = now
	user.UpdatedAt = now

	return nil
}

// GetUserByID retrieves a user by their ID
func (r *UserRepository) GetUserByID(id int32) (*model.User, error) {
	ctx := context.Background()

	result, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return r.sqlcUserToModelUser(&result), nil
}

// GetUserByEmail retrieves a user by their email
func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	ctx := context.Background()

	result, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return r.sqlcUserToModelUser(&result), nil
}

// GetAllUsers retrieves all users from the database
func (r *UserRepository) GetAllUsers() ([]model.User, error) {
	ctx := context.Background()

	results, err := r.queries.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]model.User, len(results))
	for i, result := range results {
		createdAt := time.Time{}
		if result.CreatedAt.Valid {
			createdAt = result.CreatedAt.Time
		}

		updatedAt := time.Time{}
		if result.UpdatedAt.Valid {
			updatedAt = result.UpdatedAt.Time
		}

		users[i] = model.User{
			ID:        result.ID,
			Username:  result.Username,
			Email:     result.Email,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
	}

	return users, nil
}

// UpdateUser updates an existing user
func (r *UserRepository) UpdateUser(user *model.User) error {
	ctx := context.Background()
	now := time.Now()

	// Convert time.Time to pgtype.Timestamptz for pgx
	timestamptz := pgtype.Timestamptz{
		Time:  now,
		Valid: true,
	}

	params := repository.UpdateUserParams{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		UpdatedAt: timestamptz,
	}

	err := r.queries.UpdateUser(ctx, params)
	if err != nil {
		// Check for unique constraint violation
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "users_email_unique" || pgErr.ConstraintName == "users_email_key" {
				return &pgconn.PgError{
					Code:    "23505",
					Message: "user with this email already exists",
				}
			}
			if pgErr.ConstraintName == "users_username_unique" || pgErr.ConstraintName == "users_username_key" {
				return &pgconn.PgError{
					Code:    "23505",
					Message: "user with this username already exists",
				}
			}
		}
		return err
	}

	user.UpdatedAt = now
	return nil
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(userID int32, hashedPassword string) error {
	ctx := context.Background()
	now := time.Now()

	// Convert time.Time to pgtype.Timestamptz for pgx
	timestamptz := pgtype.Timestamptz{
		Time:  now,
		Valid: true,
	}

	params := repository.UpdatePasswordParams{
		ID:        userID,
		Password:  hashedPassword,
		UpdatedAt: timestamptz,
	}

	err := r.queries.UpdatePassword(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser deletes a user by their ID
func (r *UserRepository) DeleteUser(id int32) error {
	ctx := context.Background()

	err := r.queries.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// UserExists checks if a user exists by email
func (r *UserRepository) UserExists(email string) (bool, error) {
	ctx := context.Background()

	exists, err := r.queries.UserExists(ctx, email)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetUsersWithPagination retrieves users with pagination
func (r *UserRepository) GetUsersWithPagination(limit, offset int32) ([]model.User, error) {
	ctx := context.Background()

	params := repository.GetUsersWithPaginationParams{
		Limit:  limit,
		Offset: offset,
	}

	results, err := r.queries.GetUsersWithPagination(ctx, params)
	if err != nil {
		return nil, err
	}

	users := make([]model.User, len(results))
	for i, result := range results {
		createdAt := time.Time{}
		if result.CreatedAt.Valid {
			createdAt = result.CreatedAt.Time
		}

		updatedAt := time.Time{}
		if result.UpdatedAt.Valid {
			updatedAt = result.UpdatedAt.Time
		}

		users[i] = model.User{
			ID:        result.ID,
			Username:  result.Username,
			Email:     result.Email,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
	}

	return users, nil
}

// CountUsers returns the total number of users
func (r *UserRepository) CountUsers() (int, error) {
	ctx := context.Background()

	count, err := r.queries.CountUsers(ctx)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// Helper method to convert sqlc User to model User
func (r *UserRepository) sqlcUserToModelUser(sqlcUser *repository.User) *model.User {
	createdAt := time.Time{}
	if sqlcUser.CreatedAt.Valid {
		createdAt = sqlcUser.CreatedAt.Time
	}

	updatedAt := time.Time{}
	if sqlcUser.UpdatedAt.Valid {
		updatedAt = sqlcUser.UpdatedAt.Time
	}

	return &model.User{
		ID:        sqlcUser.ID,
		Username:  sqlcUser.Username,
		Email:     sqlcUser.Email,
		Password:  sqlcUser.Password,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
