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

func (r *Repository) List(ctx context.Context, chat_id uuid.UUID) ([]model.Message, error) {
	query := "SELECT * FROM messages WHERE chat_id=$1"
	rows, err := r.db.Query(query, chat_id.String())
	if err != nil {
		return nil, err
	}
	rows.Close()

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
