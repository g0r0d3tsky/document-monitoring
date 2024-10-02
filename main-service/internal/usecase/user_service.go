package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"main-service/internal/domain"
	"time"

	"github.com/google/uuid"
)

const (
	//salt       = "some-salt"
	signingKey = "some-signing-key"
	tokenTTL   = 12 * time.Hour
)

type UserInfo struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
}

type tokenClaims struct {
	jwt.RegisteredClaims
	UserClaims UserInfo `json:"userClaims"`
}

//go:generate mockgen -source=user.go -destination=mocks/userMock.go

type UserRepo interface {
	GetUser(ctx context.Context, login string, password string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, username string) error
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	UpdateUser(ctx context.Context, username string, updatedUser *domain.User) error
}

type UserService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(ctx context.Context, login string, password string) (*domain.User, error) {
	user, err := s.repo.GetUser(ctx, login, password)
	if err != nil {
		return nil, fmt.Errorf("getting user: %w", err)
	}
	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User) error {
	err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("creating user: %w", err)
	}
	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, username string) error {
	if err := s.repo.DeleteUser(ctx, username); err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}
	return nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("getting user by username: %w", err)
	}
	return user, nil
}

func (s *UserService) GenerateToken(ctx context.Context, login string, password string) (string, error) {
	user, err := s.GetUser(ctx, login, password)
	if err != nil {
		return "", fmt.Errorf("get user: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserClaims: UserInfo{
			UserID: user.ID,
			Role:   user.Role,
		},
	})

	signedToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signedToken, nil
}

func (s *UserService) ParseToken(accessToken string) (*UserInfo, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*tokenClaims); ok && token.Valid {
		return &claims.UserClaims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *UserService) UpdateUser(ctx context.Context, username string, updatedUser *domain.User) error {
	_, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("error fetching user: %w", err)
	}

	err = s.repo.UpdateUser(ctx, username, updatedUser)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}
