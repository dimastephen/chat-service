package config

import "github.com/joho/godotenv"

type GRPCConfig interface {
	Address() string
}

type PGConfig interface {
	DSN() string
}

type SecretKey interface {
	AccessKey() []byte
	RefreshKey() []byte
}

type RedisConfig interface {
	Address() string
	Password() string
	DB() int
}

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}
	return nil
}
