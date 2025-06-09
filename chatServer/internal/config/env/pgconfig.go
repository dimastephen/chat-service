package env

import (
	"errors"
	"github.com/dimastephen/chatServer/internal/config"
	"os"
)

const (
	dsnEnvName = "PG_DSN"
)

type pgConfig struct {
	dsn string
}

func (c *pgConfig) DSN() string {
	return c.dsn
}

func NewPGConfig() (config.PGConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("Failed to parse dsn from local.env")
	}
	return &pgConfig{
		dsn: dsn,
	}, nil
}
