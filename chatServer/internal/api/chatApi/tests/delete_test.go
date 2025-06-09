package tests

import (
	"context"
	"errors"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/dimastephen/chatServer/internal/api/chatApi"
	"github.com/dimastephen/chatServer/internal/model"
	"github.com/dimastephen/chatServer/internal/service"
	"github.com/dimastephen/chatServer/internal/service/mocks"
	"github.com/dimastephen/chatServer/pkg/chatServerV1"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDelete(t *testing.T) {
	type ServiceMock func(mc *minimock.Controller) service.Service

	type args struct {
		ctx context.Context
		r   *chatServerV1.DeleteRequest
	}

	var (
		mc         = minimock.NewController(t)
		ctx        = context.Background()
		id         = gofakeit.Int64()
		serviceErr = errors.New("service err")

		req = &chatServerV1.DeleteRequest{
			Id: id,
		}

		res = &chatServerV1.DeleteResponse{}

		info = &model.DeleteInfo{
			Id: id,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *chatServerV1.DeleteResponse
		err             error
		chatServiceMock ServiceMock
	}{
		{
			"Success",
			args{ctx: ctx, r: req},
			res,
			nil,
			func(mc *minimock.Controller) service.Service {
				mock := mocks.NewServiceMock(t)
				mock.DeleteMock.Expect(ctx, info).Return(nil)
				return mock
			},
		},
		{
			"Service err",
			args{ctx: ctx, r: req},
			nil,
			serviceErr,
			func(mc *minimock.Controller) service.Service {
				mock := mocks.NewServiceMock(t)
				mock.DeleteMock.Expect(ctx, info).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chatService := tt.chatServiceMock(mc)
			api := chatApi.NewImplementation(chatService)

			resp, err := api.Delete(tt.args.ctx, tt.args.r)
			require.Equal(t, tt.want, resp)
			require.Equal(t, tt.err, err)
		})
	}

}
