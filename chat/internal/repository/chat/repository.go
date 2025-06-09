package chatRepository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/dimastephen/chatServer/internal/client/db"
	"github.com/dimastephen/chatServer/internal/model"
	"github.com/dimastephen/chatServer/internal/repository"
	"log"
)

type repo struct {
	pgClient db.Client
}

func NewRepository(db db.Client) repository.ChatRepository {
	return &repo{pgClient: db}
}

func (r *repo) Create(ctx context.Context, info *model.CreateInfo) (*model.CreateInfo, error) {
	builder := sq.Insert("chatv1").Columns("users").Values(info.Usernames).
		PlaceholderFormat(sq.Dollar).Suffix("RETURNING ID")
	query, _, err := builder.ToSql()
	if err != nil {
		log.Fatal("Failed to build CreateRequest")
	}

	resp := &model.CreateInfo{}

	err = r.pgClient.DB().QueryRowContext(ctx, db.Query{QueryRaw: query, Name: "Create"}, info.Usernames).Scan(&resp.Id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *repo) Delete(ctx context.Context, info *model.DeleteInfo) error {
	query, _, err := sq.Delete("chatv1").Where(sq.Eq{"id": info.Id}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		log.Fatalf("Error in sqling Delete chat")
	}

	_, err = r.pgClient.DB().ExecContext(ctx, db.Query{QueryRaw: query, Name: "Delete"}, info.Id)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (r *repo) LogAction(ctx context.Context, create *model.CreateInfo, delete *model.DeleteInfo) error {
	if create == nil {
		action := "DELETE"
		builder := sq.Insert("chat_log").Columns("chat_id", "action").Values(delete.Id, action).
			PlaceholderFormat(sq.Dollar)
		query, _, err := builder.ToSql()
		if err != nil {
			return err
		}
		_, err = r.pgClient.DB().ExecContext(ctx, db.Query{QueryRaw: query, Name: "LOG_DELETE"}, delete.Id, action)
		if err != nil {
			return err
		}
	} else if delete == nil {
		action := "CREATE"
		builder := sq.Insert("chat_log").Columns("chat_id", "action").Values(create.Id, action).
			PlaceholderFormat(sq.Dollar)
		query, _, err := builder.ToSql()
		if err != nil {
			return err
		}
		_, err = r.pgClient.DB().ExecContext(ctx, db.Query{QueryRaw: query, Name: "LOG_CREATE"}, create.Id, action)
		if err != nil {
			return err
		}
	}
	return nil
}
