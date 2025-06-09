package authservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/dimastephen/auth/internal/config"
	"github.com/dimastephen/auth/internal/crypto"
	"github.com/dimastephen/auth/internal/jwt"
	"github.com/dimastephen/auth/internal/models"
	"github.com/dimastephen/auth/internal/repository"
	service "github.com/dimastephen/auth/internal/service"
	"time"
)

const (
	refreshDuration = 24 * time.Hour
	accessDuration  = 2 * time.Minute
)

type authservice struct {
	repo   repository.AuthRepository
	secret config.SecretKey
}

func NewAuthService(repo repository.AuthRepository, secret config.SecretKey) service.AuthService {
	return &authservice{repo: repo, secret: secret}
}
func (a *authservice) Register(ctx context.Context, user *models.User) (int, error) {
	if user == nil {
		return 0, errors.New("user model in service is nil")
	}
	hash, err := crypto.HashPassword(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = hash
	id, err := a.repo.Register(ctx, user)
	if err != nil || id == 0 {
		return 0, fmt.Errorf("failed to register user: %w")
	}
	return id, nil
}

func (a *authservice) Login(ctx context.Context, user *models.User) (string, error) {
	if user == nil {
		return "", errors.New("user model in service is nil")
	}
	repoData, err := a.repo.Login(ctx, user)
	if err != nil {
		return "", fmt.Errorf("wrong username or password: %w", err)
	}
	err = crypto.CompareHashAndPassword(user.Password, repoData.Password)
	if err != nil {
		return "", fmt.Errorf("wrong username or password: %w", err)
	}

	tokenStr, err := jwt.GenerateToken(repoData, a.secret.RefreshKey(), refreshDuration)

	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return tokenStr, nil
}

func (a *authservice) GetRefreshToken(_ context.Context, tokenStr string) (string, error) {
	claims, err := jwt.VerifyToken(tokenStr, a.secret.RefreshKey())
	if err != nil {
		return "", fmt.Errorf("failed to verify refresh token: %w", err)
	}
	token, err := jwt.GenerateToken(&models.User{Username: claims.Username, Role: claims.Role}, a.secret.RefreshKey(), refreshDuration)
	if err != nil {
		return "", fmt.Errorf("failed to generate new refreshToken: %w", err)
	}
	return token, nil
}

func (a *authservice) GetAccessToken(_ context.Context, tokenStr string) (string, error) {
	claims, err := jwt.VerifyToken(tokenStr, a.secret.RefreshKey())
	if err != nil {
		return "", fmt.Errorf("failed to verify refresh token: %s", err)
	}
	token, err := jwt.GenerateToken(&models.User{Username: claims.Username, Role: claims.Role}, a.secret.AccessKey(), accessDuration)
	if err != nil {
		return "", fmt.Errorf("Failed to generate new accessToken: %s", err)
	}
	return token, nil
}
