package dto

import "github.com/google/uuid"

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email,max=255"`
	Password  string `json:"password" validate:"required,min=8,max=20"`
	Username  string `json:"username" validate:"required,min=3,max=42"`
	FirstName string `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string `json:"last_name" validate:"required,min=2,max=100"`
}

type RegisterResponse struct {
	UserID       uuid.UUID `json:"user_id"`
	Role         string    `json:"role"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=42"`
}
type LoginResponse struct {
	UserID       uuid.UUID `json:"user_id"`
	Role         string    `json:"role"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
type RefreshResponse struct {
	UserID       uuid.UUID `json:"user_id"`
	Role         string    `json:"role"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
