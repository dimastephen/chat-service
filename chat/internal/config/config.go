package config

import "github.com/joho/godotenv"

func Load(configPath string) error {
	err := godotenv.Load(configPath)
	if err != nil {
		return err
	}
	return nil
}

type GRPCConfig interface {
	Address() string
}

type PGConfig interface {
	DSN() string
}

type HTTPConfig interface {
	Address() string
}
