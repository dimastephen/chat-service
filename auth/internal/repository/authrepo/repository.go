package authrepo

import (
	"context"
	"errors"
	"github.com/dimastephen/auth/internal/logger"
	"github.com/dimastephen/auth/internal/models"
	"github.com/dimastephen/auth/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type authRepository struct {
	dbc *pgxpool.Pool
}

func NewAuthRepository(ctx context.Context, dsn string) (repository.AuthRepository, error) {
	dbc, err := pgxpool.New(ctx, dsn)
	if err != nil {
		logger.Error("Failed to connect ot pgxpool", zap.Error(err))
		return nil, err
	}
	err = dbc.Ping(ctx)
	if err != nil {
		logger.Error("Failed to ping db", zap.Error(err))
		return nil, err
	}
	return &authRepository{dbc: dbc}, nil
}

func (a *authRepository) Login(ctx context.Context, user *models.User) (*models.User, error) {
	query := "SELECT password,role FROM users WHERE username = ($1)"
	logger.Debug("Login query", zap.String("query", query))
	resp := models.User{}
	resp.Username = user.Username
	err := a.dbc.QueryRow(ctx, query, user.Username).Scan(&resp.Password, &resp.Role)
	if errors.As(err, &pgx.ErrNoRows) {
		return nil, errors.New("wrong username or password")
	}
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (a *authRepository) Register(ctx context.Context, user *models.User) (int, error) {
	query := "INSERT INTO users(username,password) VALUES ($1,$2) RETURNING id"
	logger.Debug("Register query", zap.String("query", query))

	var id int
	err := a.dbc.QueryRow(ctx, query, user.Username, user.Password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
