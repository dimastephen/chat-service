package app

import (
	"context"
	"github.com/dimastephen/auth/internal/config"
	"github.com/dimastephen/auth/internal/interceptor"
	"github.com/dimastephen/auth/internal/logger"
	desc2 "github.com/dimastephen/auth/pkg/access_v1"
	desc "github.com/dimastephen/auth/pkg/authV1"
	"github.com/dimastephen/utils/pkg/rate_limiter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"time"
)

type App struct {
	provider   *ServiceProvider
	grpcServer *grpc.Server
}

func NewApp(ctx context.Context, configpath string, level string) (*App, error) {
	a := App{}
	err := a.initDeps(ctx, configpath, level)
	if err != nil {
		log.Fatalf("error creating app %s", err.Error())
	}
	return &a, nil
}

func (a *App) initDeps(ctx context.Context, configPath string, level string) error {
	err := a.initConfig(ctx, configPath)
	log.Printf("configpath: %v", configPath)
	if err != nil {
		log.Fatalf("Failed to init config: %v", err.Error())
	}
	logger.Init(getCore(getLevel(level)))

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
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()), grpc.ChainUnaryInterceptor(interceptor.NewRateLimiterInterceptor(limiter).Unary, interceptor.LogInterceptor, interceptor.ValidateInterceptor))
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
	logger.Info("Running GRPC Server", zap.String("port", a.provider.GRPCConfig().Address()))

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

func getCore(level zap.AtomicLevel) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.TimeKey = "timestamp"
	developmentCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	developmentCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)

	return zapcore.NewCore(consoleEncoder, stdout, level)
}

func getLevel(level string) zap.AtomicLevel {
	var lvl zapcore.Level
	log.Printf("log level for now: %v", level)
	if err := lvl.Set(level); err != nil {
		log.Fatalf("failed to set logger lvl: %v", err)
	}
	return zap.NewAtomicLevelAt(lvl)
}
