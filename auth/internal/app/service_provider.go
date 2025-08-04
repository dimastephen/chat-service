package app

import (
	"context"
	"github.com/dimastephen/auth/internal/api/access"
	"github.com/dimastephen/auth/internal/api/auth"
	"github.com/dimastephen/auth/internal/config"
	"github.com/dimastephen/auth/internal/config/env"
	"github.com/dimastephen/auth/internal/repository"
	"github.com/dimastephen/auth/internal/repository/accessrepo"
	"github.com/dimastephen/auth/internal/repository/authrepo"
	"github.com/dimastephen/auth/internal/service"
	"github.com/dimastephen/auth/internal/service/accessservice"
	"github.com/dimastephen/auth/internal/service/authservice"
	"log"
)

type ServiceProvider struct {
	pgConfig    config.PGConfig
	grpcConfig  config.GRPCConfig
	secretkey   config.SecretKey
	redisConfig config.RedisConfig

	authService   service.AuthService
	accessService service.AccessService

	authRepo   repository.AuthRepository
	accessRepo repository.AccessRepository

	authImplementation   *auth.AuthImplementation
	accessImplementation *access.AccessImplementation
}

func NewServiceProvider() *ServiceProvider {
	return &ServiceProvider{}
}

func (s *ServiceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatal("failed to parse PGconfig\n")
		}
		s.pgConfig = cfg
	}
	return s.pgConfig
}

func (s *ServiceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatal(err)
		}
		s.grpcConfig = cfg
	}
	return s.grpcConfig
}

func (s *ServiceProvider) SecretKey() config.SecretKey {
	if s.secretkey == nil {
		secret, err := env.NewSecretConfig()
		if err != nil {
			log.Fatal(err)
		}
		s.secretkey = secret
	}
	return s.secretkey
}

func (s *ServiceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		redisConfig, err := env.NewRedisConfig()
		if err != nil {
			log.Fatal(err)
		}
		s.redisConfig = redisConfig
	}
	return s.redisConfig
}

func (s *ServiceProvider) AuthRepo(ctx context.Context) repository.AuthRepository {
	if s.authRepo == nil {
		dbc, err := authrepo.NewAuthRepository(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatal(err)
		}
		s.authRepo = dbc
	}
	return s.authRepo
}

func (s *ServiceProvider) AccessRepo(ctx context.Context) repository.AccessRepository {
	if s.accessRepo == nil {
		redis, err := accessrepo.NewAccessRepository(ctx, s.RedisConfig())
		if err != nil {
			log.Fatal(err)
		}
		s.accessRepo = redis
		err = s.accessRepo.Initialize(ctx, "roles.json")
	}
	return s.accessRepo
}

func (s *ServiceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		serv := authservice.NewAuthService(s.AuthRepo(ctx), s.SecretKey())
		s.authService = serv
	}
	return s.authService
}

func (s *ServiceProvider) AccessService(ctx context.Context) service.AccessService {
	if s.accessService == nil {
		serv := accessservice.NewAccessService(s.AccessRepo(ctx), s.SecretKey())
		s.accessService = serv
	}
	return s.accessService
}

func (s *ServiceProvider) AuthImplementation(ctx context.Context) *auth.AuthImplementation {
	if s.authImplementation == nil {
		impl := auth.NewImplementation(s.AuthService(ctx))
		s.authImplementation = impl
	}
	return s.authImplementation
}

func (s *ServiceProvider) AccessImplementation(ctx context.Context) *access.AccessImplementation {
	if s.accessImplementation == nil {
		impl := access.NewImplementation(s.AccessService(ctx))
		s.accessImplementation = impl
	}
	return s.accessImplementation
}
