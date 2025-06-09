package app

import (
	"context"
	"github.com/dimastephen/chatServer/internal/api/chatApi"
	"github.com/dimastephen/chatServer/internal/client/db"
	"github.com/dimastephen/chatServer/internal/client/db/pg"
	"github.com/dimastephen/chatServer/internal/client/db/transaction"
	"github.com/dimastephen/chatServer/internal/config"
	"github.com/dimastephen/chatServer/internal/config/env"
	"github.com/dimastephen/chatServer/internal/repository"
	chatRepository "github.com/dimastephen/chatServer/internal/repository/chat"
	"github.com/dimastephen/chatServer/internal/service"
	chatService "github.com/dimastephen/chatServer/internal/service/chat"
	"log"
)

type ServiceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig
	httpConfig config.HTTPConfig

	chatService service.Service

	dbClient  db.Client
	txManager db.TxManager
	chatRepo  repository.ChatRepository

	chatImpl *chatApi.Implementation
}

func NewServiceProvider() *ServiceProvider {
	return &ServiceProvider{}
}

func (s *ServiceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatal("error getting pgConfig in ServiceProvider")
		}
		s.pgConfig = cfg
	}
	return s.pgConfig
}

func (s *ServiceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("error getting grpcConfig in ServiceProvider %s", err.Error())
		}
		s.grpcConfig = cfg
	}
	return s.grpcConfig
}

func (s *ServiceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := env.NewHttpConfig()
		if err != nil {
			log.Fatalf("error getting httpConfig in ServiceProvider %s", err.Error())
		}
		s.httpConfig = cfg
	}
	return s.httpConfig
}

func (s *ServiceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		client, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("Error connecting DB %v", err)
		}
		s.dbClient = client
		err = s.dbClient.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("Error pinging db %v", err)
		}
	}
	return s.dbClient
}

func (s *ServiceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		tx := transaction.NewTransactionManager(s.DBClient(ctx).DB())
		s.txManager = tx
	}
	return s.txManager

}

func (s *ServiceProvider) ChatRepo(ctx context.Context) repository.ChatRepository {
	if s.chatRepo == nil {
		repo := chatRepository.NewRepository(s.DBClient(ctx))
		s.chatRepo = repo
	}
	return s.chatRepo
}

func (s *ServiceProvider) ChatService(ctx context.Context) service.Service {
	if s.chatService == nil {
		serv := chatService.NewService(s.ChatRepo(ctx), s.TxManager(ctx))
		s.chatService = serv
	}
	return s.chatService
}

func (s *ServiceProvider) ChatImpl(ctx context.Context) *chatApi.Implementation {
	if s.chatImpl == nil {
		impl := chatApi.NewImplementation(s.ChatService(ctx))
		s.chatImpl = impl
	}
	return s.chatImpl
}
