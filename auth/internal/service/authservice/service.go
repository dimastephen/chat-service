package authservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/dimastephen/auth/internal/config"
	"github.com/dimastephen/auth/internal/crypto"
	"github.com/dimastephen/auth/internal/jwt"
	"github.com/dimastephen/auth/internal/logger"
	"github.com/dimastephen/auth/internal/models"
	"github.com/dimastephen/auth/internal/repository"
	service "github.com/dimastephen/auth/internal/service"
	"go.uber.org/zap"
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
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if user == nil {
		logger.Error("User model in service is nil")
		return 0, errors.New("user model in service is nil")
	}
	hash, err := crypto.HashPassword(user.Password)
	if err != nil {
		logger.Error("Failed to hash password")
		return 0, err
	}
	user.Password = hash
	id, err := a.repo.Register(ctx, user)
	if err != nil || id == 0 {
		logger.Error("Failed to register user", zap.String("username", user.Username))
		return 0, fmt.Errorf("failed to register user")
	}
	return id, nil
}

func (a *authservice) Login(ctx context.Context, user *models.User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if user == nil {
		logger.Error("user model in service is nil")
		return "", errors.New("user model in service is nil")
	}
	repoData, err := a.repo.Login(ctx, user)
	if err != nil {
		logger.Error("Wrong username for", zap.String("username", user.Username))
		return "", fmt.Errorf("wrong username or password")
	}
	err = crypto.CompareHashAndPassword(user.Password, repoData.Password)
	if err != nil {
		logger.Error("Wrong password for", zap.String("username", user.Username))
		return "", fmt.Errorf("wrong username or password")
	}

	tokenStr, err := jwt.GenerateToken(repoData, a.secret.RefreshKey(), refreshDuration)

	if err != nil {
		logger.Error("Fail in generating refresh token", zap.Error(err))
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return tokenStr, nil
}

func (a *authservice) GetRefreshToken(ctx context.Context, tokenStr string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	claims, err := jwt.VerifyToken(tokenStr, a.secret.RefreshKey())
	if err != nil {
		logger.Error("Error in verifying old refresh token", zap.Error(err))
		return "", fmt.Errorf("failed to verify refresh token: %w", err)
	}
	token, err := jwt.GenerateToken(&models.User{Username: claims.Username, Role: claims.Role}, a.secret.RefreshKey(), refreshDuration)
	if err != nil {
		logger.Error("Error in generating new refresh token", zap.Error(err))
		return "", fmt.Errorf("failed to generate new refreshToken: %w", err)
	}
	return token, nil
}

func (a *authservice) GetAccessToken(ctx context.Context, tokenStr string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	claims, err := jwt.VerifyToken(tokenStr, a.secret.RefreshKey())
	if err != nil {
		logger.Error("Error in verifying old access token", zap.Error(err))
		return "", fmt.Errorf("failed to verify refresh token: %s", err)
	}
	token, err := jwt.GenerateToken(&models.User{Username: claims.Username, Role: claims.Role}, a.secret.AccessKey(), accessDuration)
	if err != nil {
		logger.Error("Error in generating new access token", zap.Error(err))
		return "", fmt.Errorf("failed to generate new accessToken: %s", err)
	}
	return token, nil
}
