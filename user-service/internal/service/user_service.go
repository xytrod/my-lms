package service

import (
	"context"
	"errors"
	"main/user-service/internal/dto"
	"main/user-service/internal/model"
	"main/user-service/internal/repository"
	"main/user-service/internal/user_errors"
	"net/mail"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (*model.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	UpdateUser(ctx context.Context, req dto.UpdateUserRequest) (*model.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]model.User, error)
}
type userService struct {
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*model.User, error) {
	email := normalizeEmail(req.Email)
	username := strings.TrimSpace(req.Username)
	firstName := strings.TrimSpace(req.FirstName)
	lastName := strings.TrimSpace(req.LastName)

	if err := validateUser(email, username, firstName, lastName); err != nil {
		return nil, err
	}

	if err := s.IsEmailUnique(ctx, email, uuid.Nil); err != nil {
		return nil, err
	}
	if err := s.IsUsernameUnique(ctx, username, uuid.Nil); err != nil {
		return nil, err
	}

	user := model.User{
		ID:        uuid.New(),
		Email:     email,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Role:      model.RoleStudent,
		IsActive:  true,
	}
	if err := s.repo.Create(ctx, &user); err != nil {
		return nil, err
	}
	return &user, nil
}
func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	if id == uuid.Nil {
		return nil, user_errors.ErrInvalidUserID
	}
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
func (s *userService) UpdateUser(ctx context.Context, req dto.UpdateUserRequest) (*model.User, error) {
	if req.ID == uuid.Nil {
		return nil, user_errors.ErrInvalidUserID
	}
	user, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, err
	}
	if req.Email != nil {
		email := normalizeEmail(*req.Email)

		if email == "" {
			return nil, user_errors.ErrEmailRequired
		}
		if !validatedEmail(email) {
			return nil, user_errors.ErrInvalidEmail
		}
		if email != user.Email {
			if err := s.IsEmailUnique(ctx, email, user.ID); err != nil {
				return nil, err
			}
			user.Email = email
		}
	}
	if req.Username != nil {
		username := strings.TrimSpace(*req.Username)
		if username == "" {
			return nil, user_errors.ErrUsernameRequired
		}
		if username != user.Username {
			if err := s.IsUsernameUnique(ctx, username, user.ID); err != nil {
				return nil, err
			}
			user.Username = username
		}
	}
	if req.FirstName != nil {
		firstName := strings.TrimSpace(*req.FirstName)
		if firstName == "" {
			return nil, user_errors.ErrFirstNameRequired
		}
		user.FirstName = firstName
	}
	if req.LastName != nil {
		lastName := strings.TrimSpace(*req.LastName)
		if lastName == "" {
			return nil, user_errors.ErrLastNameRequired
		}
		user.LastName = lastName
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	if err := s.repo.Update(ctx, user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return user_errors.ErrInvalidUserID
	}
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user_errors.ErrUserNotFound
		}
		return err
	}
	return nil
}
func (s *userService) List(ctx context.Context, limit int, offset int) ([]model.User, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	users, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return users, err
}

func normalizeEmail(email string) string {
	email = strings.ToLower(strings.TrimSpace(email))
	return email
}
func validatedEmail(email string) bool {
	receivedEmail, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	return receivedEmail.Address == email
}
func validateUser(email string, username string, firstName, lastName string) error {
	if email == "" {
		return user_errors.ErrEmailRequired
	}
	if !validatedEmail(email) {
		return user_errors.ErrInvalidEmail
	}
	if username == "" {
		return user_errors.ErrUsernameRequired
	}
	if firstName == "" {
		return user_errors.ErrFirstNameRequired
	}
	if lastName == "" {
		return user_errors.ErrLastNameRequired
	}
	return nil
}
func (s *userService) IsEmailUnique(ctx context.Context, email string, user_id uuid.UUID) error {
	user, err := s.repo.GetByEmail(ctx, email)
	if err == nil {
		if user_id != uuid.Nil && user.ID == user_id {
			return nil
		}
		return user_errors.ErrEmailAlreadyExists
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}
func (s *userService) IsUsernameUnique(ctx context.Context, username string, user_id uuid.UUID) error {
	user, err := s.repo.GetByUsername(ctx, username)
	if err == nil {
		if user.ID != uuid.Nil && user.ID == user_id {
			return nil
		}
		return user_errors.ErrUsernameAlreadyExists
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}
