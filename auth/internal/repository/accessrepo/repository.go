package accessrepo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dimastephen/auth/internal/config"
	"github.com/dimastephen/auth/internal/logger"
	"github.com/dimastephen/auth/internal/models"
	"github.com/dimastephen/auth/internal/repository"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"os"
	"slices"
)

type roleMapping map[string][]string

type accessRepository struct {
	client *redis.Client
}

func NewAccessRepository(ctx context.Context, config config.RedisConfig) (repository.AccessRepository, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Address(),
		Password: config.Password(),
		DB:       config.DB(),
	})

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Error("failed to ping redis", zap.Error(err))
		return nil, err
	}

	return &accessRepository{client: client}, nil
}

func (a *accessRepository) ReadRoles(ctx context.Context, endpoint string, info *models.UserClaims) error {
	key := "access:" + endpoint
	roles, err := a.client.SMembers(ctx, key).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return err
	}
	if len(roles) == 0 {
		return nil
	}
	ok := slices.Contains(roles, info.Role)
	if !ok {
		return errors.New("access denied")
	} else {
		return nil
	}
}

func (a *accessRepository) Initialize(ctx context.Context, filepath string) error {
	file, err := os.ReadFile(filepath)
	if err != nil {
		logger.Error("failed to read json with roles", zap.Error(err))
		return fmt.Errorf("failed to read json with role mapping: %v", err)
	}

	var mapping roleMapping

	err = json.Unmarshal(file, &mapping)
	if err != nil {
		logger.Error("failed to unmarshall json with roles", zap.Error(err))
		return fmt.Errorf("failed to unmarshall json: %v", err)
	}

	for endpoint, roles := range mapping {
		key := "access:" + endpoint
		existingRoles, err := a.client.SMembers(ctx, key).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			return fmt.Errorf("failed to get roles for %v: %v", endpoint, err)
		}

		for _, role := range roles {
			if !slices.Contains(existingRoles, role) {
				if err := a.client.SAdd(ctx, key, role).Err(); err != nil {
					return fmt.Errorf("failed to add roles: %v", err)
				}
			}
		}
	}
	return nil
}
