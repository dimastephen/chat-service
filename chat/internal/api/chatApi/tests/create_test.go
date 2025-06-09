package tests

import (
	"context"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/dimastephen/chatServer/internal/api/chatApi"
	"github.com/dimastephen/chatServer/internal/model"
	"github.com/dimastephen/chatServer/internal/service"
	"github.com/dimastephen/chatServer/internal/service/mocks"
	"github.com/dimastephen/chatServer/pkg/chatServerV1"
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreate(t *testing.T) {

	type chatServiceMockFunc func(mc *minimock.Controller) service.Service

	type args struct {
		ctx context.Context
		r   *chatServerV1.CreateRequest
	}

	var (
		mc         = minimock.NewController(t)
		ctx        = context.Background()
		usernames  = gofakeit.ProductAudience()
		id         = gofakeit.Int64()
		serviceErr = errors.New("err")

		req = &chatServerV1.CreateRequest{
			Usernames: usernames,
		}

		res = &chatServerV1.CreateResponse{
			Id: int64(id),
		}

		info = &model.CreateInfo{
			Usernames: usernames,
		}
	)

	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *chatServerV1.CreateResponse
		err             error
		chatServiceMock chatServiceMockFunc
	}{{
		name: "Success case",
		args: args{
			ctx: ctx,
			r:   req,
		},
		want: res,
		err:  nil,
		chatServiceMock: func(mc *minimock.Controller) service.Service {
			mock := mocks.NewServiceMock(mc)
			mock.CreateMock.Expect(ctx, info).Return(&model.CreateInfo{Usernames: usernames, Id: id}, nil)
			return mock
		},
	},
		{
			name: "Service err",
			args: args{ctx: ctx, r: req},
			want: nil,
			err:  serviceErr,
			chatServiceMock: func(mc *minimock.Controller) service.Service {
				mock := mocks.NewServiceMock(mc)
				mock.CreateMock.Expect(ctx, info).Return(nil, serviceErr)
				return mock
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ChatService := tt.chatServiceMock(mc)
			api := chatApi.NewImplementation(ChatService)

			resp, err := api.Create(tt.args.ctx, tt.args.r)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resp)

		})
	}
}
