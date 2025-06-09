package chatApi

import (
	"context"
	"github.com/dimastephen/chatServer/internal/converter"
	"github.com/dimastephen/chatServer/pkg/chatServerV1"
)

func (i *Implementation) Delete(ctx context.Context, r *chatServerV1.DeleteRequest) (*chatServerV1.DeleteResponse, error) {
	err := i.chatService.Delete(ctx, converter.ToDeleteInfoFromDesc(r))
	if err != nil {
		return nil, err
	}
	return &chatServerV1.DeleteResponse{}, nil
}
