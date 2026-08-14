package service

import (
	"context"
	"errors"
	"fmt"
	"main/auth-service/internal/apperror"
	"main/auth-service/internal/client"
	"main/auth-service/internal/dto"
	"main/auth-service/internal/hashing"
	"main/auth-service/internal/model"
	"main/auth-service/internal/repo"
	"main/auth-service/internal/token"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	Refresh(ctx context.Context, req dto.RefreshRequest) (*dto.RefreshResponse, error)
	Logout(ctx context.Context, req dto.LogoutRequest) error
}
type authService struct {
	repo         repo.AuthRepository
	session      repo.SessionRepo
	userClient   client.UserClient
	passwordHash hashing.Hasher
	tokenManager token.Manager
}

func NewAuthService(repo repo.AuthRepository, session repo.SessionRepo, userClient client.UserClient, passwordHash hashing.Hasher, tokenManager token.Manager) AuthService {
	return &authService{
		repo:         repo,
		session:      session,
		userClient:   userClient,
		passwordHash: passwordHash,
		tokenManager: tokenManager,
	}
}
func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	username := strings.TrimSpace(req.Username)
	firstName := strings.TrimSpace(req.FirstName)
	lastName := strings.TrimSpace(req.LastName)

	exists, err := s.repo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("check credential_e", err)
	}
	if exists {
		return nil, apperror.Conflict("credential_already_exists", "user with this email already exists", ErrCredentialAlreadyExists)
	}
	hashed, err := s.passwordHash.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password", err)
	}
	created, err := s.userClient.CreateUser(ctx, client.CreateUserRequest{
		Email:     email,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
	})
	if err != nil {
		if errors.Is(err, client.ErrUserAlreadyExists) {
			return nil, ErrCredentialAlreadyExists
		}
		if errors.Is(err, client.ErrUserService) {
			return nil, apperror.ServiceUnavailable("user_service_unavailable", "user service is unavailable in the moment", err)
		}
		return nil, fmt.Errorf("create user prof", err)
	}
	credential := &model.Credential{
		ID:           uuid.New(),
		UserID:       created.ID,
		Email:        email,
		PasswordHash: hashed,
	}
	if err := s.repo.Create(ctx, credential); err != nil {
		clean, cancelled := context.WithTimeout(context.WithoutCancel(ctx), 3*time.Second)
		defer cancelled()
		created_err := s.userClient.DeleteUser(clean, created.ID)
		if created_err != nil {
			return nil, fmt.Errorf("delete user prof", ErrRegistrationFailed, err, created_err)
		}
		return nil, apperror.Internal("registration_failed", "registration could not be processed", err)
	}
	tokenPair, err := s.tokenManager.GeneratePair(created.ID, created.Role)
	if err != nil {
		return nil, fmt.Errorf("generate pair", err)
	}
	if err := s.createSession(ctx, created.ID, tokenPair); err != nil {
		return nil, apperror.Internal("session_creation_failed", "session creation failed", err)
	}
	return &dto.RegisterResponse{
		UserID:       created.ID,
		Role:         created.Role,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}
func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	credential, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.Unauthorized("invalid_credential", "invalid email or password", ErrInvalidCredentials)
		}
		return nil, apperror.Internal("auth_storage_error", "auth service is unavailable in the moment", err)
	}
	if err := s.passwordHash.CompareTo(credential.PasswordHash, req.Password); err != nil {
		if errors.Is(err, hashing.ErrPasswordMismatch) {
			return nil, apperror.Unauthorized("invalid_credential", "invalid email or password", ErrInvalidCredentials)
		}
		return nil, apperror.Internal("auth_storage course_errors", "auth service is unavailable in the moment", err)
	}
	user, err := s.userClient.GetuserByID(ctx, credential.UserID)
	if err != nil {
		if errors.Is(err, client.ErrUserNotFound) {
			return nil, apperror.Internal("profile not found", "authentification could not be processed", err)
		}
		return nil, apperror.Internal("login failed", "authentification could not be processed", err)
	}
	if !user.IsActive {
		return nil, apperror.Unauthorized("user_inactive", "user is inactive", ErrUserInActive)
	}
	tokenPair, err := s.tokenManager.GeneratePair(user.ID, user.Role)
	if err != nil {
		return nil, apperror.Internal("generate pair failed", "authentification could not be processed", err)
	}
	if err := s.createSession(ctx, credential.UserID, tokenPair); err != nil {
		return nil, apperror.Internal("session_creation_failed", "session creation failed", err)
	}
	return &dto.LoginResponse{
		UserID:       user.ID,
		Role:         user.Role,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}
