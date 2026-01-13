package postgres

import (
	"context"
	"database/sql"
	"wch/services/msgs/internal/domain/model"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Save(ctx context.Context, msg *model.Message) error {
	query := "INSERT INTO messages (id, chat_id, sender_id, content, created_at) VALUES ($1, $2, $3, $4, $5);"
	if row := r.db.QueryRow(
		query,
		msg.Id,
		msg.ChatId,
		msg.SenderId,
		msg.Content,
		msg.CreatedAt,
	); row.Err() != nil {
		return row.Err()
	}

	return nil
}

func (r *Repository) List(ctx context.Context, chat_id uuid.UUID, limit int, cursor *model.Cursor) ([]model.Message, error) {
	var rows *sql.Rows
	var err error
	if cursor != nil {
		query := `
	SELECT * FROM messages WHERE chat_id=$1
	 AND (
	 	created_at < $2 OR (created_at = $2 AND id < $3)
	) ORDER BY created_at DESC, id DESC LIMIT $4;
	`
		rows, err = r.db.Query(query, chat_id.String(), cursor.CreatedAt, cursor.ID, limit)
		if err != nil {
			return nil, err
		}
	} else {
		query := `
	SELECT * FROM messages WHERE chat_id=$1
	 ORDER BY created_at DESC, id DESC LIMIT $2;
	`
		rows, err = r.db.Query(query, chat_id.String(), limit)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	msgs := make([]model.Message, 0)
	for rows.Next() {
		var msg model.Message
		var id string
		var chat_id string
		var sender_id string

		err = rows.Scan(
			&id,
			&chat_id,
			&sender_id,
			&msg.Content,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		msg.Id = uuid.MustParse(id)
		msg.ChatId = uuid.MustParse(chat_id)
		msg.SenderId = uuid.MustParse(sender_id)

		msgs = append(msgs, msg)
	}

	return msgs, nil
}
