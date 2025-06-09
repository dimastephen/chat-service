package env

import (
	"errors"
	"github.com/dimastephen/auth/internal/config"
	"os"
	"strconv"
)

const (
	redis_address  = "REDIS_ADDRESS"
	redis_password = "REDIS_PASSWORD"
	redis_db       = "REDIS_DB"
)

type redis struct {
	db       int
	address  string
	password string
}

func (r *redis) DB() int {
	return r.db
}

func (r *redis) Address() string {
	return r.address
}

func (r *redis) Password() string {
	return r.password
}

func NewRedisConfig() (config.RedisConfig, error) {
	address := os.Getenv(redis_address)
	if len(address) == 0 {
		return nil, errors.New("failed to parse redis address in config")
	}

	db := os.Getenv(redis_db)
	if len(db) == 0 {
		return nil, errors.New("failed to parse redis db in config")
	}
	intdb, err := strconv.Atoi(db)
	if err != nil {
		return nil, errors.New("failed to parse redis int db in config")
	}

	password := os.Getenv(redis_password)
	if len(password) == 0 {
		return nil, errors.New("failed to parse redis password in config")
	}

	return &redis{
		db:       intdb,
		address:  address,
		password: password,
	}, nil
}
