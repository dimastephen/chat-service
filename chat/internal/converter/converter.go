package converter

import (
	"github.com/dimastephen/chatServer/internal/model"
	desc "github.com/dimastephen/chatServer/pkg/chatServerV1"
)

func ToCreateInfoFromDesc(request *desc.CreateRequest) *model.CreateInfo {
	models := model.CreateInfo{
		Usernames: request.GetUsernames(),
	}
	return &models
}

func ToDescFromCreateInfo(info *model.CreateInfo) *desc.CreateResponse {
	return &desc.CreateResponse{
		Id: info.Id,
	}
}

func ToDeleteInfoFromDesc(request *desc.DeleteRequest) *model.DeleteInfo {
	return &model.DeleteInfo{
		Id: request.GetId(),
	}
}
