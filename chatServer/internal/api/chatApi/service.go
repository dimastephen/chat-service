package chatApi

import (
	"github.com/dimastephen/chatServer/internal/service"
	desc "github.com/dimastephen/chatServer/pkg/chatServerV1"
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
