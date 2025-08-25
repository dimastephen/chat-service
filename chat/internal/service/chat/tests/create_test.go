package tests

import (
	"context"
	"errors"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/dimastephen/chatServer/internal/model"
	"github.com/dimastephen/chatServer/internal/repository"
	"github.com/dimastephen/chatServer/internal/repository/mocks"
	chatService "github.com/dimastephen/chatServer/internal/service/chat"
	"github.com/dimastephen/utils/pkg/db"
	tx "github.com/dimastephen/utils/pkg/db/mocks"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreate(t *testing.T) {
	type mockRepofunc func(mc *minimock.MockController) repository.ChatRepository
	type mockTxfunc func(mc *minimock.MockController) db.TxManager

	var (
		mc        = minimock.NewController(t)
		ctx       = context.Background()
		usernames = gofakeit.NiceColors()
		id        = 5
		req       = &model.CreateInfo{
			Usernames: usernames,
		}
		resp = &model.CreateInfo{
			Usernames: usernames,
			Id:        int64(id),
		}
	)

	type args struct {
		ctx  context.Context
		info *model.CreateInfo
	}

	tests := []struct {
		name     string
		args     args
		repoErr  error
		logErr   error
		wantErr  error
		wantResp *model.CreateInfo
	}{
		{
			name:     "Success",
			args:     args{ctx: ctx, info: req},
			repoErr:  nil,
			logErr:   nil,
			wantErr:  nil,
			wantResp: resp,
		},
		{
			name:     "RepoErr",
			args:     args{ctx: ctx, info: req},
			repoErr:  errors.New("CustomRepoErr"),
			logErr:   nil,
			wantErr:  errors.New("CustomRepoErr"),
			wantResp: nil,
		},
		{
			name:     "LogErr",
			args:     args{ctx: ctx, info: req},
			repoErr:  nil,
			logErr:   errors.New("CustomLogErr"),
			wantErr:  errors.New("CustomLogErr"),
			wantResp: resp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repomock := mocks.NewChatRepositoryMock(mc)
			txmock := tx.NewTxManagerMock(mc)

			repomock.CreateMock.Expect(tt.args.ctx, tt.args.info).Return(tt.wantResp, tt.repoErr)
			if tt.repoErr == nil {
				repomock.LogActionMock.Expect(tt.args.ctx, tt.wantResp, nil).Return(tt.logErr)
			}
			txmock.ReadCommitedMock.Set(func(ctx context.Context, f db.Handler) (err error) {
				return f(ctx)
			})

			service := chatService.NewService(repomock, txmock)
			response, err := service.Create(tt.args.ctx, tt.args.info)
			require.Equal(t, tt.wantErr, err)
			if err == nil {
				require.Equal(t, tt.wantResp.Id, response.Id)
			}
		})
	}
}
