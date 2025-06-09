package chatApi

import (
	"context"
	"github.com/dimastephen/chatServer/internal/converter"
	"github.com/dimastephen/chatServer/pkg/chatServerV1"
	"log"
)

func (i *Implementation) Create(ctx context.Context, r *chatServerV1.CreateRequest) (*chatServerV1.CreateResponse, error) {
	info, err := i.chatService.Create(ctx, converter.ToCreateInfoFromDesc(r))
	if err != nil {
		return nil, err
	}
	log.Print("inserted id with", info.Id)
	return converter.ToDescFromCreateInfo(info), nil
}
