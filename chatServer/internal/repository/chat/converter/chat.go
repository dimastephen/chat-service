package converter

import (
	"github.com/dimastephen/chatServer/internal/repository/chat/model"
	"github.com/dimastephen/chatServer/pkg/chatServerV1"
)

func FromChatToResponse(chat *model.NewChat) *chatServerV1.CreateResponse {
	return &chatServerV1.CreateResponse{
		Id: int64(chat.Id),
	}
}
