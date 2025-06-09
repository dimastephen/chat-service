package env

import (
	"errors"
	"github.com/dimastephen/chatServer/internal/config"
	"net"
	"os"
)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"
)

type grpcConfig struct {
	host string
	port string
}

func (c *grpcConfig) Address() string {
	return net.JoinHostPort(c.host, c.port)
}

func NewGRPCConfig() (config.GRPCConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("Failed to parse host_name from local.env")
	}
	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("Failed to parse port_name from local.env")
	}
	return &grpcConfig{
		host: host,
		port: port,
	}, nil
}
