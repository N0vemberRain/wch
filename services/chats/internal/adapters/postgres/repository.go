package repository

import (
	"context"
	"database/sql"
	"errors"
	chats "wch/services/chats/internal/domain"
	"wch/services/chats/internal/domain/model"

	"github.com/google/uuid"
)

var ErrHasntBeenDone error = errors.New("the method hasn't been done yet")

type ChatRepositoryPg struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepositoryPg {
	return &ChatRepositoryPg{db}
}

func (r *ChatRepositoryPg) CreateChat(ctx context.Context, c *model.Chat) error {
	if row := r.db.QueryRow(
		`INSERT INTO chats (id, type, name, created_at, updated_at) VALUES (
		$1, $2, $3, $4, $5
		)`,
		c.ID,
		c.Type,
		c.Name,
		c.CreatedAt,
		c.UpdatedAt,
	); row.Err() != nil {
		return row.Err()
	}

	return nil
}

func (r *ChatRepositoryPg) GetChatByID(ctx context.Context, chatID uuid.UUID) (*model.Chat, error) {
	c := &model.Chat{}
	var chat_type string
	err := r.db.QueryRow(
		`SELECT id, type, name, created_at, updated_at FROM chats WHERE id=$1`,
		chatID.String(),
	).Scan(
		&c.ID,
		&chat_type,
		&c.Name,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	c.Type, err = model.ChatTypeFromString(chat_type)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (r *ChatRepositoryPg) GetChatByName(ctx context.Context, name string) (*model.Chat, error) {
	c := &model.Chat{}
	var chat_type string
	err := r.db.QueryRow(
		`SELECT id, type, name, created_at, updated_at FROM chats WHERE name=$1`,
		name,
	).Scan(
		&c.ID,
		&chat_type,
		&c.Name,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	c.Type, err = model.ChatTypeFromString(chat_type)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *ChatRepositoryPg) UpdateChat(ctx context.Context, c *model.Chat) error {
	return ErrHasntBeenDone
}

func (r *ChatRepositoryPg) DeleteChat(ctx context.Context, chatID uuid.UUID) error {
	return ErrHasntBeenDone
}

func (r *ChatRepositoryPg) AddParticipant(
	ctx context.Context, chatID uuid.UUID, p *model.ChatParticipant,
) error {
	err := r.db.QueryRow(
		`INSERT INTO chat_participants (chat_id, user_id, role, joined_at) VALUES ($1, $2, $3, $4);`,
		p.ChatID.String(),
		p.UserID.String(),
		model.ChatParticipantRoleToString(p.Role),
		p.JoinedAt,
	).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *ChatRepositoryPg) RemoveParticipant(
	ctx context.Context, chatID uuid.UUID, userID uuid.UUID,
) error {
	err := r.db.QueryRow(
		`DELETE FROM chat_participants WHERE chat_id=$1, user_id=$2;`,
		chatID.String(),
		userID.String(),
	).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *ChatRepositoryPg) GetParticipantsIDs(
	ctx context.Context, chatID uuid.UUID,
) ([]uuid.UUID, error) {
	rows, err := r.db.Query(
		`SELECT user_id FROM chat_participants WHERE chat_id=$1;`,
		chatID.String(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids = make([]uuid.UUID, 0)
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, uuid.MustParse(id))
	}

	return ids, nil
}

func (r *ChatRepositoryPg) GetParticipant(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (
	*model.ChatParticipant,
	error,
) {
	row := r.db.QueryRow(
		`SELECT * FROM chat_participants WHERE chat_id=$1 AND user_id=$2;`,
		chatID.String(),
		userID,
	)

	if row.Err() != nil {
		return nil, row.Err()
	}

	p := &model.ChatParticipant{}
	var role string
	if err := row.Scan(&p.ChatID, &p.UserID, &role, &p.JoinedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, chats.ErrParticipateNotFound
		} else {
			return nil, err
		}
	}
	p.Role = model.ChatParticipantRoleFromString(role)

	return p, nil
}
