package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"go-backend-valos-id/core/user/model"
	"go-backend-valos-id/core/user/repository"
	"go-backend-valos-id/core/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserHandler struct {
	userRepo *repository.UserRepository
}

func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// CreateUser handles the creation of a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req model.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Check if user already exists
	exists, err := h.userRepo.UserExists(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check if user exists",
		})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"error": "User with this email already exists",
		})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create user
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := h.userRepo.CreateUser(user); err != nil {
		// Check for unique constraint violation
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "users_email_unique" || pgErr.ConstraintName == "users_email_key" {
				c.JSON(http.StatusConflict, gin.H{
					"error": "User with this email already exists",
				})
				return
			}
			if pgErr.ConstraintName == "users_username_unique" || pgErr.ConstraintName == "users_username_key" {
				c.JSON(http.StatusConflict, gin.H{
					"error": "User with this username already exists",
				})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    h.toUserResponse(user),
	})
}

// GetAllUsers retrieves all users
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userRepo.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve users",
		})
		return
	}

	userResponses := make([]model.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = h.toUserResponse(&user)
	}

	c.JSON(http.StatusOK, gin.H{
		"users": userResponses,
		"count": len(userResponses),
	})
}

// GetUserByID retrieves a user by ID
func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID, err := h.parseUserID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.userRepo.GetUserByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": h.toUserResponse(user),
	})
}

// UpdateUser handles updating an existing user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, err := h.parseUserID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check if user exists first
	_, err = h.userRepo.GetUserByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user",
		})
		return
	}

	var req struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Email    string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	user := &model.User{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
	}

	if err := h.userRepo.UpdateUser(user); err != nil {
		// Check for unique constraint violation
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "users_email_unique" || pgErr.ConstraintName == "users_email_key" {
				c.JSON(http.StatusConflict, gin.H{
					"error": "User with this email already exists",
				})
				return
			}
			if pgErr.ConstraintName == "users_username_unique" || pgErr.ConstraintName == "users_username_key" {
				c.JSON(http.StatusConflict, gin.H{
					"error": "User with this username already exists",
				})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    h.toUserResponse(user),
	})
}

// DeleteUser handles deleting a user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := h.parseUserID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.userRepo.DeleteUser(userID); err != nil {
		// sqlc DeleteUser doesn't return an error for no rows affected
		// We need to check if the user exists first
		if _, checkErr := h.userRepo.GetUserByID(userID); checkErr == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// GetUsersWithPagination retrieves users with pagination
func (h *UserHandler) GetUsersWithPagination(c *gin.Context) {
	limit, err := h.parseIntQuery(c.Query("limit"), 10, 1, 100)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid limit parameter",
		})
		return
	}

	offset, err := h.parseIntQuery(c.Query("offset"), 0, 0, 1000000)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid offset parameter",
		})
		return
	}

	users, err := h.userRepo.GetUsersWithPagination(int32(limit), int32(offset))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve users",
		})
		return
	}

	total, err := h.userRepo.CountUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to count users",
		})
		return
	}

	userResponses := make([]model.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = h.toUserResponse(&user)
	}

	c.JSON(http.StatusOK, gin.H{
		"users": userResponses,
		"pagination": gin.H{
			"limit":  limit,
			"offset": offset,
			"total":  total,
		},
	})
}

// Helper methods

func (h *UserHandler) parseUserID(idStr string) (int32, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID")
	}
	if id <= 0 {
		return 0, fmt.Errorf("user ID must be positive")
	}
	return int32(id), nil
}

func (h *UserHandler) parseIntQuery(value string, defaultValue, min, max int) (int, error) {
	if value == "" {
		return defaultValue, nil
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid integer value")
	}

	if result < min {
		result = min
	}
	if result > max {
		result = max
	}

	return result, nil
}

func (h *UserHandler) toUserResponse(user *model.User) model.UserResponse {
	return model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
