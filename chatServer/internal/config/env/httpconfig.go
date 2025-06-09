package env

import (
	"github.com/dimastephen/chatServer/internal/config"
	"github.com/pkg/errors"
	"net"
	"os"
)

const (
	httpHostEnvName = "HTTP_HOST"
	httpPortEnvName = "HTTP_PORT"
)

type httpConfig struct {
	host string
	port string
}

func (h *httpConfig) Address() string {
	return net.JoinHostPort(h.host, h.port)
}

func NewHttpConfig() (config.HTTPConfig, error) {
	host := os.Getenv(httpHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("Failed to parse http Host from env")
	}
	port := os.Getenv(httpPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("Failed to parse http Post from env")
	}
	return &httpConfig{
		port: port,
		host: host,
	}, nil
}
