package chatService

import (
	"context"
	"github.com/dimastephen/chatServer/internal/logger"
	"github.com/dimastephen/chatServer/internal/model"
	"github.com/dimastephen/chatServer/internal/repository"
	def "github.com/dimastephen/chatServer/internal/service"
	"github.com/dimastephen/utils/pkg/db"
	"go.uber.org/zap"
)

type service struct {
	noteRepository repository.ChatRepository
	txManager      db.TxManager
}

func NewService(noteRepository repository.ChatRepository, txManager db.TxManager) def.Service {
	return &service{
		noteRepository: noteRepository,
		txManager:      txManager,
	}
}

func (s *service) Create(ctx context.Context, info *model.CreateInfo) (*model.CreateInfo, error) {
	resp := &model.CreateInfo{}
	logger.Debug("Creating chat with", zap.Any("usernames", info.Usernames))
	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var errTx error
		resp, errTx = s.noteRepository.Create(ctx, info)
		if errTx != nil {
			return errTx
		}

		errTx = s.noteRepository.LogAction(ctx, resp, nil)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *service) Delete(ctx context.Context, info *model.DeleteInfo) (err error) {
	logger.Debug("Deleting chat with", zap.Int("id", int(info.Id)))
	err = s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var errTx error
		errTx = s.noteRepository.Delete(ctx, info)
		if errTx != nil {
			logger.Error("Error in Tx", zap.Error(errTx))
			return errTx
		}

		errTx = s.noteRepository.LogAction(ctx, nil, info)
		if errTx != nil {
			logger.Error("Error in Tx", zap.Error(errTx))
			return errTx
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
