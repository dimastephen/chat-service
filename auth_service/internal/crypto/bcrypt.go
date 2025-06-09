package crypto

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashdPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashdPassword), nil
}

func CompareHashAndPassword(password string, hashedPassword string) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.As(err, &bcrypt.ErrMismatchedHashAndPassword) {
			return errors.New("wrong password")
		}
	}
	return nil
}
