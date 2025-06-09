package service

import (
	"context"
	"github.com/dimastephen/auth/internal/models"
)

type AuthService interface {
	Register(ctx context.Context, user *models.User) (int, error)
	Login(ctx context.Context, user *models.User) (string, error)
	GetRefreshToken(ctx context.Context, token string) (string, error)
	GetAccessToken(ctx context.Context, token string) (string, error)
}

type AccessService interface {
	Check(ctx context.Context, endpoint string) error
}
