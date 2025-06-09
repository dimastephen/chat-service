package tests

import (
	"context"
	"errors"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/dimastephen/chatServer/internal/client/db"
	txmock "github.com/dimastephen/chatServer/internal/client/db/mocks"
	"github.com/dimastephen/chatServer/internal/model"
	"github.com/dimastephen/chatServer/internal/repository"
	"github.com/dimastephen/chatServer/internal/repository/mocks"
	chatService "github.com/dimastephen/chatServer/internal/service/chat"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDelete(t *testing.T) {
	type repoMockFunc func(mc *minimock.Controller) repository.ChatRepository
	type clientMock func(mc *minimock.Controller) db.Client
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)
		id  = gofakeit.Int64()
		req = &model.DeleteInfo{
			Id: id,
		}
		logErr  = errors.New("Log error")
		repoErr = errors.New("Repo error")
	)

	type args struct {
		ctx context.Context
		r   *model.DeleteInfo
	}

	tests := []struct {
		name    string
		args    args
		repoErr error
		logErr  error
		wantErr error
	}{
		{
			name:    "Success",
			args:    args{ctx, req},
			repoErr: nil,
			logErr:  nil,
			wantErr: nil,
		},
		{
			name:    "Log Error",
			args:    args{ctx, req},
			repoErr: nil,
			logErr:  logErr,
			wantErr: logErr,
		},
		{
			name:    "RepoErr",
			args:    args{ctx, req},
			repoErr: repoErr,
			logErr:  nil,
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chatRepo := mocks.NewChatRepositoryMock(mc)
			txMock := txmock.NewTxManagerMock(mc)

			chatRepo.DeleteMock.Expect(tt.args.ctx, tt.args.r).Return(tt.repoErr)
			if tt.repoErr == nil {
				chatRepo.LogActionMock.Expect(tt.args.ctx, nil, tt.args.r).Return(tt.logErr)
			}

			txMock.ReadCommitedMock.Set(func(ctx context.Context, f db.Handler) error {
				return f(ctx)
			})

			service := chatService.NewService(chatRepo, txMock)
			err := service.Delete(tt.args.ctx, tt.args.r)
			require.Equal(t, tt.wantErr, err)
		})
	}

}
