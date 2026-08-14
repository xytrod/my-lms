package client

import (
	"context"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
type User struct {
	ID       uuid.UUID `json:"id"`
	Role     string    `json:"role"`
	IsActive bool      `json:"is_active"`
}
type UserClient interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (*User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetuserByID(ctx context.Context, id uuid.UUID) (*User, error)
}
