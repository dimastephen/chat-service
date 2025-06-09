package env

import (
	"errors"
	"github.com/dimastephen/auth/internal/config"
	"os"
)

const (
	dsnEnvName = "PG_DSN"
)

type pgConfig struct {
	dsn string
}

func (p *pgConfig) DSN() string {
	return p.dsn
}

func NewPGConfig() (config.PGConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("Failed to parse DSN from env\n")
	}

	return &pgConfig{dsn: dsn}, nil
}
