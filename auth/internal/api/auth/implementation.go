package auth

import (
	"context"
	"github.com/dimastephen/auth/internal/logger"
	"github.com/dimastephen/auth/internal/models"
	"github.com/dimastephen/auth/internal/service"
	desc "github.com/dimastephen/auth/pkg/authV1"
	"go.uber.org/zap"
)

type AuthImplementation struct {
	desc.UnimplementedAuthServer
	authService service.AuthService
}

func NewImplementation(service service.AuthService) *AuthImplementation {
	return &AuthImplementation{authService: service}
}

func (i *AuthImplementation) Register(ctx context.Context, r *desc.RegisterRequest) (*desc.RegisterResponse, error) {
	username := r.GetUsername()
	password := r.GetPassword()
	logger.Debug("Registering new user", zap.String("username", username))

	user := &models.User{Password: password, Username: username}
	id, err := i.authService.Register(ctx, user)
	if err != nil {
		logger.Error("Error registering user", zap.Error(err))
		return nil, err
	}
	logger.Info("Successfuly registered user", zap.Int("id", id))
	return &desc.RegisterResponse{Id: int64(id)}, nil
}

func (i *AuthImplementation) Login(ctx context.Context, r *desc.LoginRequest) (*desc.LoginResponse, error) {
	username := r.GetUsername()
	password := r.GetPassword()
	logger.Debug("User try to log in", zap.String("username", username))

	user := &models.User{Password: password, Username: username}
	token, err := i.authService.Login(ctx, user)
	if err != nil {
		logger.Error("Failed to log in", zap.String("username", username), zap.Error(err))
		return nil, err
	}

	logger.Info("User is successfully logged in", zap.String("username", username))
	return &desc.LoginResponse{RefreshToken: token}, nil
}

func (i *AuthImplementation) GetRefreshToken(ctx context.Context, r *desc.GetRefreshTokenRequest) (*desc.GetRefreshTokenResponse, error) {
	tokenStr := r.GetRefreshToken()
	newToken, err := i.authService.GetRefreshToken(ctx, tokenStr)
	if err != nil {
		return nil, err
	}

	logger.Info("User successfully got new refresh token")
	return &desc.GetRefreshTokenResponse{RefreshToken: newToken}, nil
}

func (i *AuthImplementation) GetAccessToken(ctx context.Context, r *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {
	refreshToken := r.GetRefreshToken()
	accessToken, err := i.authService.GetAccessToken(ctx, refreshToken)

	if err != nil {
		return nil, err
	}

	logger.Info("User successfully got new access token")
	return &desc.GetAccessTokenResponse{AccessToken: accessToken}, nil
}
