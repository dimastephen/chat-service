package app

import (
	"context"
	desc2 "github.com/dimastephen/auth/grpc/pkg/access_v1"
	desc "github.com/dimastephen/auth/grpc/pkg/authV1"
	"github.com/dimastephen/auth/internal/config"
	"github.com/dimastephen/auth/internal/interceptor"
	"github.com/dimastephen/chat-service/common/rate_limiter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

type App struct {
	provider   *ServiceProvider
	grpcServer *grpc.Server
}

func NewApp(ctx context.Context, configpath string) (*App, error) {
	a := App{}
	err := a.initDeps(ctx, configpath)
	if err != nil {
		log.Fatalf("error creating app %s", err.Error())
	}
	return &a, nil
}

func (a *App) initDeps(ctx context.Context, configpath string) error {
	err := a.initConfig(ctx, configpath)
	if err != nil {
		log.Fatalf("Failed to init config: %v", err.Error())
	}

	inits := []func(ctx context.Context) error{
		a.initServiceProvider,
		a.initServer,
	}

	for _, fun := range inits {
		err = fun(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initConfig(ctx context.Context, path string) error {
	err := config.Load(path)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) initServer(ctx context.Context) error {
	limiter := rate_limiter.NewTokenBucketLimiter(ctx, 10, time.Second)
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()), grpc.UnaryInterceptor(interceptor.NewRateLimiterInterceptor(limiter).Unary))
	reflection.Register(a.grpcServer)
	desc.RegisterAuthServer(a.grpcServer, a.provider.AuthImplementation(ctx))
	desc2.RegisterAccessServer(a.grpcServer, a.provider.AccessImplementation(ctx))
	return nil
}

func (a *App) initServiceProvider(ctx context.Context) error {
	a.provider = NewServiceProvider()
	return nil
}

func (a *App) Run() error {
	log.Printf("Running GRPC Server on: %v", a.provider.GRPCConfig().Address())

	list, err := net.Listen("tcp", a.provider.GRPCConfig().Address())
	if err != nil {
		return err
	}
	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}
