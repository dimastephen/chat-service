package chatApi

import (
	"context"
	"github.com/dimastephen/chatServer/internal/converter"
	"github.com/dimastephen/chatServer/internal/service"
	"github.com/dimastephen/chatServer/pkg/chatServerV1"
	desc "github.com/dimastephen/chatServer/pkg/chatServerV1"
	"log"
)

type Implementation struct {
	desc.UnimplementedChatServerServer
	chatService service.Service
}

func NewImplementation(serv service.Service) *Implementation {
	return &Implementation{
		chatService: serv,
	}
}

func (i *Implementation) Delete(ctx context.Context, r *chatServerV1.DeleteRequest) (*chatServerV1.DeleteResponse, error) {
	err := i.chatService.Delete(ctx, converter.ToDeleteInfoFromDesc(r))
	if err != nil {
		return nil, err
	}
	return &chatServerV1.DeleteResponse{}, nil
}

func (i *Implementation) Create(ctx context.Context, r *chatServerV1.CreateRequest) (*chatServerV1.CreateResponse, error) {
	info, err := i.chatService.Create(ctx, converter.ToCreateInfoFromDesc(r))
	if err != nil {
		return nil, err
	}
	log.Print("inserted id with", info.Id)
	return converter.ToDescFromCreateInfo(info), nil
}