func (s *authService) Refresh(ctx context.Context, req dto.RefreshRequest) (*dto.RefreshResponse, error) {
	raw := strings.TrimSpace(req.RefreshToken)
	claims, err := s.tokenManager.ParseRefreshToken(raw)
	if err != nil {
		return nil, apperror.Unauthorized("invalid_refresh_token", "invalid refresh raw", ErrInvalidRefreshToken)
	}
	session, err := s.session.GetByToken(ctx, claims.ID)
	if session.UserID != claims.UserID {
		return nil, apperror.Unauthorized("invalid_refresh_token", "invalid refresh raw", ErrInvalidRefreshToken)
	}
	if session.RevokedAt != nil {
		return nil, ErrInvalidRefreshToken
	}
	if time.Now().After(session.ExpiresAt) {
		return nil, ErrInvalidRefreshToken
	}
	if session.HashedToken != token.HashForSession(raw) {
		return nil, ErrInvalidRefreshToken
	}
	credential, err := s.repo.GetByUserId(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.Unauthorized("invalid_refresh_token", "invalid refresh raw", ErrInvalidRefreshToken)
		}
		return nil, apperror.Internal("refresh_failed", "refresh could not be processed", err)
	}
	user, err := s.userClient.GetuserByID(ctx, credential.UserID)
	if err != nil {
		if errors.Is(err, client.ErrUserNotFound) {
			return nil, apperror.Unauthorized("invalid_refresh_token", "invalid refresh raw", ErrInvalidRefreshToken)
		}
		if errors.Is(err, client.ErrUserService) {
			return nil, apperror.Unauthorized("user_service_unavailable", "user service is unavailable in the moment", err)
		}
		return nil, apperror.Internal("refresh_failed", "refresh could not be processed", err)
	}
	if !user.IsActive {
		return nil, apperror.Unauthorized("user_inactive", "user is active", ErrUserInActive)
	}
	if err := s.session.Revoke(ctx, claims.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.Unauthorized("invalid_refresh_token", "invalid refresh raw", ErrInvalidRefreshToken)
		}
		return nil, apperror.Internal("refresh_failed", "refresh could not be processed", err)
	}
	tokenPair, err := s.tokenManager.GeneratePair(user.ID, user.Role)
	if err != nil {
		return nil, apperror.Internal("generate pair failed", "authentification could not be processed", err)
	}
	if err := s.createSession(ctx, user.ID, tokenPair); err != nil {
		return nil, apperror.Internal("session_creation_failed", "session creation failed", err)
	}
	return &dto.RefreshResponse{
		UserID:       user.ID,
		Role:         user.Role,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}
func (s *authService) Logout(ctx context.Context, req dto.LogoutRequest) error {
	raw := strings.TrimSpace(req.RefreshToken)
	claims, err := s.tokenManager.ParseRefreshToken(raw)
	if err != nil {
		return apperror.Unauthorized("invalid_refresh_token", "invalid refresh token", err)
	}
	session, err := s.session.GetByToken(ctx, claims.ID)
	if err != nil {
		return apperror.Unauthorized("invalid_refresh_token", "invalid refresh token", err)
	}
	if session.HashedToken != token.HashForSession(raw) {
		return apperror.Unauthorized("invalid_refresh_token", "invalid refresh token", ErrInvalidRefreshToken)
	}
	if session.RevokedAt != nil {
		return nil
	}
	if err := s.session.Revoke(ctx, claims.ID); err != nil {
		return apperror.Internal("logout_failed", "logout was failed", err)
	}
	return nil
}
func (s *authService) createSession(ctx context.Context, userID uuid.UUID, pair *token.Token) error {
	session := &model.RefreshSession{
		ID:          uuid.New(),
		UserID:      userID,
		TokenID:     pair.RefreshTokenID,
		HashedToken: token.HashForSession(pair.RefreshToken),
		ExpiresAt:   pair.RefreshExpiresAt,
	}
	if err := s.session.Create(ctx, session); err != nil {
		return err
	}
	return nil
}
