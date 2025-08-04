package app

import (
	"context"
	"github.com/dimastephen/chatServer/internal/config"
	"github.com/dimastephen/chatServer/internal/interceptor"
	"github.com/dimastephen/chatServer/internal/logger"
	desc "github.com/dimastephen/chatServer/pkg/chatServerV1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
)

type App struct {
	provider   *ServiceProvider
	grpcServer *grpc.Server
	httpServer *http.Server
}

func NewApp(ctx context.Context, configPath string, level string) (*App, error) {
	a := App{}
	err := a.initDeps(ctx, configPath, level)
	if err != nil {
		log.Fatalf("error creating app %s", err.Error())
	}
	return &a, nil
}

func (a *App) RunGRPCServer() error {
	logger.Info("Running GRPC Server on: %v", zap.String("address", a.provider.GRPCConfig().Address()))

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

func (a *App) RunHTTPServer() error {
	logger.Info("Running HTTP Server on: %v", zap.String("address", a.provider.HTTPConfig().Address()))
	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (a *App) initDeps(ctx context.Context, configPath string, level string) error {
	err := a.initConfig(ctx, configPath)
	if err != nil {
		log.Fatalf("error loading config: %s", err.Error())
	}
	logger.Init(getCore(getLevel(level)))
	inits := []func(context.Context) error{
		a.initServiceProvider,
		a.initServer,
		a.InitHTTP,
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
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()), grpc.ChainUnaryInterceptor(interceptor.LogInterceptor, interceptor.ValidateInterceptor))

	reflection.Register(a.grpcServer)

	desc.RegisterChatServerServer(a.grpcServer, a.provider.ChatImpl(ctx))

	return nil
}

func (a *App) InitHTTP(ctx context.Context) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	err := desc.RegisterChatServerHandlerFromEndpoint(ctx, mux, a.ServiceProvider().GRPCConfig().Address(), opts)
	if err != nil {
		return err
	}
	a.httpServer = &http.Server{
		Addr:    a.ServiceProvider().HTTPConfig().Address(),
		Handler: mux,
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.provider = NewServiceProvider()
	return nil
}

func (a *App) ServiceProvider() *ServiceProvider {
	return a.provider
}

func (a *App) Run() error {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		err := a.RunGRPCServer()
		if err != nil {
			log.Fatalf("Error running GRPC server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		err := a.RunHTTPServer()
		if err != nil {
			log.Fatalf("Error running HTTP server: %v", err)
		}
	}()
	wg.Wait()
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
	err := lvl.Set(level)
	if err != nil {
		log.Fatalf("failed to set level to logger: %v", err)
	}
	return zap.NewAtomicLevelAt(lvl)
}
