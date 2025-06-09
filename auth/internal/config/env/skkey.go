package env

import (
	"errors"
	"github.com/dimastephen/auth/internal/config"
	"os"
)

const (
	secretkeyrefresh = "SECRET_KEY_REFRESH"
	secretkeyaccess  = "SECRET_KEY_ACCESS"
)

type secretKey struct {
	refreshkey []byte
	accesskey  []byte
}

func (s *secretKey) RefreshKey() []byte { return s.refreshkey }

func (s *secretKey) AccessKey() []byte { return s.accesskey }

func NewSecretConfig() (config.SecretKey, error) {
	refreshKey := os.Getenv(secretkeyrefresh)
	accessKey := os.Getenv(secretkeyaccess)
	if len(refreshKey) == 0 {
		return nil, errors.New("failed to parse refreshsecretKey from env")
	}
	if len(accessKey) == 0 {
		return nil, errors.New("failed to parse accesssecretKey from env")
	}
	return &secretKey{
		refreshkey: []byte(refreshKey),
		accesskey:  []byte(accessKey),
	}, nil
}
