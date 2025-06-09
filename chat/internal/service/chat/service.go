package chatService

import (
	"context"
	"github.com/dimastephen/chatServer/internal/client/db"
	"github.com/dimastephen/chatServer/internal/model"
	"github.com/dimastephen/chatServer/internal/repository"
	def "github.com/dimastephen/chatServer/internal/service"
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
	err = s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var errTx error
		errTx = s.noteRepository.Delete(ctx, info)
		if errTx != nil {
			return errTx
		}

		errTx = s.noteRepository.LogAction(ctx, nil, info)
		if errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
