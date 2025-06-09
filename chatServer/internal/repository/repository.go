package repository

import (
	"context"
	"github.com/dimastephen/chatServer/internal/model"
)

type ChatRepository interface {
	Create(ctx context.Context, info *model.CreateInfo) (*model.CreateInfo, error)
	Delete(ctx context.Context, info *model.DeleteInfo) error
	LogAction(ctx context.Context, create *model.CreateInfo, delete *model.DeleteInfo) error
}
