package repository

import (
	"context"
	"github.com/dimastephen/auth/internal/models"
)

type AuthRepository interface {
	Login(ctx context.Context, user *models.User) (*models.User, error)
	Register(ctx context.Context, user *models.User) (int, error)
}

type AccessRepository interface {
	ReadRoles(ctx context.Context, endpoint string, info *models.UserClaims) error
	Initialize(ctx context.Context, filepath string) error
}
