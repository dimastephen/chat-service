package env

import (
	"errors"
	"github.com/dimastephen/auth/internal/config"
	"net"
	"os"
)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "AUTH_GRPC_PORT"
)

type grpcConfig struct {
	host string
	port string
}

func (g *grpcConfig) Address() string {
	return net.JoinHostPort(g.host, g.port)
}

func NewGRPCConfig() (config.GRPCConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("Failed to parse host from env\n")
	}
	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("Failed to parse port from env\n")
	}

	return &grpcConfig{host: host, port: port}, nil
}
