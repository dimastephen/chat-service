package service

import (
	"context"
	"github.com/dimastephen/chatServer/internal/model"
)

type Service interface {
	Create(ctx context.Context, createReq *model.CreateInfo) (response *model.CreateInfo, err error)
	Delete(ctx context.Context, response *model.DeleteInfo) (err error)
}
