package dto

import (
	"main/user-service/internal/model"
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=7,max=25"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}
type UpdateUserRequest struct {
	ID        uuid.UUID `json:"-"`
	Email     *string   `json:"email" validate:"omitempty,email"`
	Username  *string   `json:"username" validate:"omitempty,min=7,max=25"`
	FirstName *string   `json:"first_name" validate:"omitempty,min=2,max=100"`
	LastName  *string   `json:"last_name" validate:"omitempty,min=2,max=100"`
	IsActive  *bool     `json:"is_active"`
}

type UserResponse struct {
	ID        uuid.UUID      `json:"id"`
	Email     string         `json:"email"`
	Username  string         `json:"username"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Role      model.UserRole `json:"role"`
	IsActive  bool           `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
