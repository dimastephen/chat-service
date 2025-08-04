package accessservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/dimastephen/auth/internal/config"
	"github.com/dimastephen/auth/internal/jwt"
	"github.com/dimastephen/auth/internal/logger"
	"github.com/dimastephen/auth/internal/repository"
	"github.com/dimastephen/auth/internal/service"
	"google.golang.org/grpc/metadata"
	"strings"
)

type accessService struct {
	accessRepo repository.AccessRepository
	secret     config.SecretKey
}

func NewAccessService(accessRepo repository.AccessRepository, secret config.SecretKey) service.AccessService {
	return &accessService{accessRepo: accessRepo, secret: secret}
}

func (a *accessService) Check(ctx context.Context, endpoint string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Error("Faild to get metadata")
		return errors.New("failed to get metadata from context")
	}
	data := md["authorization"]
	if len(data) == 0 {
		logger.Error("Failed to get token from metadata header")
		return errors.New("failed to find token in metadata")
	}

	if !strings.HasPrefix(data[0], "Bearer ") {
		logger.Error("Unknown auth format, not Bearer")
		return errors.New("unknown authorization format")
	}
	accessToken := strings.TrimPrefix(data[0], "Bearer ")

	claims, err := jwt.VerifyToken(accessToken, a.secret.AccessKey())
	if err != nil {
		return fmt.Errorf("access token is invalid: %s", err.Error())
	}

	err = a.accessRepo.ReadRoles(ctx, endpoint, claims)
	if err != nil {
		return fmt.Errorf("access denied: %s", err.Error())
	}
	return nil
}
