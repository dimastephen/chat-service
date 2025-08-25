package env

import (
	"errors"
	"github.com/dimastephen/chatServer/internal/config"
	"net"
	"os"
)

const (
	swaggerHostEnvName = "HTTP_HOST"
	swaggerPortEnvName = "SWAGGER_PORT"
)

type swaggerConfig struct {
	host string
	port string
}

func NewSwaggerConfig() (config.SwaggerConfig, error) {
	host := os.Getenv(swaggerHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("failed to parse swagger env")
	}
	port := os.Getenv(swaggerPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("failed to parse swagger env")
	}
	return &swaggerConfig{host: host, port: port}, nil
}

func (s *swaggerConfig) Address() string {
	return net.JoinHostPort(s.host, s.port)
}
