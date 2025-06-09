package auth

import (
	"context"
	desc "github.com/dimastephen/auth/grpc/pkg/authV1"
	"github.com/dimastephen/auth/internal/models"
	"github.com/dimastephen/auth/internal/service"
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

	user := &models.User{Password: password, Username: username}
	id, err := i.authService.Register(ctx, user)
	if err != nil {
		return nil, err
	}
	return &desc.RegisterResponse{Id: int64(id)}, nil
}

func (i *AuthImplementation) Login(ctx context.Context, r *desc.LoginRequest) (*desc.LoginResponse, error) {
	username := r.GetUsername()
	password := r.GetPassword()

	user := &models.User{Password: password, Username: username}
	token, err := i.authService.Login(ctx, user)
	if err != nil {
		return nil, err
	}
	return &desc.LoginResponse{RefreshToken: token}, nil
}

func (i *AuthImplementation) GetRefreshToken(ctx context.Context, r *desc.GetRefreshTokenRequest) (*desc.GetRefreshTokenResponse, error) {
	tokenStr := r.GetRefreshToken()
	newToken, err := i.authService.GetRefreshToken(ctx, tokenStr)
	if err != nil {
		return nil, err
	}

	return &desc.GetRefreshTokenResponse{RefreshToken: newToken}, nil
}

func (i *AuthImplementation) GetAccessToken(ctx context.Context, r *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {
	refreshToken := r.GetRefreshToken()
	accessToken, err := i.authService.GetAccessToken(ctx, refreshToken)

	if err != nil {
		return nil, err
	}

	return &desc.GetAccessTokenResponse{AccessToken: accessToken}, nil
}
